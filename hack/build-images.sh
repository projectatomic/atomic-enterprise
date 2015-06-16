#!/bin/bash

# This script builds all images locally except the base and release images,
# which are handled by hack/build-base-images.sh.

set -o errexit
set -o nounset
set -o pipefail

OS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${OS_ROOT}/hack/common.sh"

# Go to the top of the tree.
cd "${OS_ROOT}"

# Get the latest Linux release
if [[ ! -d _output/local/releases ]]; then
  echo "No release has been built. Run hack/build-release.sh"
  exit 1
fi

# Extract the release achives to a staging area.
os::build::detect_local_release_tars "linux"

echo "Building images from release tars for commit ${OS_RELEASE_COMMIT}:"
echo " primary: $(basename ${OS_PRIMARY_RELEASE_TAR})"
echo " image:   $(basename ${OS_IMAGE_RELEASE_TAR})"

imagedir="_output/imagecontext"
rm -rf "${imagedir}"
mkdir -p "${imagedir}"
tar xzf "${OS_PRIMARY_RELEASE_TAR}" -C "${imagedir}"
tar xzf "${OS_IMAGE_RELEASE_TAR}" -C "${imagedir}"

# Copy primary binaries to the appropriate locations.
cp -f "${imagedir}/openshift" images/origin/bin
cp -f "${imagedir}/openshift" images/router/haproxy/bin
cp -f "${imagedir}/openshift" images/ipfailover/keepalived/bin

# Copy image binaries to the appropriate locations.
cp -f "${imagedir}/pod" images/pod/bin
cp -f "${imagedir}/hello-atomic" examples/hello-atomic/bin
cp -f "${imagedir}/dockerregistry" images/dockerregistry/bin

# builds an image and tags it two ways - with latest, and with the release tag
function image {
  echo "--- $1 ---"
  docker build -t $1:latest $2
  docker tag -f $1:latest $1:${OS_RELEASE_COMMIT}
}

# images that depend on scratch
image openshift/origin-pod                   images/pod
# images that depend on openshift/origin-base
image openshift/origin                       images/origin
image openshift/origin-haproxy-router        images/router/haproxy
image openshift/origin-keepalived-ipfailover images/ipfailover/keepalived
image openshift/origin-docker-registry       images/dockerregistry
# images that depend on openshift/origin
image openshift/origin-deployer              images/deployer
image openshift/origin-docker-builder        images/builder/docker/docker-builder
image openshift/origin-gitserver             examples/gitserver
image openshift/origin-sti-builder           images/builder/docker/sti-builder
# extra images (not part of infrastructure)
image openshift/hello-atomic              examples/hello-atomic
# unpublished images
image openshift/origin-custom-docker-builder images/builder/docker/custom-docker-builder
image openshift/sti-image-builder            images/builder/docker/sti-image-builder

echo "++ Active images"
docker images | grep openshift/
