package v1

import corev1 "k8s.io/api/core/v1"

type RemoteObjectReference struct {
	corev1.LocalObjectReference `json:",inline"`

	// +kubebuilder:validation:Required
	// Namespace represents the namespace of the referenced object.
	Namespace string `json:"namespace"`
}
