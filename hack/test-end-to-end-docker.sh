#!/bin/bash

# This script tests the high level end-to-end functionality demonstrated
# as part of the examples/sample-app

set -o errexit
set -o nounset
set -o pipefail

OS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${OS_ROOT}/hack/util.sh"

echo "[INFO] Starting containerized end-to-end test"

# Use either the latest release built images, or latest.
if [[ -z "${USE_IMAGES-}" ]]; then
  tag="latest"
  if [[ -e "${OS_ROOT}/_output/local/releases/.commit" ]]; then
    COMMIT="$(cat "${OS_ROOT}/_output/local/releases/.commit")"
    tag="${COMMIT}"
  fi
  USE_IMAGES="openshift/origin-\${component}:${tag}"
fi

unset KUBECONFIG

if [[ -z "${BASETMPDIR-}" ]]; then
  TMPDIR="${TMPDIR:-"/tmp"}"
  BASETMPDIR="${TMPDIR}/openshift-e2e/containerized"
  sudo rm -rf "${BASETMPDIR}"
  mkdir -p "${BASETMPDIR}"
fi
VOLUME_DIR="${BASETMPDIR}/volumes"
FAKE_HOME_DIR="${BASETMPDIR}/openshift.local.home"
LOG_DIR="${LOG_DIR:-${BASETMPDIR}/logs}"
ARTIFACT_DIR="${ARTIFACT_DIR:-${BASETMPDIR}/artifacts}"
mkdir -p $LOG_DIR
mkdir -p $ARTIFACT_DIR

GO_OUT="${OS_ROOT}/_output/local/go/bin"

# set path so OpenShift is available
export PATH="${GO_OUT}:${PATH}"


function cleanup()
{
  out=$?
  echo
  if [ $out -ne 0 ]; then
    echo "[FAIL] !!!!! Test Failed !!!!"
  else
    echo "[INFO] Test Succeeded"
  fi
  echo

  set +e
  echo "[INFO] Dumping container logs to ${LOG_DIR}"
  docker logs origin >"${LOG_DIR}/openshift.log" 2>&1

  if [[ -z "${SKIP_TEARDOWN-}" ]]; then
    echo "[INFO] Tearing down test"
    docker stop origin
    docker rm origin

    echo "[INFO] Stopping k8s docker containers"; docker ps | awk 'index($NF,"k8s_")==1 { print $1 }' | xargs -l -r docker stop
    if [[ -z "${SKIP_IMAGE_CLEANUP-}" ]]; then
      echo "[INFO] Removing k8s docker containers"; docker ps -a | awk 'index($NF,"k8s_")==1 { print $1 }' | xargs -l -r docker rm
    fi
    set -u
  fi
  set -e

  echo "[INFO] Exiting"
  exit $out
}

trap "exit" INT TERM
trap "cleanup" EXIT

out=$(
  set +e
  docker stop origin 2>&1
  docker rm origin 2>&1
  set -e
)

# Setup
echo "[INFO] `openshift version`"
echo "[INFO] Using images:              ${USE_IMAGES}"

echo "[INFO] Starting OpenShift containerized server"
sudo docker run -d --name="origin" \
  --privileged --net=host \
  -v /:/rootfs:ro -v /var/run:/var/run:rw -v /sys:/sys:ro -v /var/lib/docker:/var/lib/docker:rw \
  -v "/var/lib/openshift/openshift.local.volumes:/var/lib/openshift/openshift.local.volumes" \
  "openshift/origin:${tag}" start --images="${USE_IMAGES}"

export HOME="${FAKE_HOME_DIR}"
# This directory must exist so Docker can store credentials in $HOME/.dockercfg
mkdir -p ${FAKE_HOME_DIR}

CURL_EXTRA="-k"

wait_for_url "https://localhost:8443/healthz/ready" "apiserver(ready): " 0.25 160

# install the router
echo "[INFO] Installing the router"
sudo docker exec origin openshift admin router --create --credentials="./openshift.local.config/master/openshift-router.kubeconfig" --images="${USE_IMAGES}"

# install the registry. The --mount-host option is provided to reuse local storage.
echo "[INFO] Installing the registry"
sudo docker exec origin openshift admin registry --create --credentials="./openshift.local.config/master/openshift-registry.kubeconfig" --images="${USE_IMAGES}"

registry="$(dig @localhost "docker-registry.default.svc.cluster.local." +short A | head -n 1)"
[ -n "${registry}" ]
echo "[INFO] Verifying the docker-registry is up at ${registry}"
wait_for_url_timed "http://${registry}:5000/healthz" "[INFO] Docker registry says: " $((2*TIME_MIN))


echo "[INFO] Login"
oc login localhost:8443 -u test -p test --insecure-skip-tls-verify
oc new-project test

echo "[INFO] Applying application config"
oc new-app -f examples/hello-atomic/all-in-one.tmpl.yaml

wait_for_command "oc get -n test pods -l name=hello-atomic | grep -i Running" $((60*TIME_SEC))

echo "[INFO] Validating app response..."
app_svc_ip=$(oc get --template="{{ .spec.portalIP }}:{{ (index .spec.ports 0).port }}" service hello-atomic-service)
validate_response "-s -k http://${app_svc_ip}" "Hello Atomic!" 0.2 50
