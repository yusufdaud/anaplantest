package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// YusufSpec defines the desired state of Yusuf
type YusufSpec struct {
	Name string `json:"name"`
}

// YusufStatus defines the observed state of Yusuf
type YusufStatus struct {
	Happy bool `json:"happy,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Yusuf is the Schema for the yusufs API
type Yusuf struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   YusufSpec   `json:"spec,omitempty"`
	Status YusufStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// YusufList contains a list of Yusuf
type YusufList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Yusuf `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Yusuf{}, &YusufList{})
}
