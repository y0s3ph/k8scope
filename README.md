# k8scope

**Opinionated observability stack for Kubernetes.**

Deploy Prometheus, Grafana, Loki, and Alertmanager with battle-tested defaults in one command. Stop spending days configuring YAML — start observing your cluster in minutes.

> [!WARNING]
> k8scope is in early development. APIs and configuration may change.

## The Problem

Setting up observability on Kubernetes means:

- Configuring `kube-prometheus-stack` (4000+ lines of values.yaml)
- Adding Loki separately with its own Helm chart
- Connecting datasources in Grafana manually
- Importing dashboards that may or may not work
- Writing alerting rules from scratch (or living with 200 noisy defaults)
- Figuring out retention, storage, ingress, and auth

**This takes 2-5 days for an experienced SRE.** k8scope reduces it to one command.

## Deployment Modes

k8scope provides opinionated defaults for every stage of your infrastructure:

| Mode | Target | Replicas | Storage | Retention | Auth |
|------|--------|----------|---------|-----------|------|
| `dev` | Local testing | 1 | Ephemeral | Session | None |
| `startup` | Small clusters | 1 | 10Gi PVC | 7 days | Basic |
| `production` | Growing teams | 2-3 (HA) | 50Gi PVC | 30 days | Basic |
| `enterprise` | Large orgs | 2-3 (HA) | External (S3/GCS) | 90 days | OIDC/SSO |

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

| Component | Purpose |
|-----------|---------|
| **Prometheus** | Metrics collection and alerting engine |
| **Grafana** | Dashboards and visualization |
| **Loki** | Log aggregation |
| **Alertmanager** | Alert routing and deduplication |
| **Node Exporter** | Host-level metrics |
| **kube-state-metrics** | Kubernetes object metrics |

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
- [ ] Curated Grafana dashboards (cluster, node, pod, networking, logs)
- [ ] Curated alerting rules by severity
- [ ] `dev` mode with Docker Compose
- [ ] `status` command with health checks
- [ ] `upgrade` command with safe rollouts
- [ ] Ingress and TLS configuration
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
