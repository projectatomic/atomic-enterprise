Atomic Enterprise Examples
==========================

This directory contains examples of using Atomic Enterprise and explaining the new concepts
available on top of Kubernetes and Docker.

* [Hello Atomic](./hello-atomic) is a simple Hello World style application that can be used to start a simple pod
* [Atomic Enterprise Sample](./sample-app) is an end-to-end application demonstrating the full
  Atomic Enterprise concept chain - images, deployments, and templates.
* [Jenkins Example](./jenkins) demonstrates how to enhance the [sample-app](./sample-app) by deploying a Jenkins pod on Atomic Enterprise and thereby enable continuous integration for incoming changes to the codebase and trigger deployments when integration succeeds.
* [Node.js echo Sample](https://github.com/openshift/nodejs-ex) highlights the simple workflow from creating project, new app from GitHub, building, deploying, running and updating.
* [Project Quotas and Resource Limits](./project-quota) demonstrates how quota and resource limits can be applied to resources in an Atomic Enterprise project.
* [Replicated Zookeper Template](./zookeeper) provides a template for an Atomic Enterprise service that exposes a simple set of primitives that distributed applications can build upon to implement higher level services for synchronization, configuration maintenance, and groups and naming.
* [Database Templates](./db-templates) provide templates for ephemeral and persistent storage on Atomic Enterprise using MongoDB, MySQL, and PostgreSQL.
* [Clustered Etcd Template](./etcd) provides a template for setting up a clustered instance of the [Etcd](https://github.com/coreos/etcd) key-value store as a service on Atomic Enterprise.
* [Configurable Git Server](./gitserver) sets up a serivce capable of automatic mirroring of Git repositories, intended for use within a container or Kubernetes pod.
