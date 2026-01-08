# Semester 3 Pico Cluster

This repository contains the infrastructure and application code for a Kubernetes-based Raspberry Pi Pico cluster project. The project demonstrates cloud-native deployment patterns, IoT sensor data processing, and enterprise-grade observability and authentication solutions.

## 📁 Project Structure

The repository is organized into two main directories representing different deployment environments:

```
Semster-3-pico-cluster/
├── .github/
│   └── workflows/          # GitHub Actions CI/CD workflows
├── Local/                  # Local cluster deployments and applications
│   ├── Grafana/           # Grafana monitoring dashboards
│   ├── Sensor-verwerker/  # Sensor data processing application (Go)
│   ├── api-server/        # REST API server for sensor data
│   ├── argocd/            # ArgoCD GitOps configuration
│   └── postgress-nativePG/ # CloudNativePG PostgreSQL cluster configuration
└── Remote/                 # Remote cluster deployments
    ├── argocd/            # ArgoCD GitOps configuration for remote cluster
    ├── authentik/         # Authentik SSO/authentication platform
    ├── glpi/              # GLPI IT asset management
    ├── grafana/           # Grafana monitoring for remote cluster
    ├── semgrep/           # Semgrep security scanning
    └── ssl/               # SSL/TLS certificate management
```

## 🔧 Components Overview

### Local Cluster (`/Local`)

The local cluster hosts the core application components and infrastructure services:

#### **Grafana** (`/Local/Grafana`)
- Monitoring and visualization dashboards
- Configured to display metrics from the sensor data processing pipeline
- Integrates with Prometheus for metrics collection

#### **Sensor-verwerker** (`/Local/Sensor-verwerker`)
- **Language**: Go
- **Purpose**: Processes IoT sensor data from Raspberry Pi Pico devices
- **Components**:
  - `main.go` - Core application logic
  - `main_test.go` - Unit tests
  - `Dockerfile` - Container image definition
  - `manifest/` - Kubernetes deployment manifests
  - `test-deployment.sh` - Deployment testing script

#### **API Server** (`/Local/api-server`)
- REST API endpoint for sensor data ingestion and retrieval
- **Configuration Files**:
  - `cluster.yml` - PostgreSQL cluster configuration
  - `deployment.yaml` - Kubernetes deployment specification
  - `api-server/` - Application source code

#### **ArgoCD** (`/Local/argocd`)
- GitOps continuous deployment tool
- Manages automated application deployments to the local cluster
- Syncs cluster state with Git repository definitions

#### **PostgreSQL CloudNativePG** (`/Local/postgress-nativePG`)
- Cloud-native PostgreSQL operator for Kubernetes
- Provides high-availability database clusters
- Stores sensor data and application state

### Remote Cluster (`/Remote`)

The remote cluster provides centralized services, authentication, and security:

#### **ArgoCD** (`/Remote/argocd`)
- GitOps deployment management for remote infrastructure
- Orchestrates application lifecycle across remote services

#### **Authentik** (`/Remote/authentik`)
- Identity provider and SSO platform
- Handles OAuth2, SAML, and LDAP authentication
- Provides centralized access management for all cluster services

#### **GLPI** (`/Remote/glpi`)
- IT Asset and Service Management platform
- Tracks hardware inventory and infrastructure components
- Manages tickets and change requests

#### **Grafana** (`/Remote/grafana`)
- Remote monitoring and observability dashboards
- Aggregates metrics from both local and remote clusters
- Provides centralized visibility across the entire infrastructure

#### **Semgrep** (`/Remote/semgrep`)
- Static application security testing (SAST)
- Automated code security scanning
- Enforces security policies across repositories

#### **SSL/TLS** (`/Remote/ssl`)
- Certificate management and SSL/TLS configuration
- Handles secure communication between services
- Manages cert-manager or manual certificate deployments

## 🚀 Deployment Architecture

### Local Cluster Workflow
1. **Sensor Data Collection**: Raspberry Pi Pico devices send sensor readings to the API server
2. **Data Processing**: The Sensor-verwerker application processes and transforms raw sensor data
3. **Data Storage**: Processed data is stored in the CloudNativePG PostgreSQL cluster
4. **Visualization**: Grafana dashboards display real-time and historical sensor metrics
5. **GitOps**: ArgoCD ensures all deployments stay synchronized with the Git repository

### Remote Cluster Workflow
1. **Authentication**: Authentik provides SSO for all remote services
2. **Monitoring**: Grafana aggregates metrics from local and remote clusters
3. **Security**: Semgrep scans code repositories for vulnerabilities
4. **Asset Management**: GLPI tracks infrastructure components and configurations
5. **Certificate Management**: SSL/TLS ensures encrypted communication

## 🛠️ Technology Stack

- **Container Orchestration**: Kubernetes
- **GitOps**: ArgoCD
- **Programming Language**: Go (Golang)
- **Database**: PostgreSQL with CloudNativePG operator
- **Monitoring**: Grafana + Prometheus
- **Authentication**: Authentik (OAuth2/SAML/LDAP)
- **Security Scanning**: Semgrep
- **IT Service Management**: GLPI
- **Hardware**: Raspberry Pi Pico cluster

## 📋 Prerequisites

- Kubernetes cluster (local and remote)
- kubectl CLI configured
- ArgoCD installed and configured
- Helm 3.x (for chart-based deployments)
- Docker (for building container images)
- Go 1.x+ (for local development of Sensor-verwerker)

## 🔐 Security Considerations

- All inter-service communication uses TLS/SSL encryption
- Authentik provides centralized authentication and authorization
- Semgrep performs automated security scanning on all code changes
- Secrets are managed through Kubernetes secrets or external secret managers
- Network policies enforce least-privilege access between services

## 📊 Monitoring & Observability

- **Metrics**: Prometheus scrapes metrics from all applications
- **Visualization**: Grafana provides dashboards for system and application metrics
- **Logs**: Centralized logging through Kubernetes logging infrastructure
- **Alerting**: Alert rules configured in Prometheus/Grafana for critical events

## 🤝 Contributing

When contributing to this repository:

1. Create a feature branch from `beta` or `main`
2. Make your changes and test locally
3. Ensure all tests pass (`go test` for Go applications)
4. Submit a pull request with a clear description of changes
5. ArgoCD will automatically sync approved changes to the cluster

## 📝 License

This project is developed as part of Semester 3 coursework.

## 📧 Contact

For questions or issues, please open a GitHub issue in this repository.
