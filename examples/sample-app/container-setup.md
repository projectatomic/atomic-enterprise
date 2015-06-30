# Container Setup for the Sample Application
Atomic Enterprise is available as a [Docker](https://www.docker.io) container. It
has all of the software prebuilt and pre-installed, but you do need to do a few
things to get it going.

## Download and Run Atomic Enterprise Origin
If you have not already, perform the following to (download and) run the Origin
Docker container:

[//]: # (TODO: Update image name in the future)

    $ docker run -d --name "atomic-enterprise" --net=host --privileged \
    -v /var/run/docker.sock:/var/run/docker.sock \
    openshift/origin start

[//]: # (TODO: Update sharedstatedir in the future)

Note that this won't hold any data after a restart, so you'll need to use a data
container or mount a volume at `/var/lib/openshift` to preserve that data. For
example, create a `/var/lib/openshift` folder on your Docker host, and then
start origin with the following:

[//]: # (TODO: Update image name in the future)

    $ docker run -d --name "atomic-enterprise" --net=host --privileged \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /var/lib/openshift:/var/lib/openshift \
    openshift/origin start

## Preparing the Docker Host
On your **Docker host** you will need to fetch some images. You can do so like so:

[//]: # (TODO: Update image name in the future)

    docker pull openshift/origin-docker-registry
    docker pull openshift/origin-sti-builder
    docker pull openshift/origin-deployer


This will fetch several Docker images that are used as part of the Sample
Application.

Next, be sure to follow the **Setup** instructions for the Sample Application
regarding an "insecure" Docker registry.

## Connect to the Atomic Enterprise Container
Once the container is started, you need to attach to it in order to execute
commands:

    $ docker exec -it atomic-enterprise bash

You may or may not want to change the bash prompt inside this container so that
you know where you are:

    $ PS1="atomic-enterprise-dock: [\u@\h \W]\$ "

## Get the Sample Application Code
Inside the Atomic Enterprise Docker container, you'll need to fetch some of the code
bits that are used in the sample app.

[//]: # (TODO: Update sharedstatedir in the future)
[//]: # (TODO: Update image name in the future)

    $ cd /var/lib/openshift
    $ mkdir -p examples/sample-app
    $ wget \
    https://raw.githubusercontent.com/projectatomic/atomic-enterprise/master/examples/sample-app/application-template-stibuild.json \
    -O examples/sample-app/application-template-stibuild.json

## Configure client security

[//]: # (TODO: Update sharedstatedir in the future)

    $ export CURL_CA_BUNDLE=`pwd`/openshift.local.config/master/ca.crt

For more information on this step, see [Application Deploy and Update
Flow](https://github.com/projectatomic/atomic-enterprise/blob/master/examples/sample-app/README.md#application-deploy-and-update-flow),
step #3.

## Deploy the private docker registry

    $ oadm registry --create --credentials="${KUBECONFIG}"
    $ cd examples/sample-app

For more information on this step, see [Application Build, Deploy, and Update
Flow](https://github.com/projectatomic/atomic-admin/blob/master/examples/sample-app/README.md#application-build-deploy-and-update-flow),
step #4.

## Continue With Sample Application
At this point you can continue with the steps in the [Sample
Application](https://github.com/projectatomic/atomic-enterprise/blob/master/examples/sample-app/README.md),
starting from [Application Deploy and Update
Flow](https://github.com/projectatomic/atomic-enterprise/blob/master/examples/sample-app/README.md#application-deploy-and-update-flow),
step #5.

You can watch the Atomic Enterprise logs by issuing the following on your **Docker
host**:

    $ docker attach atomic-enterprise
