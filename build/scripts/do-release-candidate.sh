#!/bin/sh -ex

if [ -z "${GITHUB_REF_NAME}" ] || [ "${GITHUB_REF_TYPE}" != "tag" ] ; then
    echo "Expected a tag push event, skipping release workflow";
    exit 1;
fi

RELEASE_NOTES_FILE="docs/release_notes/${GITHUB_REF_NAME/-rc.*}.md"

if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes file ${RELEASE_NOTES_FILE} does not exist. Exiting..."
    exit 1
fi

export RELEASE_DESCRIPTION="${GITHUB_REF_NAME}"

make build-all
#goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

#docker login --username weaveworkseksctlci --password "${DOCKER_HUB_PASSWORD}"
#EKSCTL_IMAGE_VERSION="${CIRCLE_TAG}" make -f Makefile.docker push-eksctl-image || true
