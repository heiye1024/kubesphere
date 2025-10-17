# KubeSphere Virtualization Extension

## 前置条件
- 目标集群必须启用 KVM（`/dev/kvm` 映射到工作节点）。
- SR-IOV 需在节点打上 `sriov.capable=true` 标签，且部署 SR-IOV CNI 插件。
- HugePages 在节点上通过 `hugepagesz=1G hugepages=4` 等参数预留，KubeVirt 可见 `hugepages-1Gi` 资源。
- 存储需提供支持快照的 `VolumeSnapshotClass`（CSI）。

## 安装规划
- 使用 `values.yaml` 中的 `clusterSelector` 控制生效集群，确保不具备虚拟化能力的成员集群不部署相关组件。
- `overrides[]` 支持按集群/命名空间覆盖镜像、FeatureFlag 与节点选择器。

## 常见错误与修复
| 错误 | 原因 | 修复 |
| --- | --- | --- |
| `LiveMigration blocked: storage not shared` | 源宿主机未共享同一 RWX 存储 | 配置 RWX 存储或启用块迁移 | 
| `no sriov.capable=true nodes available` | Admission 校验发现节点缺标签 | 在节点上添加标签并部署 PF/VF | 
| `Hugepages request unsatisfied` | 节点未分配 HugePages | 调整内核启动参数并重启节点 |

## 备份与恢复
1. 使用 VMSnapshot 结合 CSI Snapshot 生成点时间备份。
2. 利用 Velero 配置 `Schedule` 备份 `virtualization.kubesphere.io` 资源以及 `VolumeSnapshotContent`。
3. 恢复时先导入存储卷，再回放 `VirtualMachine` 资源。

## 灰度与回滚
- Chart/Operator 支持滚动升级：先在灰度集群设置 `featureFlags.liveMigration=true` 验证。
- 回滚时使用 `helm rollback`，CRD 向后兼容（v1beta1 storage version + conversion webhook）。
- 扩展版本升级时遵循 `kubectl cert-manager` 证书轮转流程。

## Air-gap 安装
1. 使用 `values.yaml` 中 `airgap.registries` 列出私有仓库地址。
2. 执行 `make virtualization-build` 构建镜像，推送到离线仓库。
3. 校验 SBOM：`syft kubesphere/virtualization-controller:latest`。
4. Cosign 验证：`cosign verify --key <key> kubesphere/virtualization-controller:latest`。

## 观测与告警
- PrometheusRule 位于 `config/monitoring/prometheusrule.yaml`。
- Grafana Dashboard JSON 位于 `config/monitoring/grafana-dashboard.json`。
- 日志与链路追踪字段要求见 `config/monitoring/logging-tempo.md`。

