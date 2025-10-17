package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=vmsnap;vmsnaps,scope=Namespaced
//+kubebuilder:subresource:status

type VMSnapshot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMSnapshotSpec   `json:"spec,omitempty"`
	Status VMSnapshotStatus `json:"status,omitempty"`
}

type VMSnapshotSpec struct {
	// +kubebuilder:validation:Required
	SourceRef NamespacedName `json:"sourceRef"`

	// +kubebuilder:validation:MinItems=0
	IncludedDisks []string `json:"includedDisks,omitempty"`

	// +kubebuilder:validation:Enum=Retain;Delete
	// +kubebuilder:default=Retain
	RetainPolicy string `json:"retainPolicy,omitempty"`
}

type VMSnapshotStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	ReadyToUse bool               `json:"readyToUse,omitempty"`
}

//+kubebuilder:object:root=true

type VMSnapshotList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMSnapshot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMSnapshot{}, &VMSnapshotList{})
}
