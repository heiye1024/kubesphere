package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// PowerState defines desired VM power state.
// +kubebuilder:validation:Enum=Running;Stopped
type PowerState string

const (
	PowerStateRunning PowerState = "Running"
	PowerStateStopped PowerState = "Stopped"
)

const (
	ConditionTypePowerOperation = "PowerOperation"
	ConditionTypeMigration      = "Migration"
)

const (
	ConditionReasonProgressing = "Progressing"
	ConditionReasonReady       = "Ready"
	ConditionReasonFailed      = "Failed"
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
	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+m)$`
	CPU string `json:"cpu"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+)(Mi|Gi)$`
	Memory string `json:"memory"`

	// +kubebuilder:validation:Pattern=`^$|[A-Za-z0-9_.-]+$`
	CPUModel string `json:"cpuModel,omitempty"`

	// +kubebuilder:default=false
	DedicatedCPUPlacement bool `json:"dedicatedCPUPlacement,omitempty"`

	// +kubebuilder:validation:Enum=none;strict;best-effort
	// +kubebuilder:default="none"
	NUMAPolicy string `json:"numaPolicy,omitempty"`

	NUMA *NUMASpec `json:"numa,omitempty"`

	// +kubebuilder:validation:Enum="";"1Gi";"2Mi"
	Hugepages string `json:"hugepages,omitempty"`

	GPUs []GPUDevice `json:"gpus,omitempty"`

	// +kubebuilder:validation:MinItems=1
	Disks []VirtualMachineDisk `json:"disks"`

	// +kubebuilder:validation:MinItems=1
	Nets []VirtualMachineNetwork `json:"nets"`

	CloudInit *CloudInitSpec `json:"cloudInit,omitempty"`

	Console *ConsoleDevices `json:"console,omitempty"`

	LiveMigration *LiveMigrationSpec `json:"liveMigration,omitempty"`

	LivenessProbe *Probe `json:"livenessProbe,omitempty"`

	ReadinessProbe *Probe `json:"readinessProbe,omitempty"`

	// +kubebuilder:validation:Enum=Running;Stopped
	// +kubebuilder:default=Running
	PowerState PowerState `json:"powerState,omitempty"`

	KubeVirt *runtime.RawExtension `json:"kubeVirt,omitempty"`
}

type Probe struct {
	PeriodSeconds    int32 `json:"periodSeconds,omitempty"`
	TimeoutSeconds   int32 `json:"timeoutSeconds,omitempty"`
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

	// +kubebuilder:default=false
	IOThread bool `json:"iothread,omitempty"`

	// +kubebuilder:validation:Minimum=1
	BootOrder *int32 `json:"bootOrder,omitempty"`

	// +kubebuilder:default=false
	Hotplug bool `json:"hotplug,omitempty"`

	DiskRef LocalObjectReference `json:"diskRef"`
}

type LocalObjectReference struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
}

type VirtualMachineNetwork struct {
	// +kubebuilder:validation:Enum=bridge;masquerade;sriov
	Type string `json:"type"`

	NADRef *NamespacedName `json:"nadRef,omitempty"`

	// +kubebuilder:validation:Pattern=`^$|[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Model string `json:"model,omitempty"`

	// +kubebuilder:validation:Pattern=`^$|([0-9]+)(K|M|G)?bps$`
	Bandwidth string `json:"bandwidth,omitempty"`

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
	Conditions     []metav1.Condition `json:"conditions,omitempty"`
	PowerState     PowerState         `json:"powerState,omitempty"`
	Phase          string             `json:"phase,omitempty"`
	MigrationState string             `json:"migrationState,omitempty"`
}

type NUMASpec struct {
	// +kubebuilder:validation:MinItems=1
	Cells []NUMACell `json:"cells"`
}

type NUMACell struct {
	// +kubebuilder:validation:Minimum=0
	ID int32 `json:"id"`

	// +kubebuilder:validation:Pattern=`^[0-9]+(-[0-9]+)?(,[0-9]+(-[0-9]+)?)*$`
	CPUs string `json:"cpus"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi)$`
	Memory string `json:"memory"`

	// +kubebuilder:validation:Minimum=1
	ThreadsPerCore *int32 `json:"threadsPerCore,omitempty"`
}

type GPUDevice struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// +kubebuilder:validation:Enum=vgpu;passthrough
	DeviceType string `json:"deviceType"`

	// +kubebuilder:validation:Pattern=`^[a-z0-9.-]+/[a-z0-9.-]+$`
	ResourceName string `json:"resourceName"`
}

type CloudInitSpec struct {
	UserData string `json:"userData,omitempty"`

	NetworkData string `json:"networkData,omitempty"`

	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`

	UserDataSecretRef *NamespacedName `json:"userDataSecretRef,omitempty"`

	NetworkDataSecretRef *NamespacedName `json:"networkDataSecretRef,omitempty"`
}

type ConsoleDevices struct {
	// +kubebuilder:default=true
	VNC bool `json:"vnc,omitempty"`

	// +kubebuilder:default=true
	Serial bool `json:"serial,omitempty"`

	// +kubebuilder:validation:Enum="";"spice"
	Type string `json:"type,omitempty"`
}

type LiveMigrationSpec struct {
	// +kubebuilder:default=false
	Enabled bool `json:"enabled,omitempty"`

	// +kubebuilder:validation:Pattern=`^$|([1-9][0-9]*)(Mi|Gi)$`
	Bandwidth string `json:"bandwidth,omitempty"`

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
