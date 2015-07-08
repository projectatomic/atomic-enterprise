#!/bin/bash

# This script tests the high level end-to-end functionality demonstrated
# as part of the examples/sample-app

if [[ -z "$(which iptables)" ]]; then
	echo "IPTables not found - the end-to-end test requires a system with iptables for Kubernetes services."
	exit 1
fi
iptables --list > /dev/null 2>&1
if [ $? -ne 0 ]; then
	sudo iptables --list > /dev/null 2>&1
	if [ $? -ne 0 ]; then
		echo "You do not have iptables or sudo privileges.	Kubernetes services will not work without iptables access.	See https://github.com/GoogleCloudPlatform/kubernetes/issues/1859.	Try 'sudo hack/test-end-to-end.sh'."
		exit 1
	fi
fi

set -o errexit
set -o nounset
set -o pipefail

OS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${OS_ROOT}/hack/util.sh"

echo "[INFO] Starting end-to-end test"

# Use either the latest release built images, or latest.
if [[ -z "${USE_IMAGES-}" ]]; then
	USE_IMAGES='openshift/origin-${component}:latest'
	if [[ -e "${OS_ROOT}/_output/local/releases/.commit" ]]; then
		COMMIT="$(cat "${OS_ROOT}/_output/local/releases/.commit")"
		USE_IMAGES="openshift/origin-\${component}:${COMMIT}"
	fi
fi

ROUTER_TESTS_ENABLED="${ROUTER_TESTS_ENABLED:-true}"
TEST_ASSETS="${TEST_ASSETS:-false}"

if [[ -z "${BASETMPDIR-}" ]]; then
	TMPDIR="${TMPDIR:-"/tmp"}"
	BASETMPDIR="${TMPDIR}/openshift-e2e"
	sudo rm -rf "${BASETMPDIR}"
	mkdir -p "${BASETMPDIR}"
fi
ETCD_DATA_DIR="${BASETMPDIR}/etcd"
VOLUME_DIR="${BASETMPDIR}/volumes"
FAKE_HOME_DIR="${BASETMPDIR}/openshift.local.home"
LOG_DIR="${LOG_DIR:-${BASETMPDIR}/logs}"
ARTIFACT_DIR="${ARTIFACT_DIR:-${BASETMPDIR}/artifacts}"
mkdir -p $LOG_DIR
mkdir -p $ARTIFACT_DIR

DEFAULT_SERVER_IP=`ifconfig | grep -Ev "(127.0.0.1|172.17.42.1)" | grep "inet " | head -n 1 | sed 's/adr://' | awk '{print $2}'`
API_HOST="${API_HOST:-${DEFAULT_SERVER_IP}}"
API_PORT="${API_PORT:-8443}"
API_SCHEME="${API_SCHEME:-https}"
MASTER_ADDR="${API_SCHEME}://${API_HOST}:${API_PORT}"
PUBLIC_MASTER_HOST="${PUBLIC_MASTER_HOST:-${API_HOST}}"
KUBELET_SCHEME="${KUBELET_SCHEME:-https}"
KUBELET_HOST="${KUBELET_HOST:-127.0.0.1}"
KUBELET_PORT="${KUBELET_PORT:-10250}"

SERVER_CONFIG_DIR="${BASETMPDIR}/openshift.local.config"
MASTER_CONFIG_DIR="${SERVER_CONFIG_DIR}/master"
NODE_CONFIG_DIR="${SERVER_CONFIG_DIR}/node-${KUBELET_HOST}"

# use the docker bridge ip address until there is a good way to get the auto-selected address from master
# this address is considered stable
# used as a resolve IP to test routing
CONTAINER_ACCESSIBLE_API_HOST="${CONTAINER_ACCESSIBLE_API_HOST:-172.17.42.1}"

AIO_CONFIG_FILE="${LOG_DIR}/all-in-one-config.json"
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
	for container in $(docker ps -aq); do
		docker logs "$container" >&"${LOG_DIR}/container-$container.log"
	done

	echo "[INFO] Dumping build log to ${LOG_DIR}"

	oc get -n test builds --output-version=v1beta3 -t '{{ range .items }}{{.metadata.name}}{{ "\n" }}{{end}}' | xargs -r -l oc build-logs -n test >"${LOG_DIR}/stibuild.log"
	oc get -n docker builds --output-version=v1beta3 -t '{{ range .items }}{{.metadata.name}}{{ "\n" }}{{end}}' | xargs -r -l oc build-logs -n docker >"${LOG_DIR}/dockerbuild.log"
	oc get -n custom builds --output-version=v1beta3 -t '{{ range .items }}{{.metadata.name}}{{ "\n" }}{{end}}' | xargs -r -l oc build-logs -n custom >"${LOG_DIR}/custombuild.log"

	echo "[INFO] Dumping etcd contents to ${ARTIFACT_DIR}/etcd_dump.json"
	set_curl_args 0 1
	curl ${clientcert_args} -L "${API_SCHEME}://${API_HOST}:4001/v2/keys/?recursive=true" > "${ARTIFACT_DIR}/etcd_dump.json"
	echo

	if [[ -z "${SKIP_TEARDOWN-}" ]]; then
		echo "[INFO] Deleting test constructs"
		oc delete -n test all --all
		oc delete -n docker all --all
		oc delete -n custom all --all
		oc delete -n cache all --all
		oc delete -n default all --all

		echo "[INFO] Tearing down test"
		pids="$(jobs -pr)"
		echo "[INFO] Children: ${pids}"
		sudo kill ${pids}
		sudo ps f
		set +u
		echo "[INFO] Stopping k8s docker containers"; docker ps | awk 'index($NF,"k8s_")==1 { print $1 }' | xargs -l -r docker stop
		if [[ -z "${SKIP_IMAGE_CLEANUP-}" ]]; then
			echo "[INFO] Removing k8s docker containers"; docker ps -a | awk 'index($NF,"k8s_")==1 { print $1 }' | xargs -l -r docker rm
		fi
		set -u
	fi
	set -e

	# clean up zero byte log files
	# Clean up large log files so they don't end up on jenkins
	find ${ARTIFACT_DIR} -name *.log -size +20M -exec echo Deleting {} because it is too big. \; -exec rm -f {} \;
	find ${LOG_DIR} -name *.log -size +20M -exec echo Deleting {} because it is too big. \; -exec rm -f {} \;
	find ${LOG_DIR} -name *.log -size 0 -exec echo Deleting {} because it is empty. \; -exec rm -f {} \;

	echo "[INFO] Exiting"
	exit $out
}

trap "exit" INT TERM
trap "cleanup" EXIT

# Setup
stop_openshift_server
echo "[INFO] `openshift version`"
echo "[INFO] Server logs will be at:    ${LOG_DIR}/openshift.log"
echo "[INFO] Test artifacts will be in: ${ARTIFACT_DIR}"
echo "[INFO] Volumes dir is:            ${VOLUME_DIR}"
echo "[INFO] Config dir is:             ${SERVER_CONFIG_DIR}"
echo "[INFO] Using images:              ${USE_IMAGES}"

# Start All-in-one server and wait for health
echo "[INFO] Create certificates for the OpenShift server"
# find the same IP that openshift start will bind to.  This allows access from pods that have to talk back to master
ALL_IP_ADDRESSES=`ifconfig | grep "inet " | sed 's/adr://' | awk '{print $2}'`
SERVER_HOSTNAME_LIST="${PUBLIC_MASTER_HOST},localhost"
while read -r IP_ADDRESS
do
	SERVER_HOSTNAME_LIST="${SERVER_HOSTNAME_LIST},${IP_ADDRESS}"
done <<< "${ALL_IP_ADDRESSES}"

openshift admin create-master-certs \
	--overwrite=false \
	--cert-dir="${MASTER_CONFIG_DIR}" \
	--hostnames="${SERVER_HOSTNAME_LIST}" \
	--master="${MASTER_ADDR}" \
	--public-master="${API_SCHEME}://${PUBLIC_MASTER_HOST}:${API_PORT}"

openshift admin create-node-config \
	--listen="${KUBELET_SCHEME}://0.0.0.0:${KUBELET_PORT}" \
	--node-dir="${NODE_CONFIG_DIR}" \
	--node="${KUBELET_HOST}" \
	--hostnames="${KUBELET_HOST}" \
	--master="${MASTER_ADDR}" \
	--node-client-certificate-authority="${MASTER_CONFIG_DIR}/ca.crt" \
	--certificate-authority="${MASTER_CONFIG_DIR}/ca.crt" \
	--signer-cert="${MASTER_CONFIG_DIR}/ca.crt" \
	--signer-key="${MASTER_CONFIG_DIR}/ca.key" \
	--signer-serial="${MASTER_CONFIG_DIR}/ca.serial.txt"

oadm create-bootstrap-policy-file --filename="${MASTER_CONFIG_DIR}/policy.json"

# create openshift config
openshift start \
	--write-config=${SERVER_CONFIG_DIR} \
	--create-certs=false \
    --listen="${API_SCHEME}://0.0.0.0:${API_PORT}" \
    --master="${MASTER_ADDR}" \
    --public-master="${API_SCHEME}://${PUBLIC_MASTER_HOST}:${API_PORT}" \
    --hostname="${KUBELET_HOST}" \
    --volume-dir="${VOLUME_DIR}" \
    --etcd-dir="${ETCD_DATA_DIR}" \
    --images="${USE_IMAGES}"


echo "[INFO] Starting OpenShift server"
sudo env "PATH=${PATH}" OPENSHIFT_PROFILE=web OPENSHIFT_ON_PANIC=crash openshift start \
	--master-config=${MASTER_CONFIG_DIR}/master-config.yaml \
	--node-config=${NODE_CONFIG_DIR}/node-config.yaml \
    --loglevel=4 \
    &> "${LOG_DIR}/openshift.log" &
OS_PID=$!

export HOME="${FAKE_HOME_DIR}"
# This directory must exist so Docker can store credentials in $HOME/.dockercfg
mkdir -p ${FAKE_HOME_DIR}

export KUBECONFIG="${MASTER_CONFIG_DIR}/admin.kubeconfig"
CLUSTER_ADMIN_CONTEXT=$(oc config view --flatten -o template -t '{{index . "current-context"}}')

if [[ "${API_SCHEME}" == "https" ]]; then
	export CURL_CA_BUNDLE="${MASTER_CONFIG_DIR}/ca.crt"
	export CURL_CERT="${MASTER_CONFIG_DIR}/admin.crt"
	export CURL_KEY="${MASTER_CONFIG_DIR}/admin.key"

	# Make oc use ${MASTER_CONFIG_DIR}/admin.kubeconfig, and ignore anything in the running user's $HOME dir
	sudo chmod -R a+rwX "${KUBECONFIG}"
	echo "[INFO] To debug: export KUBECONFIG=$KUBECONFIG"
fi


wait_for_url "${KUBELET_SCHEME}://${KUBELET_HOST}:${KUBELET_PORT}/healthz" "[INFO] kubelet: " 0.5 60
wait_for_url "${API_SCHEME}://${API_HOST}:${API_PORT}/healthz" "apiserver: " 0.25 80
wait_for_url "${API_SCHEME}://${API_HOST}:${API_PORT}/healthz/ready" "apiserver(ready): " 0.25 80
wait_for_url "${API_SCHEME}://${API_HOST}:${API_PORT}/api/v1beta3/nodes/${KUBELET_HOST}" "apiserver(nodes): " 0.25 80

# add e2e-user as a viewer for the default namespace so we can see infrastructure pieces appear
openshift admin policy add-role-to-user view e2e-user --namespace=default

# create test project so that this shows up in the console
openshift admin new-project test --description="This is an example project to demonstrate OpenShift v3" --admin="e2e-user"
openshift admin new-project docker --description="This is an example project to demonstrate OpenShift v3" --admin="e2e-user"
openshift admin new-project custom --description="This is an example project to demonstrate OpenShift v3" --admin="e2e-user"
openshift admin new-project cache --description="This is an example project to demonstrate OpenShift v3" --admin="e2e-user"

echo "The console should be available at ${API_SCHEME}://${PUBLIC_MASTER_HOST}:${API_PORT}/console."
echo "Log in as 'e2e-user' to see the 'test' project."

# install the router
echo "[INFO] Installing the router"
openshift admin router --create --credentials="${MASTER_CONFIG_DIR}/openshift-router.kubeconfig" --images="${USE_IMAGES}"

# install the registry. The --mount-host option is provided to reuse local storage.
echo "[INFO] Installing the registry"
openshift admin registry --create --credentials="${MASTER_CONFIG_DIR}/openshift-registry.kubeconfig" --images="${USE_IMAGES}"

echo "[INFO] Pre-pulling and pushing hello-atomic"
docker pull atomicenterprise/hello-atomic:latest
echo "[INFO] Pulled hello-atomic"

echo "[INFO] Waiting for Docker registry pod to start"
# TODO: simplify when #4702 is fixed upstream
wait_for_command '[[ "$(oc get endpoints docker-registry --output-version=v1beta3 -t "{{ if .subsets }}{{ len .subsets }}{{ else }}0{{ end }}" || echo "0")" != "0" ]]' $((5*TIME_MIN))

# services can end up on any IP.	Make sure we get the IP we need for the docker registry
DOCKER_REGISTRY=$(oc get --output-version=v1beta3 --template="{{ .spec.portalIP }}:{{ with index .spec.ports 0 }}{{ .port }}{{ end }}" service docker-registry)

# TODO: do this in a temporary manner that also deals with commented out lines
sed -i "s/'--insecure-registry .*'/'--insecure-registry ${DOCKER_REGISTRY}'/" /etc/sysconfig/docker
systemctl restart docker

registry="$(dig @${API_HOST} "docker-registry.default.svc.cluster.local." +short A | head -n 1)"
[[ -n "${registry}" && "${registry}:5000" == "${DOCKER_REGISTRY}" ]]

echo "[INFO] Verifying the docker-registry is up at ${DOCKER_REGISTRY}"
wait_for_url_timed "http://${DOCKER_REGISTRY}/healthz" "[INFO] Docker registry says: " $((2*TIME_MIN))

[ "$(dig @${API_HOST} "docker-registry.default.local." A)" ]

# Client setup (log in as e2e-user and set 'test' as the default project)
# This is required to be able to push to the registry!
echo "[INFO] Logging in as a regular user (e2e-user:pass) with project 'test'..."
oc login -u e2e-user -p pass
[ "$(oc whoami | grep 'e2e-user')" ]
oc project test
token=$(oc config view --flatten -o template -t '{{with index .users 0}}{{.user.token}}{{end}}')
[[ -n ${token} ]]

echo "[INFO] Docker login as e2e-user to ${DOCKER_REGISTRY}"
docker login -u e2e-user -p ${token} -e e2e-user@openshift.com ${DOCKER_REGISTRY}
echo "[INFO] Docker login successful"

echo "[INFO] Tagging and pushing hello-atomic to ${DOCKER_REGISTRY}/test/hello-atomic"
docker tag -f atomicenterprise/hello-atomic:latest ${DOCKER_REGISTRY}/test/hello-atomic:latest
docker push ${DOCKER_REGISTRY}/test/hello-atomic:latest
echo "[INFO] Pushed hello-openshift"

echo "[INFO] Back to 'default' project with 'admin' user..."
oc project ${CLUSTER_ADMIN_CONTEXT}
[ "$(oc whoami | grep 'system:admin')" ]

# The build requires a dockercfg secret in the builder service account in order
# to be able to push to the registry.  Make sure it exists first.
echo "[INFO] Waiting for dockercfg secrets to be generated in project 'test' before building"
wait_for_command "oc get -n test serviceaccount/builder -o yaml | grep dockercfg > /dev/null" $((60*TIME_SEC))

# Process template and create

echo "[INFO] Back to 'test' context with 'e2e-user' user"
oc project test

# create a deployment, service, and route
oc process -n test -f examples/hello-atomic/all-in-one.tmpl.yaml -v "IMAGE=${DOCKER_REGISTRY}/test/hello-atomic:latest" > "${AIO_CONFIG_FILE}"
oc create -f "${AIO_CONFIG_FILE}"
wait_for_command "oc get -n test pods -l name=hello-atomic | grep -i Running" $((60*TIME_SEC))

echo "[INFO] Back to 'default' project with 'admin' user..."
oc project ${CLUSTER_ADMIN_CONTEXT}

# ensure the router is started
# TODO: simplify when #4702 is fixed upstream
wait_for_command '[[ "$(oc get endpoints router --output-version=v1beta3 -t "{{ if .subsets }}{{ len .subsets }}{{ else }}0{{ end }}" || echo "0")" != "0" ]]' $((5*TIME_MIN))

echo "[INFO] Validating routed app response..."
validate_response "-s -k --resolve www.example.com:443:${CONTAINER_ACCESSIBLE_API_HOST} https://www.example.com" "Hello Atomic!" 0.2 50

echo "[INFO] Confirming service is listed in DNS correctly..."
dns_app_svc_ip="$(dig @${API_HOST} "hello-atomic-service.test.svc.cluster.local." +short A | head -n 1)"
app_svc_ip="$(oc get service -n test -o template hello-atomic-service --template="{{.spec.portalIP}}")"
[[ -n "${app_svc_ip}" && "${dns_app_svc_ip}" == "${app_svc_ip}" ]]
echo "[INFO] Service correctly listed at ${dns_app_svc_ip}"

# Remote command execution
echo "[INFO] Validating exec"
registry_pod=$(oc get pod -l deploymentconfig=docker-registry -t '{{(index .items 0).metadata.name}}')
# when running as a restricted pod the registry will run with a pre-allocated
# user in the neighborhood of 1000000+.  Look for a substring of the pre-allocated uid range
oc exec -p ${registry_pod} id | grep 10

# Port forwarding (needs socat)
echo "[INFO] Validating port-forward"
oc port-forward -p ${registry_pod} 5001:5000  &> "${LOG_DIR}/port-forward.log" &
wait_for_url_timed "http://localhost:5001/healthz" "[INFO] Docker registry says: " $((10*TIME_SEC))

# UI e2e tests can be found in assets/test/e2e
if [[ "$TEST_ASSETS" == "true" ]]; then
	echo "[INFO] Running UI e2e tests..."
	pushd ${OS_ROOT}/assets > /dev/null
		grunt test-e2e
	popd > /dev/null
fi
