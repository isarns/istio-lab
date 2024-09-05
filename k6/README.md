# k6

This directory contains everything needed to run load tests using [k6](https://k6.io/) in the Kubernetes cluster.

## Contents

- **`script.js`**: This is the main file that defines the k6 load test script.
- **`test-run.yaml`**: The Kubernetes manifest that defines the job for running the k6 test in the cluster.
- **Taskfile**: Automates common tasks related to the k6 test, such as starting, deleting, and reloading the test run.

## Taskfile Commands

You can use the following commands in the Taskfile to manage the k6 tests:

- **Start the Test**: 
  Use `task start` to:
  1. Create a ConfigMap from the `script.js`.
  2. Deploy the k6 test job in the cluster.
  3. Stream logs from the k6 test.

- **Delete the Test**:
  Use `task delete` to:
  1. Delete the k6 test job.
  2. Remove the ConfigMap used for the test.

- **Reload the Test**:
  Use `task reload` to:
  1. Delete the current test.
  2. Re-deploy the test from scratch.
  3. Stream the test logs.