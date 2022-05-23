#!/bin/bash -ex

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
source "${DIR}/tag-common.sh"

if [ -z "${GITHUB_REF_NAME}" ] || [ "${GITHUB_REF_TYPE}" != "tag" ] ; then
  echo "Expected a tag push event, skipping release workflow"
  exit 1
fi

tag="${GITHUB_REF_NAME}"
RELEASE_NOTES_FILE="docs/release_notes/${tag/-rc.*}.md"

if [ ! -f "${RELEASE_NOTES_FILE}" ]; then
    echo "Release notes file ${RELEASE_NOTES_FILE} does not exist. Exiting..."
    exit 1
fi

# Update eksctl version to release-candidate
pre_release_id=$(release_generate pre-release-id)

echo "PRE_RELEASE_ID is ${pre_release_id}"

export RELEASE_DESCRIPTION="${tag}"

GORELEASER_CURRENT_TAG="v${tag}" PRE_RELEASE_ID="${pre_release_id}" goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"
