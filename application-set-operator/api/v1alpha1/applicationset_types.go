package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ApplicationSetConfig Configurable properties for the ApplicationSet
type ApplicationSetConfig struct {
	// +kubebuilder:default=2
	// +kubebuilder:validation:Optional
	Replicas int32 `json:"replicas"`

	// +kubebuilder:default=80
	// +kubebuilder:validation:Optional
	Port int32 `json:"port"`

	Image string `json:"image"`
}

// ApplicationSetSpec defines the desired state of ApplicationSet
type ApplicationSetSpec struct {
	Frontend ApplicationSetConfig `json:"frontend"`
	Backend  ApplicationSetConfig `json:"backend"`
}

// ApplicationSetStatus defines the observed state of ApplicationSet
type ApplicationSetStatus struct {
	Created string `json:"created"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApplicationSet is the Schema for the applicationsets API
type ApplicationSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSetSpec   `json:"spec,omitempty"`
	Status ApplicationSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationSetList contains a list of ApplicationSet
type ApplicationSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApplicationSet{}, &ApplicationSetList{})
}
