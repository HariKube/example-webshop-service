/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RegistryTokenSpec defines the desired state of RegistryToken.
type RegistryTokenSpec struct {
	// +kubebuilder:validation:Required
	// OwnerRef represents the owner reference of the order.
	OwnerRef metav1.OwnerReference `json:"ownerRef"`

	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the addon.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the addon.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// User represents the associated user information.
	User User `json:"user"`

	// +kubebuilder:validation:Required
	// ExpirationTimestamp represents the expiration time of the registry token.
	ExpirationTimestamp metav1.Time `json:"expirationTimestamp"`
}

// RegistryTokenStatus defines the observed state of RegistryToken.
type RegistryTokenStatus struct {
	ErrorMessage   string                       `json:"errorMessage,omitempty"`
	ErrorTimestamp metav1.Time                  `json:"errorTimestamp,omitempty"`
	TokenRef       *corev1.LocalObjectReference `json:"tokenRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Owner",type="string",JSONPath=".spec.ownerRef.name"
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Expiring",type="string",JSONPath=".spec.expirationTimestamp"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage"
// +kubebuilder:selectablefield:JSONPath=".spec.ownerRef.name"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"

// RegistryToken is the Schema for the registrytokens API.
type RegistryToken struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistryTokenSpec   `json:"spec,omitempty"`
	Status RegistryTokenStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistryTokenList contains a list of RegistryToken.
type RegistryTokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RegistryToken `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RegistryToken{}, &RegistryTokenList{})
}
