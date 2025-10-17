package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=vmdisk;vmdisks,scope=Namespaced
//+kubebuilder:subresource:status

type VirtualDisk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualDiskSpec   `json:"spec,omitempty"`
	Status VirtualDiskStatus `json:"status,omitempty"`
}

type VirtualDiskSpec struct {
	// +kubebuilder:validation:Enum=DataVolume;PVC;Blank
	Backing string `json:"backing"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi|Ti)$`
	Size string `json:"size"`

	// +kubebuilder:validation:Enum=ReadWriteOnce;ReadOnlyMany;ReadWriteMany
	// +kubebuilder:default=ReadWriteOnce
	AccessMode string `json:"accessMode"`

	// +kubebuilder:validation:Enum=Block;Filesystem
	// +kubebuilder:default=Filesystem
	VolumeMode string `json:"volumeMode"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	StorageClass string `json:"storageClass,omitempty"`
}

type VirtualDiskStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true

type VirtualDiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualDisk `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualDisk{}, &VirtualDiskList{})
}
