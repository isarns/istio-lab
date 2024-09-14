# IstioLab

IstioLab is a streamlined environment designed to help you understand and experiment with [Istio](https://istio.io/). By providing pre-built tools and configurable microservices, IstioLab ensures that Istio is the primary variable in your Kubernetes environment, allowing for focused learning and experimentation in a controlled setting.

## Features

- **One-Command Setup**: Quickly set up a Kubernetes cluster using [Kind](https://kind.sigs.k8s.io/) with a single command.
- **Integrated Tools**: Automatic configuration of [Istio](https://istio.io/), [Kiali](https://kiali.io/), [Grafana](https://grafana.com/), [k6](https://k6.io/), and a custom application.
- **Pre-Built Scenarios**: Mimic real-world pod-to-pod communication with predefined scenarios.
- **Easy Teardown**: Clean up your environment effortlessly to reset and start fresh.

## Getting Started

### Prerequisites

Before you start, ensure the following tools are installed:

- [**Kind**](https://kind.sigs.k8s.io/) - For managing Kubernetes clusters.
- [**Docker**](https://www.docker.com/) - For container management.
- [**Task**](https://taskfile.dev/) - For task automation.
- [**Helm**](https://helm.sh/) - For managing Kubernetes applications.

#### Optional (For Modifying the App)

If you plan to modify the custom application, you will also need:

- [**Go**](https://go.dev/) - For application development.
- [**Kustomize**](https://kustomize.io/) - For customizing Kubernetes resources.

### Project Structure

- **[app](./app/README.md)**: Contains the Go code and Kubernetes manifests for the microservices.
- **[istio](./istio/README.md)**: Holds all Istio-related configurations and files.
- **[k6](./k6/README.md)**: Includes the k6 `script.js` and Kubernetes manifest to run load tests.
- **kind**: Contains the Kind cluster configuration files.
- **[monitoring](./monitoring/README.md)**: Includes Grafana Kubernetes manifests and Grafonnet files to define dashboards.

### Installation and Usage

#### 1. Create the Cluster

- Run the following command to set up the cluster:
  ```bash
  task create
  ```
- **First-Time Setup**: Update your `/etc/hosts` file as prompted at the end of the task. Add:
  ```
  127.0.0.1 proxy.local kiali.local grafana.local
  ```

#### 2. Run a Test

- Use the command `task test-{delay}-{count}-{scenario}` to run a test, where:
  - `{delay}`: Go duration (e.g., `20ms`, `1s`, `1h`).
  - `{count}`: Number of requests (integer).
  - `{scenario}`: Scenario identifier (`A`, `B`, `C`, or `D`).
- **Example**:
  ```bash
  task test-20ms-20-A
  ```
  This runs Scenario A with a 20ms delay between each request, totaling 20 requests.
- **Looping the Test**:
  ```bash
  while true; do task test-20ms-20-A && sleep 30; done
  ```

#### 3. View Logs

- Run `task logs-{deployment}` to view logs from a specific deployment.
  - `{deployment}`: Name of the deployment (e.g., `app-a`, `app-b`).
- **Example**:
  ```bash
  task logs-app-a
  ```

#### 4. Delete the Cluster

- Tear down the cluster and clean up resources:
  ```bash
  task delete
  ```

## Architecture

![Simple Architecture](./assets/simple%20architecture.png)

*Figure 1: Basic architecture diagram showing the deployments connected through the Istio Gateway.*

![Data Flow](./assets/kiali.gif)

*Figure 2: Data flow visualized in Kiali.*

## Scenarios

The scenarios are designed to mimic real-world pod-to-pod communication within a Kubernetes cluster. Each scenario starts at the proxy, moves on to `app-a`, and then, depending on the scenario, proceeds to other services. The end service in the communication chain runs a highly intensive computation for 5 seconds by default.

### Scenario Descriptions

- **Scenario A (`/scenarioA`)**:
  - **Flow**: `Proxy` → `App A` → `App B`.
  - **Description**: App A sends traffic to App B in a simple, one-directional flow.
  - **Details**: App B performs intensive computation for 5 seconds upon receiving a request from App A.

- **Scenario B (`/scenarioB`)**:
  - **Flow**: `Proxy` → `App A` → `App B` (with loopback).
  - **Description**: Similar to Scenario A, but App B has a feedback loop, sending a response back to itself.

- **Scenario C (`/scenarioC`)**:
  - **Flow**: `Proxy` → `App A` → `App B` → `App C`.
  - **Description**: App A sends traffic to App B, which then forwards the traffic to App C, creating a chain of requests.

- **Scenario D (`/scenarioD`)**:
  - **Flow**: `Proxy` → `App A` → `App C`.
  - **Description**: App A directly sends traffic to App C, bypassing App B.

### Visual Representation

![Scenarios](./assets/Scenarios.png)

*Figure 3: Visual representation of the different scenarios.*

## Monitoring and Observability

IstioLab comes with integrated monitoring tools to help you visualize and analyze the traffic within your cluster.

- **Kiali**: Access at `http://kiali.local/` to visualize the service mesh topology and traffic flow.
- **Grafana**: Access at `http://grafana.local/` to view performance metrics and dashboards.

> **Note**: Ensure your `/etc/hosts` file is updated as mentioned in the [Prerequisites](#prerequisites) to access these tools via the local URLs.

## Contributing

We welcome contributions from the community! If you'd like to improve the project or add new features:

1. **Fork the Repository**: Click the "Fork" button at the top of the repository page to create your own copy.
2. **Create a New Branch**: For your feature or bug fix:
   ```bash
   git checkout -b feature-name
   ```
3. **Make Changes**: Implement your feature or fix, and ensure it's well-tested.
4. **Submit a Pull Request**: Open a pull request to the main repository, describing your changes and the rationale behind them.

Feel free to open issues or feature requests to discuss ideas or report bugs.

---

## Further Reading

For deeper insights into Istio and related topics, check out [Isar Nasimov's Blog](https://isar-nasimov.medium.com/), where topics like Istio, Kubernetes, and more are covered in detail. Be sure to explore the Istio-related posts for additional tips and best practices!
