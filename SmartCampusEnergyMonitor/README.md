# Smart Campus Energy Monitor ⚡️🏢

A distributed microservices system for real-time energy monitoring, anomaly detection, and alerting across a smart campus.

## 🚀 Overview

The **Smart Campus Energy Monitor** is designed to ingest high-frequency energy data from IoT sensors, detect power spikes or brownouts in real-time using statistical analysis (Z-Score), and visualize consumption metrics. It demonstrates a production-grade architecture using **gRPC** for low-latency communication and **Prometheus/Grafana** for observability.

## 🏗 Architecture

The system consists of three main microservices:

1.  **Sensor Simulator (Python)**: Generates realistic energy readings (wattage, voltage) with occasional anomalies and streams them via gRPC.
2.  **Aggregator Service (Go)**: The core engine that receives gRPC streams, calculates real-time statistics, detects anomalies (Z-Score > 3), and exposes metrics to Prometheus.
3.  **Alert Service (Python/FastAPI)**: Receives critical alert payloads from the Aggregator and logs them to an AWS S3 bucket (mocked/local for dev).

## 🛠 Tech Stack

-   **Languages**: Go (Golang), Python
-   **Communication**: gRPC (Protobuf)
-   **Monitoring**: Prometheus, Grafana, Alertmanager
-   **Infrastructure**: Docker, Docker Compose, Kubernetes (Minikube manifests included)

## ⚡️ Quick Start

### Prerequisites
-   Docker & Docker Compose

### Run Locally

1.  **Start the stack**:
    ```bash
    docker-compose up --build
    ```

2.  **Access the Dashboard**:
    -   **Grafana**: [http://localhost:3001](http://localhost:3001)
        -   **User/Pass**: `admin` / `admin`
        -   Navigate to **Dashboards > Smart Campus Energy** to see real-time heatmaps and power usage.

3.  **View Metrics & Alerts**:
    -   **Prometheus**: [http://localhost:9090](http://localhost:9090)
    -   **Alert Service**: Listens on port `8000`.

### 🧪 Simulate Anomalies
The `Sensor Simulator` is programmed to generate a power spike (~5% chance) roughly every few seconds. Watch the **Aggregator logs** to see anomalies being detected:

```bash
docker-compose logs -f aggregator-service
```

## 📦 Project Structure

```
├── aggregator-service/  # Go gRPC server & anomaly detection
├── alert-service/       # Python FastAPI for critical logging
├── sensor-simulator/    # Python gRPC client (data generator)
├── proto/               # Protobuf definitions
├── monitoring/          # Prometheus & Grafana config
├── k8s/                 # Kubernetes deployment manifests
└── docker/              # Dockerfiles
```

## ☸️ Kubernetes Deployment

Manifests are provided in the `k8s/` directory.

```bash
kubectl apply -f k8s/deployment.yaml
```

---
*Created for the "Smart Campus" Distributed Systems Project.*
