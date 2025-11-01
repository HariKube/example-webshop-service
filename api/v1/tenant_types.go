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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TenantSpec defines the desired state of Tenant.
type TenantSpec struct {
	// +kubebuilder:validation:Required
	// OwnerRef represents the owner reference of the tenant.
	OwnerRef metav1.OwnerReference `json:"ownerRef"`

	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the addon.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the addon.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// Users represents the associated users information.
	UserRefs []User `json:"userRefs"`
}

// TenantStatus defines the observed state of Tenant.
type TenantStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Owner",type="string",JSONPath=".spec.ownerRef.name"
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:selectablefield:JSONPath=".spec.ownerRef.name"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"

// Tenant is the Schema for the tenants API.
type Tenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TenantSpec   `json:"spec,omitempty"`
	Status TenantStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TenantList contains a list of Tenant.
type TenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Tenant `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Tenant{}, &TenantList{})
}
