# OmniGuard Production Readiness & Commercial Strategy Roadmap

This document outlines the strategic technical roadmap to transform OmniGuard from an Advanced Proof-of-Concept (PoC) into a production-ready, high-scale XDR/SIEM platform capable of competing with commercial solutions.

## 1. Commercial-Grade Feature Gap Analysis
To compete with systems like Sentinel, Splunk, or IBM QRadar, OmniGuard must close the following gaps in both the Frontend and Backend.

### A. Frontend & UX Capabilities (Current UI is ~15% complete)
*   **[CRITICAL] Log Explorer / Threat Hunter:**
    - Current: None.
    - Required: A powerful search interface with Lucene-like syntax, field discovery, and time-series histograms.
*   **[CRITICAL] Detection Rule Manager:**
    - Current: None (Rules are hardcoded in Go).
    - Required: A UI to create, test (against historical logs), and deploy Sigma rules or custom Python correlation logic.
*   **[HIGH] Real-time Dashboards:**
    - Current: None.
    - Required: Drag-and-drop widgets (Pie charts, maps, heatmaps) powered by ClickHouse/OpenSearch aggregations.
*   **[MEDIUM] Asset Inventory (CMDB):**
    - Current: None.
    - Required: A view of all hosts, users, and cloud entities being monitored, with their associated risk scores.
*   **[UX] Investigation UX:**
    - Current: Static lists.
    - Required: Graph-based relationship mapping, side-by-side evidence comparison, and responsive, low-latency interactions.

### B. Missing Commercial Features
*   **Threat Intelligence Integration:** Automated ingestion of IOCs from MISP, OpenCTI, or commercial feeds.
*   **Reporting & Compliance:** Automated PDF/CSV reports for SOC2, HIPAA, or GDPR.
*   **Multi-Cluster Management:** Ability to manage multiple data collectors from a single central console.
*   **Forensic Artifact Collection:** Triggering remote agents to collect evidence during an investigation.

---

## 2. Technical Roadmap

### Phase 1: Security & Identity (Zero Trust Foundation)
- [ ] **Unified Authentication:** Enforce OIDC/JWT across *every* microservice (Search, Ingestion, SOAR, AI, Alerting).
- [ ] **Strict Tenant Isolation:** Enforce mandatory `tenant_id` filtering at the database and search engine layers.
- [ ] **Infrastructure Hardening:** mTLS between services and enabled security plugins for OpenSearch/ClickHouse.

### Phase 2: ISP-Scale Ingestion & Processing
- [ ] **High-Performance Normalization:** Refactor to use compiled schemas and horizontal scaling via Kafka partitions.
- [ ] **Advanced Correlation Engine:** Replace $O(N \times M)$ matching with the **Rete algorithm** or Hyperscan.
- [ ] **Pipeline Resilience:** Implement Dead Letter Queues (DLQ) and a Schema Registry.

### Phase 3: Operationalization & Frontend Integration
- [ ] **Live API Integration:** Replace all `mockData` with real-time backend API calls.
- [ ] **Investigative Workflow:** Build a "Case Timeline" where analysts can add notes, tags, and evidence.
- [ ] **Auth Integration:** Connect Next.js to Keycloak via `next-auth`, removing all "Mock Login" logic.
- [ ] **UX Refinement:** Dark mode optimization, global keyboard shortcuts, and performance-tuned tables.

### Phase 4: AI & SOAR Maturity (Automation)
- [ ] **Sandboxed SOAR:** Run automation actions in isolated WebAssembly or gVisor environments.
- [ ] **AI Safety:** Implement Prompt Shields to prevent injection attacks and ensure tenant-aware RAG.
- [ ] **Human-in-the-Loop:** Approval workflows for destructive actions.

### Phase 5: Observability & Scale Testing
- [ ] **Full Observability:** Prometheus/Grafana dashboards for pipeline performance (EPS, Latency).
- [ ] **ISP-Level Load Testing:** Stress test the system to 100k+ EPS and identify bottlenecks.
