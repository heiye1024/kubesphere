package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=vmtemp;vmtemps,scope=Namespaced
//+kubebuilder:subresource:status

type VMTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VMTemplateSpec   `json:"spec,omitempty"`
	Status VMTemplateStatus `json:"status,omitempty"`
}

type VMTemplateSpec struct {
	Parameters  TemplateParameters  `json:"parameters"`
	Constraints TemplateConstraints `json:"constraints,omitempty"`
	UIHints     map[string]string   `json:"uiHints,omitempty"`
}

type TemplateParameters struct {
	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+m)$`
	CPU string `json:"cpu"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi)$`
	Memory string `json:"memory"`

	// +kubebuilder:validation:MinLength=1
	OS string `json:"os"`

	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	Networks []string       `json:"networks"`
	Disks    []TemplateDisk `json:"disks"`
}

type TemplateDisk struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi|Ti)$`
	Size string `json:"size"`

	// +kubebuilder:validation:Enum=system;data;ephemeral
	Type string `json:"type"`
}

type TemplateConstraints struct {
	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+m)$`
	MinCPU string `json:"minCPU,omitempty"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*|[0-9]+m)$`
	MaxCPU string `json:"maxCPU,omitempty"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi)$`
	MinMemory string `json:"minMemory,omitempty"`

	// +kubebuilder:validation:Pattern=`^([1-9][0-9]*)(Mi|Gi)$`
	MaxMemory string `json:"maxMemory,omitempty"`
}

type VMTemplateStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true

type VMTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VMTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VMTemplate{}, &VMTemplateList{})
}
