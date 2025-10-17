package dto

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	virtualizationv1beta1 "kubesphere.io/kubesphere/api/virtualization/v1beta1"
)

const (
	// LabelWorkspace marks the owning workspace for a virtualization resource.
	LabelWorkspace = "kubesphere.io/workspace"
	// LabelProject marks the owning project for a virtualization resource.
	LabelProject = "kubesphere.io/project"
	// LabelCluster marks the hosting cluster for a virtualization resource.
	LabelCluster = "kubesphere.io/cluster"
)

// Envelope represents a normalized response payload returned by the gateway.
type Envelope struct {
	Data    any    `json:"data"`
	Total   int    `json:"total,omitempty"`
	TraceID string `json:"traceID"`
	AuditID string `json:"auditID"`
	Message string `json:"message,omitempty"`
}

// VirtualMachineRequest describes VM creation payload.
type VirtualMachineRequest struct {
	Name       string                                        `json:"name"`
	Workspace  string                                        `json:"workspace"`
	Project    string                                        `json:"project"`
	Cluster    string                                        `json:"cluster"`
	CPU        string                                        `json:"cpu"`
	Memory     string                                        `json:"memory"`
	CPUModel   string                                        `json:"cpuModel,omitempty"`
	Dedicated  bool                                          `json:"dedicatedCPUPlacement,omitempty"`
	NUMAPolicy string                                        `json:"numaPolicy"`
	NUMA       *virtualizationv1beta1.NUMASpec               `json:"numa,omitempty"`
	Hugepages  string                                        `json:"hugepages"`
	GPUs       []virtualizationv1beta1.GPUDevice             `json:"gpus,omitempty"`
	Disks      []virtualizationv1beta1.VirtualMachineDisk    `json:"disks"`
	Nets       []virtualizationv1beta1.VirtualMachineNetwork `json:"nets"`
	CloudInit  *virtualizationv1beta1.CloudInitSpec          `json:"cloudInit,omitempty"`
	Console    *virtualizationv1beta1.ConsoleDevices         `json:"console,omitempty"`
	Migration  *virtualizationv1beta1.LiveMigrationSpec      `json:"liveMigration,omitempty"`
	Liveness   *virtualizationv1beta1.Probe                  `json:"livenessProbe,omitempty"`
	Readiness  *virtualizationv1beta1.Probe                  `json:"readinessProbe,omitempty"`
	PowerState virtualizationv1beta1.PowerState              `json:"powerState"`
	KubeVirt   *runtime.RawExtension                         `json:"kubeVirt,omitempty"`
}

type VirtualDiskRequest struct {
	Name         string `json:"name"`
	Workspace    string `json:"workspace"`
	Project      string `json:"project"`
	Cluster      string `json:"cluster"`
	Backing      string `json:"backing"`
	Size         string `json:"size"`
	AccessMode   string `json:"accessMode"`
	VolumeMode   string `json:"volumeMode"`
	StorageClass string `json:"storageClass"`
}

type VirtualNetRequest struct {
	Name          string `json:"name"`
	Workspace     string `json:"workspace"`
	Project       string `json:"project"`
	Cluster       string `json:"cluster"`
	NADTemplate   string `json:"nadTemplate"`
	Bandwidth     *int32 `json:"bandwidthLimit"`
	VLAN          *int32 `json:"vlan"`
	SRIOVResource string `json:"sriovResource"`
}

type VMSnapshotRequest struct {
	Name          string   `json:"name"`
	Workspace     string   `json:"workspace"`
	Project       string   `json:"project"`
	Cluster       string   `json:"cluster"`
	SourceName    string   `json:"sourceName"`
	IncludedDisks []string `json:"includedDisks"`
	RetainPolicy  string   `json:"retainPolicy"`
}

type VMTemplateRequest struct {
	Name        string                                    `json:"name"`
	Workspace   string                                    `json:"workspace"`
	Project     string                                    `json:"project"`
	Cluster     string                                    `json:"cluster"`
	Parameters  virtualizationv1beta1.TemplateParameters  `json:"parameters"`
	Constraints virtualizationv1beta1.TemplateConstraints `json:"constraints"`
	UIHints     map[string]string                         `json:"uiHints"`
}

func ToVirtualMachine(namespace string, req VirtualMachineRequest) *virtualizationv1beta1.VirtualMachine {
	labels := map[string]string{}
	if req.Workspace != "" {
		labels[LabelWorkspace] = req.Workspace
	}
	if req.Project != "" {
		labels[LabelProject] = req.Project
	}
	if req.Cluster != "" {
		labels[LabelCluster] = req.Cluster
	}
	return &virtualizationv1beta1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: virtualizationv1beta1.VirtualMachineSpec{
			CPU:                   req.CPU,
			Memory:                req.Memory,
			CPUModel:              req.CPUModel,
			DedicatedCPUPlacement: req.Dedicated,
			NUMAPolicy:            req.NUMAPolicy,
			NUMA:                  req.NUMA,
			Hugepages:             req.Hugepages,
			GPUs:                  req.GPUs,
			Disks:                 req.Disks,
			Nets:                  req.Nets,
			CloudInit:             req.CloudInit,
			Console:               req.Console,
			LiveMigration:         req.Migration,
			LivenessProbe:         req.Liveness,
			ReadinessProbe:        req.Readiness,
			PowerState:            req.PowerState,
			KubeVirt:              req.KubeVirt,
		},
	}
}

func ToVirtualDisk(namespace string, req VirtualDiskRequest) *virtualizationv1beta1.VirtualDisk {
	labels := map[string]string{}
	if req.Workspace != "" {
		labels[LabelWorkspace] = req.Workspace
	}
	if req.Project != "" {
		labels[LabelProject] = req.Project
	}
	if req.Cluster != "" {
		labels[LabelCluster] = req.Cluster
	}
	return &virtualizationv1beta1.VirtualDisk{
		ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: namespace, Labels: labels},
		Spec: virtualizationv1beta1.VirtualDiskSpec{
			Backing:      req.Backing,
			Size:         req.Size,
			AccessMode:   req.AccessMode,
			VolumeMode:   req.VolumeMode,
			StorageClass: req.StorageClass,
		},
	}
}

func ToVirtualNet(namespace string, req VirtualNetRequest) *virtualizationv1beta1.VirtualNet {
	labels := map[string]string{}
	if req.Workspace != "" {
		labels[LabelWorkspace] = req.Workspace
	}
	if req.Project != "" {
		labels[LabelProject] = req.Project
	}
	if req.Cluster != "" {
		labels[LabelCluster] = req.Cluster
	}
	return &virtualizationv1beta1.VirtualNet{
		ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: namespace, Labels: labels},
		Spec: virtualizationv1beta1.VirtualNetSpec{
			NADTemplate:    req.NADTemplate,
			BandwidthLimit: req.Bandwidth,
			VLAN:           req.VLAN,
			SRIOVResource:  req.SRIOVResource,
		},
	}
}

func ToVMSnapshot(namespace string, req VMSnapshotRequest) *virtualizationv1beta1.VMSnapshot {
	labels := map[string]string{}
	if req.Workspace != "" {
		labels[LabelWorkspace] = req.Workspace
	}
	if req.Project != "" {
		labels[LabelProject] = req.Project
	}
	if req.Cluster != "" {
		labels[LabelCluster] = req.Cluster
	}
	return &virtualizationv1beta1.VMSnapshot{
		ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: namespace, Labels: labels},
		Spec: virtualizationv1beta1.VMSnapshotSpec{
			SourceRef:     virtualizationv1beta1.NamespacedName{Name: req.SourceName, Namespace: namespace},
			IncludedDisks: req.IncludedDisks,
			RetainPolicy:  req.RetainPolicy,
		},
	}
}

func ToVMTemplate(namespace string, req VMTemplateRequest) *virtualizationv1beta1.VMTemplate {
	labels := map[string]string{}
	if req.Workspace != "" {
		labels[LabelWorkspace] = req.Workspace
	}
	if req.Project != "" {
		labels[LabelProject] = req.Project
	}
	if req.Cluster != "" {
		labels[LabelCluster] = req.Cluster
	}
	return &virtualizationv1beta1.VMTemplate{
		ObjectMeta: metav1.ObjectMeta{Name: req.Name, Namespace: namespace, Labels: labels},
		Spec: virtualizationv1beta1.VMTemplateSpec{
			Parameters:  req.Parameters,
			Constraints: req.Constraints,
			UIHints:     req.UIHints,
		},
	}
}

func FromVMList(list *virtualizationv1beta1.VirtualMachineList) []virtualizationv1beta1.VirtualMachine {
	return append([]virtualizationv1beta1.VirtualMachine{}, list.Items...)
}

func FromDiskList(list *virtualizationv1beta1.VirtualDiskList) []virtualizationv1beta1.VirtualDisk {
	return append([]virtualizationv1beta1.VirtualDisk{}, list.Items...)
}

func FromNetList(list *virtualizationv1beta1.VirtualNetList) []virtualizationv1beta1.VirtualNet {
	return append([]virtualizationv1beta1.VirtualNet{}, list.Items...)
}

func FromSnapshotList(list *virtualizationv1beta1.VMSnapshotList) []virtualizationv1beta1.VMSnapshot {
	return append([]virtualizationv1beta1.VMSnapshot{}, list.Items...)
}

func FromTemplateList(list *virtualizationv1beta1.VMTemplateList) []virtualizationv1beta1.VMTemplate {
	return append([]virtualizationv1beta1.VMTemplate{}, list.Items...)
}

func OpenAPISchema(base string) map[string]any {
	projectPath := base + "/projects/{namespace}"
	return map[string]any{
		"openapi": "3.0.0",
		"info": map[string]any{
			"title":       "KubeSphere Virtualization API",
			"version":     "v1beta1",
			"description": "Aggregated endpoints for VM, Disk, Network, Snapshot and Template management.",
		},
		"paths": map[string]any{
			projectPath + "/vms": map[string]any{
				"get":  map[string]any{"summary": "List VirtualMachines"},
				"post": map[string]any{"summary": "Create VirtualMachine"},
			},
			projectPath + "/vms/{name}:powerOn": map[string]any{
				"post": map[string]any{"summary": "Power on VirtualMachine"},
			},
			projectPath + "/vms/{name}:powerOff": map[string]any{
				"post": map[string]any{"summary": "Power off VirtualMachine"},
			},
			projectPath + "/vms/{name}:migrate": map[string]any{
				"post": map[string]any{"summary": "Migrate VirtualMachine"},
			},
			projectPath + "/vms/{name}:console": map[string]any{
				"post": map[string]any{"summary": "Open VirtualMachine console"},
			},
			projectPath + "/disks": map[string]any{
				"get":  map[string]any{"summary": "List VirtualDisks"},
				"post": map[string]any{"summary": "Create VirtualDisk"},
			},
			projectPath + "/nets": map[string]any{
				"get":  map[string]any{"summary": "List VirtualNets"},
				"post": map[string]any{"summary": "Create VirtualNet"},
			},
			projectPath + "/snapshots": map[string]any{
				"get":  map[string]any{"summary": "List VMSnapshots"},
				"post": map[string]any{"summary": "Create VMSnapshot"},
			},
			projectPath + "/templates": map[string]any{
				"get":  map[string]any{"summary": "List VMTemplates"},
				"post": map[string]any{"summary": "Create VMTemplate"},
			},
		},
	}
}
