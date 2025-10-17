package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=vmnet;vmnets,scope=Namespaced
//+kubebuilder:subresource:status

type VirtualNet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualNetSpec   `json:"spec,omitempty"`
	Status VirtualNetStatus `json:"status,omitempty"`
}

type VirtualNetSpec struct {
	// NADTemplate stores a NetworkAttachmentDefinition manifest.
	// +kubebuilder:validation:MinLength=1
	NADTemplate string `json:"nadTemplate"`

	// BandwidthLimit in Mbps.
	// +kubebuilder:validation:Minimum=1
	BandwidthLimit *int32 `json:"bandwidthLimit,omitempty"`

	// VLAN ID
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4094
	VLAN *int32 `json:"vlan,omitempty"`

	// SRIOVResource indicates resource name.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^$|[a-z0-9.-]+/[a-z0-9.-]+$`
	SRIOVResource string `json:"sriovResource,omitempty"`
}

type VirtualNetStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true

type VirtualNetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualNet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtualNet{}, &VirtualNetList{})
}
