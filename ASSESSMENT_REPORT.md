# OmniGuard: Enterprise Technical Assessment Report

## 1. Executive Technical Assessment
| Metric | Score | Status |
| :--- | :--- | :--- |
| **Overall Architecture** | 4/10 | **WARNING** |
| **Security** | 1/10 | **CRITICAL FAILURE** |
| **Scalability** | 3/10 | **POOR** |
| **Operational Maturity** | 2/10 | **VERY LOW** |
| **AI Safety** | 1/10 | **DANGEROUS** |
| **Production Readiness** | 0/10 | **NOT READY** |

---

## 2. Critical Findings

### [CRIT-01] Universal Authentication Bypass in Core APIs
- **Affected Module:** `services/search`, `services/ingestion`, `services/soar`, `services/ai-assistant`
- **Technical Description:** These services expose HTTP/gRPC endpoints without any authentication middleware. While an `Authenticator` exists in the `identity` service, it is not actually imported or applied to the routes in the other services.
- **Risk:** Complete unauthorized access to all SOC data and capabilities.
- **Attack Scenario:** An attacker on the network can query all logs, inject fraudulent security events, or trigger SOAR workflows by simply sending unauthenticated POST requests.
- **Impact:** Total compromise of data confidentiality and integrity.
- **Remediation:** Implement and enforce the `AuthInterceptor` across all service entry points.
- **Priority:** CRITICAL

### [CRIT-02] Cross-Tenant Data Leakage (Broken Tenant Isolation)
- **Affected Module:** `services/search`, `services/ai-assistant` (RAG)
- **Technical Description:** The Search and AI services query OpenSearch using wildcards (`normalized-events-*`) without applying filters based on the `tenant_id` from the user context.
- **Risk:** Multi-tenancy is "logic-only" at the identity layer but non-existent at the data layer.
- **Attack Scenario:** A user from Tenant A can craft a search query or ask the AI Assistant a question that retrieves sensitive security logs belonging to Tenant B.
- **Impact:** Regulatory failure (GDPR/SOC2), massive privacy breach, and total loss of customer trust.
- **Remediation:** Enforce mandatory `tenant_id` filters in every OpenSearch/ClickHouse query and use per-tenant indices/aliases.
- **Priority:** CRITICAL

### [CRIT-03] Server-Side Request Forgery (SSRF) in SOAR Workflow Executor
- **Affected Module:** `services/soar/core/executor.py`
- **Technical Description:** The `http_request` action allows arbitrary URLs and performs string formatting using the entire workflow `context` without validation.
- **Risk:** Remote Code Execution (RCE) or internal reconnaissance via SSRF.
- **Attack Scenario:** An attacker triggers a workflow with a malicious URL or context that forces the SOAR service to hit `http://169.254.169.254/latest/meta-data/` to steal cloud provider credentials.
- **Impact:** Infrastructure takeover and credential theft.
- **Remediation:** Implement a URL allowlist, restrict network egress for the SOAR worker, and sanitize context before templating.
- **Priority:** CRITICAL

### [HIGH-01] Insecure Data Store Defaults
- **Affected Module:** `deployments/docker-compose/docker-compose.yaml`
- **Technical Description:** OpenSearch is configured with `DISABLE_SECURITY_PLUGIN=true`. This disables authentication, RBAC, and TLS for the primary data store.
- **Risk:** Unauthenticated data access at the infrastructure level.
- **Impact:** Any container or user on the internal network has full administrative access to the SIEM's data.
- **Remediation:** Enable the Security Plugin, configure TLS, and rotate default credentials.
- **Priority:** HIGH

---

## 3. Architecture Review

### Strengths
- **Decoupled Microservices:** Good separation of concerns between ingestion, normalization, and correlation.
- **Standard Stack:** Use of Go, Kafka, and OpenSearch is appropriate for the domain.
- **Shared Models:** Adoption of ECS/OCSF models in `libs/models` provides a solid foundation for data normalization.

### Weaknesses & Anti-patterns
- **Broken Trust Boundaries:** Services assume internal traffic is safe, leading to a complete lack of "Zero Trust" architecture.
- **Coupling Issues:** The `normalization` and `correlation` services are tightly coupled to the Kafka schema but lack a Schema Registry.
- **Missing Capabilities:** No Dead Letter Queues (DLQ), no Rate Limiting, and no Schema Validation at the ingestion point.

---

## 4. Security Audit Report

### Vulnerabilities
- **RCE via SOAR:** Unrestricted HTTP actions can be chained to exploit internal services.
- **Log Poisoning:** The `normalization` service blindly trusts `RawLog.Payload`, which can lead to XSS in the frontend or injection in the search engine.
- **JWT Handling:** The Identity service's OIDC verifier uses `SkipClientIDCheck: true`, increasing the risk of token substitution attacks.

### Multi-Tenant Isolation Review
- **FAIL:** Isolation is purely cosmetic. The backend services do not propagate or enforce `tenant_id` at the database or streaming layers.

---

## 5. Kubernetes & Infrastructure Audit

- **Hardening Gaps:**
    - No `NetworkPolicy` to restrict inter-service communication.
    - No `PodSecurityPolicy`/`AdmissionController` to prevent root container execution.
- **HA Gaps:**
    - Kafka and OpenSearch are configured as single-node instances in the default manifests, providing zero resilience.
- **Resilience Gaps:**
    - Missing liveness/readiness probes in several service templates.
    - No resource requests/limits defined, making the cluster vulnerable to "noisy neighbor" issues and OOM kills.

---

## 6. AI Risk Assessment

- **Unsafe Workflows:** The AI Assistant has direct access to cross-tenant data via the RAG module.
- **Prompt Injection:** The system does not sanitize alert details before sending them to the LLM "Analysis" logic, allowing attackers to "jailbreak" the triage summary.
- **Explainability:** There is no mechanism to verify the "recommendations" provided by the AI, which could lead to dangerous autonomous actions (e.g., isolating a Domain Controller).

---

## 7. Detection Engineering Assessment

- **Rule Quality:** The correlation engine uses a naive `O(n*m)` matching loop. At 10,000 EPS and 500 rules, the service will likely fail or introduce massive latency.
- **Coverage:** No evidence of a rule lifecycle management system (testing, versioning, deployment).
- **Operational Viability:** Without a proper Sigma compiler that targets OpenSearch Query DSL or ClickHouse SQL directly, real-time correlation is functionally limited.

---

## 8. SOC Operations Assessment

- **Analyst Workflow:** The frontend is a "static shell". It lacks the required interactivity for real-world DFIR (e.g., no PCAP viewer, no process tree visualization, no case timeline).
- **Usability:** The mocked authentication prevents real SOC testing.
- **DFIR Readiness:** **ZERO.** No integration with forensic collectors or endpoint isolation tools (beyond the vulnerable SOAR action).

---

## 9. Production Readiness Checklist

- **High Ingestion Scaling:** **FAIL** (No Kafka tuning, single-node OpenSearch)
- **Multi-Tenancy:** **FAIL** (Critical leakage)
- **Security & RBAC:** **FAIL** (Unauthenticated APIs)
- **Observability:** **PARTIAL** (Basic metrics, no business/SIEM logic alerts)
- **Disaster Recovery:** **FAIL** (No backup strategy or persistent volume management)

---

## 10. Refactoring Recommendations

1. **Identity & Auth:** Move the `AuthInterceptor` into a shared library and enforce it as a mandatory middleware in every service.
2. **Tenant Context:** Implement a "Tenant Context" object that is passed through the entire Go call stack and used as a mandatory filter in all DB/Search drivers.
3. **SOAR Hardening:** Implement a sandboxed executor (e.g., using WebAssembly or restricted containers) for SOAR actions and add a "Human-in-the-Loop" approval gate.
4. **Data Pipeline:** Add a Schema Registry for Kafka and implement a proper DLQ for failed normalization events.

---

## 11. Immediate Remediation Roadmap

### Phase 1: Critical Fixes (Week 1)
- **Emergency:** Apply authentication middleware to ALL service endpoints.
- **Emergency:** Hardcode tenant filtering into the Search and RAG modules.
- **Security:** Enable OpenSearch security and change all default passwords.

### Phase 2: Security Hardening (Week 2-3)
- **SSRF Mitigation:** Implement a URL allowlist and egress filtering for SOAR.
- **Infrastructure:** Deploy K8s NetworkPolicies and PodSecurityContexts.
- **Auth:** Fix the OIDC client ID check and remove mocks from the Frontend.

### Phase 3: Scalability & Reliability (Week 4+)
- **Pipeline:** Implement Kafka partitions and DLQs.
- **Detection:** Refactor the Correlation engine to use a Rete-like algorithm or native DB-side matching.
- **Testing:** Target 80% unit test coverage for core pipeline and auth logic.

---
**Assessment Conclusion:** The OmniGuard platform currently exists as a "Proof of Concept" with a modern UI, but it is **technically dangerous** to deploy in any environment. The lack of authentication and tenant isolation makes it a liability rather than a security asset.
