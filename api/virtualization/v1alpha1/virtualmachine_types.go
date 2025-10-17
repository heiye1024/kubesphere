package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PowerState defines desired VM power state
// +kubebuilder:validation:Enum=Running;Stopped
type PowerState string

const (
	PowerStateRunning PowerState = "Running"
	PowerStateStopped PowerState = "Stopped"
)

const (
	// ConditionTypePowerOperation tracks power transitions.
	ConditionTypePowerOperation = "PowerOperation"
	// ConditionTypeMigration reports live migration progress.
	ConditionTypeMigration = "Migration"
)

const (
	// ConditionReasonProgressing indicates an ongoing operation.
	ConditionReasonProgressing = "Progressing"
	// ConditionReasonReady indicates completion.
	ConditionReasonReady = "Ready"
	// ConditionReasonFailed indicates failure.
	ConditionReasonFailed = "Failed"
)

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=vm;vms,scope=Namespaced
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Power",type=string,JSONPath=`.status.powerState`
//+kubebuilder:printcolumn:name="CPU",type=string,JSONPath=`.spec.cpu`
//+kubebuilder:printcolumn:name="Memory",type=string,JSONPath=`.spec.memory`
//+kubebuilder:storageversion

type VirtualMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineSpec   `json:"spec,omitempty"`
	Status VirtualMachineStatus `json:"status,omitempty"`
}

type VirtualMachineSpec struct {
	// CPU specifies the number of vCPUs.
	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+m)$`
	CPU string `json:"cpu"`

	// Memory specifies memory in Kubernetes quantity format.
	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+)(Mi|Gi)$`
	Memory string `json:"memory"`

	// CPUModel exposes KubeVirt CPU model selection.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^$|[A-Za-z0-9_.-]+$`
	CPUModel string `json:"cpuModel,omitempty"`

	// DedicatedCPUPlacement pins vCPUs to pCPUs when enabled.
	// +kubebuilder:default=false
	DedicatedCPUPlacement bool `json:"dedicatedCPUPlacement,omitempty"`

	// NUMAPolicy configures NUMA scheduling.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=none;strict;best-effort
	// +kubebuilder:default="none"
	NUMAPolicy string `json:"numaPolicy,omitempty"`

	// NUMA allows fine-grained NUMA cell configuration.
	// +kubebuilder:validation:Optional
	NUMA *NUMASpec `json:"numa,omitempty"`

	// Hugepages configures hugepage size.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="";"1Gi";"2Mi"
	// +kubebuilder:default=""
	Hugepages string `json:"hugepages,omitempty"`

	// GPUs assigns virtual GPUs or PCI passthrough devices.
	// +kubebuilder:validation:Optional
	GPUs []GPUDevice `json:"gpus,omitempty"`

	// Disks defines data/system disks.
	// +kubebuilder:validation:MinItems=1
	Disks []VirtualMachineDisk `json:"disks"`

	// Nets defines attached networks.
	// +kubebuilder:validation:MinItems=1
	Nets []VirtualMachineNetwork `json:"nets"`

	// CloudInit configures guest customization via cloud-init.
	// +kubebuilder:validation:Optional
	CloudInit *CloudInitSpec `json:"cloudInit,omitempty"`

	// Console exposes supported consoles for UI integrations.
	// +kubebuilder:validation:Optional
	Console *ConsoleDevices `json:"console,omitempty"`

	// LiveMigration declares migration preferences.
	// +kubebuilder:validation:Optional
	LiveMigration *LiveMigrationSpec `json:"liveMigration,omitempty"`

	// LivenessProbe checks VM health.
	// +kubebuilder:validation:Optional
	LivenessProbe *Probe `json:"livenessProbe,omitempty"`

	// ReadinessProbe checks service readiness.
	// +kubebuilder:validation:Optional
	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`

	// DesiredPowerState indicates preferred running state.
	// +kubebuilder:default=Running
	// +kubebuilder:validation:Enum=Running;Stopped
	PowerState PowerState `json:"powerState,omitempty"`

	// KubeVirt extends the spec with raw KubeVirt VM fragments.
	// +kubebuilder:validation:Optional
	KubeVirt *runtime.RawExtension `json:"kubeVirt,omitempty"`
}

type Probe struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=10
	PeriodSeconds int32 `json:"periodSeconds,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=30
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty"`
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=3
	FailureThreshold int32 `json:"failureThreshold,omitempty"`
}

type VirtualMachineDisk struct {
	// +kubebuilder:validation:Enum=system;data;ephemeral
	Type string `json:"type"`

	// +kubebuilder:validation:Enum=virtio;sata;scsi
	// +kubebuilder:default=virtio
	Bus string `json:"bus,omitempty"`

	// +kubebuilder:validation:Enum=none;writeback;writethrough;directsync
	// +kubebuilder:default=none
	Cache string `json:"cache,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	IOThread bool `json:"iothread,omitempty"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	BootOrder *int32 `json:"bootOrder,omitempty"`

	// +kubebuilder:default=false
	Hotplug bool `json:"hotplug,omitempty"`

	// Source references disk resource
	// +kubebuilder:validation:Required
	DiskRef LocalObjectReference `json:"diskRef"`
}

type LocalObjectReference struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

type VirtualMachineNetwork struct {
	// +kubebuilder:validation:Enum=bridge;masquerade;sriov
	Type string `json:"type"`

	// NAD reference when using bridge/sriov.
	// +kubebuilder:validation:Optional
	NADRef *NamespacedName `json:"nadRef,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^$|[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Model string `json:"model,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^$|([0-9]+)(K|M|G)?bps$`
	Bandwidth string `json:"bandwidth,omitempty"`

	// SRIOVResource is required when type=sriov.
	// +kubebuilder:validation:Optional
	SRIOVResource string `json:"sriovResource,omitempty"`

	// +kubebuilder:default=false
	Multiqueue bool `json:"multiqueue,omitempty"`
}

type NamespacedName struct {
	// +kubebuilder:validation:MinLength=1
	Namespace string `json:"namespace"`
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

type VirtualMachineStatus struct {
	// Conditions represents the latest available observations.
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// PowerState represents observed state.
	// +kubebuilder:validation:Optional
	PowerState PowerState `json:"powerState,omitempty"`

	// Phase string for compatibility
	Phase string `json:"phase,omitempty"`

	// MigrationState reflects the current live migration progress.
	MigrationState string `json:"migrationState,omitempty"`
}

// NUMASpec describes virtual NUMA topology.
type NUMASpec struct {
	// +kubebuilder:validation:MinItems=1
	Cells []NUMACell `json:"cells"`
}

// NUMACell represents a NUMA cell definition.
type NUMACell struct {
	// +kubebuilder:validation:Minimum=0
	ID int32 `json:"id"`

	// cpus expressed as a CPU set, e.g. "0-3".
	// +kubebuilder:validation:Pattern=`^[0-9]+(-[0-9]+)?(,[0-9]+(-[0-9]+)?)*$`
	CPUs string `json:"cpus"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi)$`
	Memory string `json:"memory"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Optional
	ThreadsPerCore *int32 `json:"threadsPerCore,omitempty"`
}

// GPUDevice declares vGPU or passthrough attachment.
type GPUDevice struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=vgpu;passthrough
	DeviceType string `json:"deviceType"`

	// +kubebuilder:validation:Pattern=`^[a-z0-9.-]+/[a-z0-9.-]+$`
	ResourceName string `json:"resourceName"`
}

// CloudInitSpec configures guest personalization.
type CloudInitSpec struct {
	// +kubebuilder:validation:Optional
	UserData string `json:"userData,omitempty"`

	// +kubebuilder:validation:Optional
	NetworkData string `json:"networkData,omitempty"`

	// +kubebuilder:validation:Optional
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`

	// +kubebuilder:validation:Optional
	UserDataSecretRef *NamespacedName `json:"userDataSecretRef,omitempty"`

	// +kubebuilder:validation:Optional
	NetworkDataSecretRef *NamespacedName `json:"networkDataSecretRef,omitempty"`
}

// ConsoleDevices toggles remote consoles.
type ConsoleDevices struct {
	// +kubebuilder:default=true
	VNC bool `json:"vnc,omitempty"`

	// +kubebuilder:default=true
	Serial bool `json:"serial,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum="";"spice"
	Type string `json:"type,omitempty"`
}

// LiveMigrationSpec governs VM migrations.
type LiveMigrationSpec struct {
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^$|([1-9][0-9]*)(Mi|Gi)$`
	Bandwidth string `json:"bandwidth,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	CompletionTimeoutSeconds *int32 `json:"completionTimeoutSeconds,omitempty"`

	// +kubebuilder:default=false
	AllowPostCopy bool `json:"allowPostCopy,omitempty"`

	// +kubebuilder:default=false
	AutoConverge bool `json:"autoConverge,omitempty"`
}

//+kubebuilder:object:root=true

type VirtualMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualMachine{}, &VirtualMachineList{})
}
