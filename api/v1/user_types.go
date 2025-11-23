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

// UserSpec defines the desired state of User.
type UserSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=125
	// +kubebuilder:validation:Pattern=`^[\p{L}][\p{L}\p{M}\s'\-]*$`
	// FirstName represents the first name of the user.
	FirstName string `json:"firstName"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=125
	// +kubebuilder:validation:Pattern=`^[\p{L}][\p{L}\p{M}\s'\-]*$`
	// LastName represents the last name of the user.
	LastName string `json:"lastName"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=email
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	// Email represents the email address of the user.
	Email string `json:"email"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:Pattern=`^\+?[1-9]\d{1,14}$`
	// PhoneNumber represents the phone number of the user.
	PhoneNumber *string `json:"phoneNumber,omitempty"`
}

// UserStatus defines the observed state of User.
type UserStatus struct {
	LastGeneration int64 `json:"lastGeneration,omitempty"`
	// +kubebuilder:validation:Enum=Pending;Validated
	// +kubebuilder:default=Pending
	Phase       string                       `json:"phase,omitempty"`
	PasswordRef *corev1.LocalObjectReference `json:"passwordRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FirstName",type=string,JSONPath=`.spec.firstName`
// +kubebuilder:printcolumn:name="LastName",type=string,JSONPath=`.spec.lastName`
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:selectablefield:JSONPath=".spec.firstName"
// +kubebuilder:selectablefield:JSONPath=".spec.lastName"
// +kubebuilder:selectablefield:JSONPath=".spec.email"

// User is the Schema for the users API.
type User struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserSpec   `json:"spec,omitempty"`
	Status UserStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UserList contains a list of User.
type UserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []User `json:"items"`
}

func init() {
	SchemeBuilder.Register(&User{}, &UserList{})
}
