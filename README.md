# k8scope

**Opinionated observability stack for Kubernetes.**

Deploy Prometheus, Grafana, Loki, Alertmanager, and OpenTelemetry Collector with battle-tested defaults in one command. Stop spending days configuring YAML — start observing your cluster in minutes.

> [!WARNING]
> k8scope is in early development. APIs and configuration may change.

## The Problem

Setting up observability on Kubernetes means:

- Configuring `kube-prometheus-stack` (4000+ lines of values.yaml)
- Adding Loki separately with its own Helm chart
- Setting up OpenTelemetry Collector pipelines for app telemetry
- Connecting datasources in Grafana manually
- Importing dashboards that may or may not work
- Writing alerting rules from scratch (or living with 200 noisy defaults)
- Figuring out retention, storage, ingress, and auth

**This takes 2-5 days for an experienced SRE.** k8scope reduces it to one command.

## Architecture

k8scope deploys a hybrid observability architecture:

- **Prometheus** scrapes Kubernetes infrastructure metrics (kubelet, kube-state-metrics, node-exporter)
- **OpenTelemetry Collector** receives application telemetry via OTLP and routes it to the appropriate backends
- **Loki** aggregates logs from both sources
- **Grafana** provides unified visualization across all signals

```
                         ┌──→ Prometheus (metrics storage)
                         │
Apps ──── OTLP ──→ OTel  ├──→ Loki (log storage)          ──→ Grafana
                  Collector                                     ↑
                         └──→ Tempo (traces - roadmap)          │
                                                                │
K8s infra ──→ Prometheus (scrape) ──────────────────────────────┘
                    │
                    └──→ Alertmanager ──→ Slack / PagerDuty / Email
```

## Deployment Modes

k8scope provides opinionated defaults for every stage of your infrastructure:

| Mode | Target | Replicas | Storage | Retention | OTel Collector | Auth |
|------|--------|----------|---------|-----------|----------------|------|
| `dev` | Local testing | 1 | Ephemeral | Session | No | None |
| `startup` | Small clusters | 1 | 10Gi PVC | 7 days | DaemonSet | Basic |
| `production` | Growing teams | 2-3 (HA) | 50Gi PVC | 30 days | DaemonSet + Gateway | Basic |
| `enterprise` | Large orgs | 2-3 (HA) | External (S3/GCS) | 90 days | Gateway (multi-tenant) | OIDC/SSO |

## Quick Start

```bash
# Install the lightweight stack for a small cluster
k8scope install --mode startup

# Preview what would be installed
k8scope install --mode production --dry-run

# Check stack health
k8scope status

# Remove everything
k8scope uninstall
```

## What Gets Installed

| Component | Purpose | Modes |
|-----------|---------|-------|
| **Prometheus** | Metrics collection and alerting engine | All |
| **Grafana** | Dashboards and visualization | All |
| **Loki** | Log aggregation | All |
| **Alertmanager** | Alert routing and deduplication | startup+ |
| **OTel Collector** | Unified telemetry pipeline (OTLP) | startup+ |
| **Node Exporter** | Host-level metrics | All |
| **kube-state-metrics** | Kubernetes object metrics | All |

### OpenTelemetry Collector Modes

The OTel Collector deployment scales with your needs:

| Mode | Deployment | Role |
|------|------------|------|
| `startup` | DaemonSet | Collects node logs and receives app OTLP data |
| `production` | DaemonSet + Gateway | DaemonSet collects, Gateway processes and routes |
| `enterprise` | Gateway (HA) | Tenant-aware routing, sampling, and filtering |

Applications instrumented with OpenTelemetry SDKs can send metrics, logs, and traces to the Collector's OTLP endpoint out of the box.

### Curated Defaults

All components come with:

- **Curated dashboards** — Not 47 generic ones, but the 5 you actually need
- **Battle-tested alerts** — Categorized as critical (pages), warning (next business day), and info (no notification)
- **Sane retention defaults** — Per mode, with storage cost estimation
- **Ingress configuration** — Ready to expose with your ingress controller

## Configuration

k8scope accepts configuration via CLI flags, a YAML config file, or both (flags override YAML):

```bash
# Using flags
k8scope install --mode production --namespace monitoring

# Using a config file
k8scope install --config k8scope.yaml
```

Example `k8scope.yaml`:

```yaml
mode: production
namespace: monitoring
kubeconfig: ~/.kube/config
```

## For GitOps Users

k8scope also publishes its Helm chart for direct use with ArgoCD, Flux, or plain Helm:

```bash
helm install k8scope oci://ghcr.io/y0s3ph/k8scope --values custom-values.yaml
```

## Roadmap

- [x] CLI scaffolding with mode-based installation plans
- [ ] Helm SDK integration for actual deployments
- [ ] OpenTelemetry Collector pipelines per mode
- [ ] Curated Grafana dashboards (cluster, node, pod, networking, logs)
- [ ] Curated alerting rules by severity
- [ ] `dev` mode with Docker Compose
- [ ] `status` command with health checks
- [ ] `upgrade` command with safe rollouts
- [ ] Ingress and TLS configuration
- [ ] Tempo integration for distributed tracing
- [ ] OIDC authentication (enterprise mode)
- [ ] External storage backends (enterprise mode)
- [ ] Multi-tenant isolation (enterprise mode)
- [ ] Helm chart publication to OCI registry

## Building from Source

```bash
go build -o k8scope ./cmd/k8scope
```

## License

Apache License 2.0 — see [LICENSE](LICENSE) for details.
