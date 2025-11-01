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

// LicenceSpec defines the desired state of Licence.
type LicenceSpec struct {
	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the licence.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the licence.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// Addons represents the list of addons associated with this licence.
	Addons []Addon `json:"addons"`

	// +kubebuilder:validation:Required
	// ExpireTimestamp represents the expiration date of the licence.
	ExpireTimestamp metav1.Time `json:"expireTimestamp,omitempty"`
}

// LicenceStatus defines the observed state of Licence.
type LicenceStatus struct {
	LastGeneration int64                        `json:"lastGeneration,omitempty"`
	ErrorMessage   string                       `json:"errorMessage,omitempty"`
	ErrorTimestamp metav1.Time                  `json:"errorTimestamp,omitempty"`
	LicenceRef     *corev1.LocalObjectReference `json:"licenceKey,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Expire",type="date",JSONPath=".spec.expireTimestamp"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"

// Licence is the Schema for the licences API.
type Licence struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LicenceSpec   `json:"spec,omitempty"`
	Status LicenceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LicenceList contains a list of Licence.
type LicenceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Licence `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Licence{}, &LicenceList{})
}
