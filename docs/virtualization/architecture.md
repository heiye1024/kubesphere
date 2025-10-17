# KubeSphere Virtualization Plugin Architecture

## Overview

The KubeSphere Virtualization plugin extends the platform with a KubeVirt-compatible virtualization stack that follows the same modular design as ks-devops. The plugin targets both host and member clusters and ensures feature parity with upstream KubeVirt components while integrating with KubeSphere RBAC, audit, multi-cluster and observability frameworks.

## Architectural Goals

- **Unified Go Backend** – All backend services are written in Go (Go 1.23+) using controller-runtime, client-go and Gin to align with the existing KubeSphere ecosystem.
- **Host Cluster Enablement** – Virtual machines can run on the host cluster when it is labeled `virtualization.kubesphere.io/enabled=true`. Host deployments directly talk to the in-cluster APIs while member clusters are accessed via cluster-proxy.
- **KubeVirt Feature Coverage** – The API surface supports CPU models, NUMA, hugepages, SR-IOV, vGPU, LiveMigration, VNC/serial console and other native KubeVirt features to avoid capability gaps.
- **KubeSphere Integration** – The solution respects workspace/project scoping, impersonation-based authorization, `/kapis` API exposure, centralized auditing and monitoring.

## Module Breakdown

### Custom Resources

CRDs under `virtualization.kubesphere.io` define virtual machines, disks, networks, snapshots and templates. They include workspace/project/cluster labels, expose power and migration conditions, and provide webhook conversions between `v1alpha1` and `v1beta1`. Validation covers NUMA placement, SR-IOV/vGPU requirements, LiveMigration settings and cluster capability labels.

### Controllers & Webhooks

Kubebuilder-generated controllers reconcile VirtualMachine, VirtualDisk, VirtualNet, VMSnapshot and VMTemplate resources. Controllers create or update KubeVirt VMs/VMIs, CDI DataVolumes and Multus NADs locally on the host or via cluster-proxy for member clusters. Admission webhooks prevent scheduling onto clusters without virtualization, reject SR-IOV/vGPU requests without matching capabilities, enforce migration prerequisites and ensure quotas.

### Aggregated REST Gateway

A Gin-based gateway exposes `/kapis/virtualization.kubesphere.io/v1beta1/...` APIs for VMs, disks, networks, snapshots and templates. Requests support the `?cluster=` selector, impersonate the caller, wrap responses in `{data,total,message,traceID,auditID}` and collect metrics and audit logs. Aggregated list queries combine host and member cluster data while respecting namespace permissions.

### Frontend Extension

The frontend registers an `ExtensionEntry` that adds virtualization menus under workspace and project scopes. Pages provide VM lifecycle management, disk and network lists, snapshot operations, template-based VM creation wizards, and VNC/serial consoles. The UI reads RBAC hints and displays backend `X-Deny-Reason` headers for disabled actions. A standalone shell allows local development without backend availability.

### Packaging & Deployment

A Helm chart (`ks-virtualization`) deploys controllers, webhooks, gateway and frontend assets. Values support `clusterSelector`, per-cluster overrides and feature flags (hostVM, liveMigration, sriov, vgpu, hugepages, numa). Documentation includes air-gap mirroring steps, image signatures (Cosign) and SBOM generation. Install plans avoid clusters without virtualization support.

### CI/CD & Testing

GitHub Actions and Jenkins pipelines lint Go/TS/YAML, run unit tests, execute kind-based integration scenarios and cover a matrix of host/member/hybrid clusters plus feature toggles (SR-IOV, vGPU, LiveMigration). Ginkgo E2E tests validate VM provisioning, migration, snapshot restore, SR-IOV rejections and RBAC enforcement.

### Observability & Operations

PrometheusRule alerts cover migration failures, DataVolume timeouts, VM heartbeat loss and disk latency. Grafana dashboards track cluster, project and VM metrics with labels `workspace`, `project`, `cluster`, `vm_name`. Loki/Tempo pipelines capture logs/traces using `trace_id`, `vm_name`, `operation`. Operations manuals document prerequisites, troubleshooting and rollback steps aligned with Helm/Operator workflows.

## Interaction Flow

1. Console calls `/kapis/virtualization.kubesphere.io/v1beta1/projects/:namespace/vms?cluster=member01`.
2. Gateway authenticates via JWT, impersonates the user and forwards through cluster-proxy when targeting member clusters.
3. Controllers reconcile CRs into KubeVirt, CDI and Multus resources.
4. Status updates propagate back through the CRDs and aggregated APIs to refresh the UI.
5. Metrics and audit entries capture user, workspace, project, cluster, operation and outcome for compliance and troubleshooting.

## Compliance Checklist Alignment

- Host/member clusters deploy virtualization workloads with feature flags and selectors.
- All APIs are served under `/kapis` and audited.
- RBAC mappings maintain workspace/project isolation; unauthorized actions return clear deny reasons.
- E2E matrix covers virtualization/no-virtualization, storage, CNI and feature toggles.
- Monitoring artifacts surface migration/import/snapshot failures; dashboards highlight resource health.
- Helm/Operator packages support rollback and backwards-compatible CRD upgrades.

