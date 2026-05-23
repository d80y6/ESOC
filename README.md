# OmniGuard: Enterprise-Grade Open Source SOC/XDR/SIEM/SOAR Platform

OmniGuard is a high-performance, horizontally scalable, multi-tenant security operations platform.

## Architecture

OmniGuard follows a microservices architecture, event-driven, and cloud-native.

### Core Components

- **Identity & Access**: Keycloak-integrated SSO/OIDC with multi-tenant RBAC.
- **Ingestion Pipeline**: High-throughput Kafka-based ingestion with Stream processing (Go).
- **Storage Layer**: Hybrid storage using OpenSearch (logs/search) and ClickHouse (metrics/analytics).
- **Correlation Engine**: Real-time Sigma-based correlation and behavioral analytics.
- **Threat Intelligence**: MISP/OpenCTI integration for IOC enrichment.
- **SOAR**: Workflow automation engine (integrating Shuffle/StackStorm concepts).
- **AI Assistant**: RAG-based LLM orchestration for alert triage and investigations.
- **Frontend**: Next.js-based enterprise SOC dashboard.

## Repository Structure

- `services/`: Microservices (Go, Python, Rust)
- `frontend/`: Next.js web application
- `libs/`: Shared libraries and utilities
- `deployments/`: K8s manifests, Helm charts, Terraform, Docker Compose
- `docs/`: Technical documentation and architecture ADRs
