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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddonSpec defines the desired state of Addon.
type AddonSpec struct {
	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the addon.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the addon.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// Price represents the price of the addon in cents.
	Price int64 `json:"price"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=product;service;support;onetime
	// AddonType represents the type of the addon.
	AddonType string `json:"addonType"`
}

// AddonStatus defines the observed state of Addon.
type AddonStatus struct {
	LastGeneration int64 `json:"lastGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Price",type="number",JSONPath=".spec.price"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.addonType"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"
// +kubebuilder:selectablefield:JSONPath=".spec.addonType"

// Addon is the Schema for the addons API.
type Addon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddonSpec   `json:"spec,omitempty"`
	Status AddonStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AddonList contains a list of Addon.
type AddonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Addon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Addon{}, &AddonList{})
}
