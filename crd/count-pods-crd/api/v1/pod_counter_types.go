package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodCounterSpec defines the desired state
type PodCounterSpec struct {
	Interval int32 `json:"interval,omitempty"`
}

// PodCounterStatus defines the observed state
type PodCounterStatus struct {
	MonitoredNamespaces []string `json:"monitoredNamespaces,omitempty"`
	LastChecked         string   `json:"lastChecked,omitempty"`
	PodCounts           map[string]int32 `json:"podCounts,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type PodCounter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PodCounterSpec   `json:"spec,omitempty"`
	Status            PodCounterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type PodCounterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodCounter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodCounter{}, &PodCounterList{})
}
