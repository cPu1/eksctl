package managed

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// A Service provides methods for managing managed nodegroups
type Service struct {
	provider        v1alpha5.ClusterProvider
	clusterName     string
	stackCollection *manager.StackCollection
}

// HealthIssue represents a health issue with a managed nodegroup
type HealthIssue struct {
	Message string
	Code    string
}

// TODO use goformation types
const (
	labelsPath         = "Resources.ManagedNodeGroup.Properties.Labels"
	releaseVersionPath = "Resources.ManagedNodeGroup.Properties.ReleaseVersion"
)

// NewService creates a new Service
func NewService(provider v1alpha5.ClusterProvider, stackCollection *manager.StackCollection, clusterName string) *Service {
	return &Service{provider: provider, stackCollection: stackCollection, clusterName: clusterName}
}

// GetHealth fetches the health status for a nodegroup
func (m *Service) GetHealth(nodeGroupName string) ([]HealthIssue, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	}

	output, err := m.provider.EKS().DescribeNodegroup(input)
	if err != nil {
		if isNotFound(err) {
			return nil, errors.Wrapf(err, "could not find a managed nodegroup with name %q", nodeGroupName)
		}
		return nil, err
	}

	health := output.Nodegroup.Health
	if health == nil || len(health.Issues) == 0 {
		// No health issues
		return nil, nil
	}

	var healthIssues []HealthIssue
	for _, issue := range health.Issues {
		healthIssues = append(healthIssues, HealthIssue{
			Message: *issue.Message,
			Code:    *issue.Code,
		})
	}

	return healthIssues, nil
}

// UpdateLabels adds or removes labels for a nodegroup
func (m *Service) UpdateLabels(nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return err
	}

	newLabels, err := extractLabels(template)
	if err != nil {
		return err
	}

	for k, v := range labelsToAdd {
		newLabels[k] = v
	}

	for _, k := range labelsToRemove {
		delete(newLabels, k)
	}

	template, err = sjson.Set(template, labelsPath, newLabels)
	if err != nil {
		return err
	}

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template)
}

// GetLabels fetches the labels for a nodegroup
func (m *Service) GetLabels(nodeGroupName string) (map[string]string, error) {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return nil, err
	}
	return extractLabels(template)
}

// UpgradeNodeGroup upgrades nodegroup to the latest AMI release for the specified Kubernetes version, or
// the current Kubernetes version if the version isn't specified
func (m *Service) UpgradeNodeGroup(nodeGroupName, kubernetesVersion string, waitTimeout time.Duration) error {
	nodegroupOutput, err := m.provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
	})

	if err != nil {
		return errors.Wrap(err, "failed to describe nodegroup")
	}

	updateInput := &eks.UpdateNodegroupVersionInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
	}

	if kubernetesVersion != "" {
		_, err := semver.ParseTolerant(kubernetesVersion)
		if err != nil {
			return errors.Wrap(err, "error parsing Kubernetes version")
		}
		updateInput.Version = aws.String(kubernetesVersion)
	} else {
		updateInput.Version = nodegroupOutput.Nodegroup.Version
	}

	nodegroupUpdate, err := m.provider.EKS().UpdateNodegroupVersion(updateInput)
	if err != nil {
		return errors.Wrap(err, "error ugprading nodegroup")
	}

	if updateErrors := nodegroupUpdate.Update.Errors; len(updateErrors) > 0 {
		var allErrors []string
		for _, updateError := range updateErrors {
			allErrors = append(allErrors, updateError.String())
		}
		return errors.Errorf("failed to upgrade nodegroup:\n%v", strings.Join(allErrors, "\n"))
	}

	for _, param := range nodegroupUpdate.Update.Params {
		if *param.Type == eks.UpdateParamTypeReleaseVersion {
			newReleaseVersion := *param.Value
			if newReleaseVersion == *nodegroupOutput.Nodegroup.ReleaseVersion {
				logger.Info("nodegroup %q is already up-to-date (release version: %v)", nodeGroupName, *nodegroupOutput.Nodegroup.ReleaseVersion)
				return nil
			}
			logger.Info("upgrading nodegroup to release version %v", newReleaseVersion)
		}
	}

	if waitTimeout > 0 {
		ctx, cancelFunc := context.WithTimeout(context.Background(), waitTimeout)
		defer cancelFunc()
		return m.waitForUpdate(ctx, nodeGroupName, nodegroupUpdate.Update.Id)
	}
	return nil
}

func (m *Service) waitForUpdate(ctx context.Context, nodeGroupName string, updateID *string) error {
	logger.Debug("waiting for nodegroup update to complete (updateID: %v)", *updateID)

	const retryAfter = 5 * time.Second

	for {
		describeOutput, err := m.provider.EKS().DescribeUpdate(&eks.DescribeUpdateInput{
			Name:          aws.String(m.clusterName),
			NodegroupName: aws.String(nodeGroupName),
			UpdateId:      updateID,
		})

		if err != nil {
			describeErr := errors.Wrap(err, "error describing nodegroup update")
			if !request.IsErrorRetryable(err) {
				return describeErr
			}
			logger.Warning(describeErr.Error())
		}

		logger.Debug("DescribeUpdate output: %v", describeOutput.Update.String())

		switch status := *describeOutput.Update.Status; status {
		case eks.UpdateStatusSuccessful:
			logger.Debug("nodegroup successfully upgraded")
			return nil

		case eks.UpdateStatusCancelled:
			return errors.New("nodegroup update cancelled")

		case eks.UpdateStatusFailed:
			var aggregatedErrors []string
			for _, updateError := range describeOutput.Update.Errors {
				aggregatedErrors = append(aggregatedErrors, updateError.String())
			}
			return errors.Errorf("nodegroup update failed: %v", strings.Join(aggregatedErrors, "\n"))

		case eks.UpdateStatusInProgress:
			logger.Debug("nodegroup update in progress")

		default:
			return errors.Errorf("unexpected nodegroup update status: %q", status)

		}

		timer := time.NewTimer(retryAfter)
		select {
		case <-ctx.Done():
			timer.Stop()
			return errors.Errorf("timed out waiting for nodegroup update to complete: %v", ctx.Err())
		case <-timer.C:
		}
	}
}

func (m *Service) updateNodeGroupVersion(nodeGroupName, releaseVersion string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return err
	}

	template, err = sjson.Set(template, releaseVersionPath, releaseVersion)
	if err != nil {
		return err
	}

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template)
}

func isNotFound(err error) bool {
	awsError, ok := err.(awserr.Error)
	return ok && awsError.Code() == eks.ErrCodeResourceNotFoundException
}

var (
	kubeVersionRegex = regexp.MustCompile(`\(k8s:\s([\d.]+),`)
	amiVersionRegex  = regexp.MustCompile(`v(\d+)$`)
)

// extractKubeVersion extracts the full Kubernetes version (including patch number) from the image description
// format: "EKS Kubernetes Worker AMI with AmazonLinux2 image, (k8s: 1.13.11, docker:18.06)"
func extractKubeVersion(description string) (string, error) {
	match := kubeVersionRegex.FindStringSubmatch(description)
	if len(match) != 2 {
		return "", fmt.Errorf("expected 2 matching items; got %d", len(match))
	}
	return match[1], nil
}

// extractAMIReleaseVersion extracts the AMI release version from the image name
// the format of the image name is amazon-eks-node-1.14-v20190927
func extractAMIReleaseVersion(imageName string) (string, error) {
	match := amiVersionRegex.FindStringSubmatch(imageName)
	if len(match) != 2 {
		return "", fmt.Errorf("expected 2 matching items; got %d", len(match))
	}
	return match[1], nil
}

// TODO switch to using goformation types
func extractLabels(template string) (map[string]string, error) {
	labelsValue := gjson.Get(template, labelsPath)
	if !labelsValue.Exists() {
		return nil, errors.New("failed to find labels")
	}
	values, ok := labelsValue.Value().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for labels: %T", labelsValue.Value())
	}

	labels := make(map[string]string)
	for k, v := range values {
		value, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for label value: %T", value)
		}
		labels[k] = value
	}

	return labels, nil
}

func makeReleaseVersion(kubernetesVersion, amiVersion string) string {
	return fmt.Sprintf("%s-%s", kubernetesVersion, amiVersion)
}
