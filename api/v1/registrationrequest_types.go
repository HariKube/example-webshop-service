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

// RegistrationRequestSpec defines the desired state of RegistrationRequest.
type RegistrationRequestSpec struct {
	// +kubebuilder:validation:Required
	// User contains the user information for the registration request.
	User UserSpec `json:"user"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=8
	// Password represents the password for the user.
	Password string `json:"password"`
	// +kubebuilder:validation:Required
	// Tenant contains the tenant information for the registration request.
	Tenant TenantSpec `json:"tenant"`
}

// RegistrationRequestStatus defines the observed state of RegistrationRequest.
type RegistrationRequestStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// RegistrationRequest is the Schema for the registrationrequests API.
type RegistrationRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegistrationRequestSpec   `json:"spec,omitempty"`
	Status RegistrationRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RegistrationRequestList contains a list of RegistrationRequest.
type RegistrationRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RegistrationRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RegistrationRequest{}, &RegistrationRequestList{})
}
