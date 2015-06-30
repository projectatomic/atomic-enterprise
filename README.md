Atomic Enterprise Platform
==========================

[![GoDoc](https://godoc.org/github.com/projectatomic/atomic-enterprise?status.png)](https://godoc.org/github.com/projectatomic/atomic-enterprise)
[![Travis](https://travis-ci.org/projectatomic/atomic-enterprise.svg?branch=master)](https://travis-ci.org/projectatomic/atomic-enterprise)

This is the source repository for [Atomic Enterprise](https://projectatomic.io), based on top of [Docker](https://www.docker.io) containers and the
[Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes) container cluster manager.
Atomic Enterprise adds operational centric tools on top of Kubernetes to enable easy deployment and scaling and
long-term lifecycle maintenance for small and large teams and applications.

**Features:**

* Build web-scale applications with integrated service discovery, DNS, load balancing, failover, health checking, persistent storage, and fast scaling
* Templatize the components of your system, reuse them, and iteratively deploy them over time
* Centralized administration and management of application component libraries
  * Roll out changes to software stacks to your entire organization in a controlled fashion
* Team and user isolation of containers, builds, and network communication in an easy multi-tenancy system
  * Allow developers to run containers securely by preventing root access and isolating containers with SELinux
  * Limit, track, and manage the resources teams are using

Security!!!
-----------
You should be aware of the inherent security risks associated with performing `docker build` operations on arbitrary images as they have effective root access. **Only build and run code you trust.**

For more information on the security of containers, see these articles:

* http://opensource.com/business/14/7/docker-security-selinux
* https://docs.docker.com/articles/security/

Consider using images from trusted parties, building them yourself, or only running containers that run as non-root users.


Getting Started
---------------
The easiest way to run Atomic Enterprise is in a Docker container (Atomic Enterprise requires Docker 1.6 or higher or 1.6.2 on CentOS/RHEL):

[//]: # (TODO: Update sharedstatedir when the image is updated)

    $ sudo docker run -d --name "origin" \
        --privileged --net=host \
        -v /:/rootfs:ro -v /var/run:/var/run:rw -v /sys:/sys:ro -v /var/lib/docker:/var/lib/docker:rw \
        -v /var/lib/openshift/openshift.local.volumes:/var/lib/openshift/openshift.local.volumes \
        openshift/origin start

*Security!* Why do we need to mount your host, run privileged, and get access to your Docker directory? Atomic Enterprise runs as a host agent (like Docker) and starts and stops Docker containers, mounts remote volumes, and monitors the system (/sys) to report performance and health info. You can strip all of these options off and Atomic Enterprise will still start, but you won't be able to run pods (which is kind of the point).

Once the container is started, you can jump into a console inside the container and run the CLI.

[//]: # (TODO: Update sharedstatedir when the image is updated)

    $ sudo docker exec -it origin bash

    # Start the Atomic Enterprise integrated registry in a container
    $ oadm registry --credentials=./openshift.local.config/master/openshift-registry.kubeconfig

    # Use the CLI to login, create a project, and then create your app.
    $ oc --help
    $ oc login
    Username: test
    Password: test
    $ oc new-project test

[//]: # (TODO: Add a command to run an image)

    # See everything you just created!
    $ oc status

Any username and password are accepted by default (with no credential system configured).

You can also use the Docker container to run our CLI (`sudo docker exec -it origin cli --help`) or download the `oc` command-line client from the [releases](https://github.com/projectatomic/atomic-enterprise/releases) page Linux and login from your host with `oc login`.

You can reset your server by stopping the `origin` container and then removing it via Docker. The contents of `/var/lib/atomic-enterprise` can then be removed. 

### Next Steps

We highly recommend trying out the [Atomic Enterprise Walkthrough](https://github.com/projectatomic/atomic-enterprise/blob/master/examples/sample-app/README.md), which shows some of the lower level pieces of of Atomic Enterprise that will be the foundation for user applications.

### Troubleshooting

If you run into difficulties running Atomic Enterprise, start by reading through the [troubleshooting guide](https://github.com/projectatomic/atomic-enterprise/blob/master/docs/debugging-atomic-enterprise.md).


API
---

The Atomic Enterprise APIs are exposed at `https://localhost:8443/oapi/v1/*`.

[//]: # (TODO: Update image name when ready)

To experiment with the API, you can get a token to act as a user:

    $ sudo docker exec -it openshift-origin bash
    $ oc login
    Username: test
    Password: test
    $ oc whoami -t
    <prints a token>
    $ exit
    # from your host
    $ curl -H "Authorization: bearer <token>" https://localhost:8443/oapi/v1/...

FAQ
---

1. How does Atomic Enterprise relate to Kubernetes?

    Atomic Enterprise embeds Kubernetes and adds additional functionality to offer a simple, powerful, and
    easy-to-approach operator experience for deploying and scaling applications in containers.
    Kubernetes today is focused around composing containerized applications - Atomic Enterprise adds
    managing images and integrating them into deployment flows.  Our goal is to do
    most of that work upstream, with integration and final packaging occurring in Atomic Enterprise.

[//]: # (TODO: Add "How does Atomic Enterprise releate to Openshift")

2. What can I run on Atomic Enterprise?

    Atomic Enterprise is designed to run *any* existing Docker images.  In addition you can define builds that will produce new Docker images from a Dockerfile.

[//]: # (TODO: Update image locations in the future)

    Your application image can be easily extended with a database service with Openshift's [database images](http://docs.openshift.org/latest/using_images/db_images/overview.html). The available database images are:

    * [MySQL](https://github.com/openshift/mysql)
    * [MongoDB](https://github.com/openshift/mongodb)
    * [PostgreSQL](https://github.com/openshift/postgresql)

Contributing
------------

You can develop [locally on your host](CONTRIBUTING.md#develop-locally-on-your-host) or with a [virtual machine](CONTRIBUTING.md#develop-on-virtual-machine-using-vagrant), or if you want to just try out Atomic Enterprise [download the latest Linux server, or Windows and Mac OS X client pre-built binaries](CONTRIBUTING.md#download-from-github).

First, **get up and running with the** [**Contributing Guide**](CONTRIBUTING.md).

All contributions are welcome - Atomic Enterprise uses the Apache 2 license and does not require any contributor agreement to submit patches.  Please open issues for any bugs or problems you encounter or get involved in the [Kubernetes project](https://github.com/GoogleCloudPlatform/kubernetes) at the container runtime layer.

See [HACKING.md](https://github.com/projectatomic/atomic-enterprise/blob/master/HACKING.md) for more details on developing on Atomic Enterprise including how different tests are setup.

If you want to run the test suite, make sure you have your environment set up, and from the `atomic-enterprise` directory run:

```
# run the unit tests
$ make check

# run a simple server integration test
$ hack/test-cmd.sh

# run the integration server test suite
$ hack/test-integration.sh

# run the end-to-end test suite
$ hack/test-end-to-end.sh

# run all of the tests above
$ make test
```

You'll need [etcd](https://github.com/coreos/etcd) installed and on your path for the integration and end-to-end tests to run, and Docker must be installed to run the end-to-end tests.  To install etcd you should be able to run:

```
$ hack/install-etcd.sh
```

Some of the components of Atomic Enterprise run as Docker images, including the builders and deployment tools in `images/builder/docker/*` and 'images/deploy/*`.  To build them locally run

```
$ hack/build-images.sh
```

License
-------

Atomic Enterprise is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/).
