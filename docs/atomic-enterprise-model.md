# Atomic Enterprise Model

## Overview
Atomic Enterprise extends the base Kubernetes model to provide a more feature rich development lifecycle platform.

## Build

## BuildConfig

## BuildLog

Model object that stores the logs from a particular build for later inspection.

## Deployment

A deployment is a specially annotated [replicationController](https://github.com/GoogleCloudPlatform/kubernetes/blob/master/docs/replication-controller.md), specifying the desired configuration of that controller. See the [deployments](deployments.md) document.

## DeploymentConfig

A DeploymentConfig specifies an existing deployment, triggers that can result in replacing that deployment with a new one, the strategy for doing so, and the history of such changes. See the [deployments](deployments.md) document.

## Image

Metadata added to the concept of a Docker image (such as repository, tag, environment variables, etc.) that runs in a container.

## ImageRepository

## ImageRepositoryMapping

## Template

## TemplateConfig

## Route

A named method of accessing an externally-exposed Kubernetes Service that represents an external endpoint (such as a web server, message queue, or database). See the [routing document](routing.md).

## Project

An Atomic Enterprise-level grouping of pod deployments and attendant resources. May be used for identifying components in an "application" and authorizing collaboration on same.

## User

A user identity that may be authenticated and authorized for a set of capabilities. Can correspond to an actual person or a service account. 

