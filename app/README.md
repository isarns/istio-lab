# app

This directory handles two main components: the source code and Kubernetes manifests for the microservices.

## Kubernetes Manifests

The manifests are organized into two directories: `base` and `live`.

- **base**: Contains the foundational "skeleton" for the app deployments.
- **live**: Contains the specific app deployments. These are built by patching the base manifests and adding ConfigMaps.

You can customize the application behavior by modifying the `TIME_TO_SLEEP` value in the ConfigMaps located in the `live` directory.

### Taskfile

The following tasks are available to manage the application deployments:

- `task deployApps`: Deploys all the applications to the cluster.
- `task restartApps`: Restarts the deployed applications.

## Source Code

The `src` directory contains:

- The source code for all microservices.
- A common Dockerfile for building the services.
- Utility libraries shared across the services.

### Taskfile

In the taskfile, there are tasks to build, push, and manage Docker images for the microservices.
