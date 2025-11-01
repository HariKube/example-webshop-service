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
	// FirstName represents the first name of the user.
	// Allow basic Unicode letters, spaces, apostrophes, and hyphens.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern=`^[\p{L}][\p{L}\p{M}\s'\-]*$`
	FirstName string `json:"firstName"`

	// +kubebuilder:validation:Required
	// LastName represents the last name of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern=`^[\p{L}][\p{L}\p{M}\s'\-]*$`
	LastName string `json:"lastName"`

	// +kubebuilder:validation:Required
	// Email represents the email address of the user.
	// +kubebuilder:validation:Format=email
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	Email string `json:"email"`

	// +kubebuilder:validation:Optional
	// CompanyName represents the company name of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	CompanyName string `json:"companyName,omitempty"`

	// +kubebuilder:validation:Required
	// Country represents the country of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=3
	Country string `json:"country"`

	// +kubebuilder:validation:Required
	// City represents the city of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=50
	City string `json:"city"`

	// Address represents the address of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	Address string `json:"address"`

	// +kubebuilder:validation:Required
	// PostalCode represents the postal code of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=30
	PostalCode string `json:"postalCode"`

	// +kubebuilder:validation:Optional
	// TaxNumber represents the tax number of the user.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	TaxNumber string `json:"taxNumber,omitempty"`

	// +kubebuilder:validation:Optional
	// PhoneNumber represents the phone number of the user.
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:Pattern=`^\+?[1-9]\d{1,14}$`
	PhoneNumber *string `json:"phoneNumber,omitempty"`

	// +kubebuilder:validation:Optional
	// Tenants represents the list of tenants associated with the user.
	Tenants []corev1.LocalObjectReference `json:"tenants,omitempty"`
}

// UserStatus defines the observed state of User.
type UserStatus struct {
	// +kubebuilder:validation:Enum=Pending;Validated
	// +kubebuilder:default=Pending
	Phase string `json:"phase,omitempty"`

	ValidatedEmail string `json:"validatedEmail,omitempty"`

	PasswordRef *corev1.LocalObjectReference `json:"passwordRef,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="FirstName",type=string,JSONPath=`.spec.firstName`
// +kubebuilder:printcolumn:name="LastName",type=string,JSONPath=`.spec.lastName`
// +kubebuilder:printcolumn:name="Email",type=string,JSONPath=`.spec.email`
// +kubebuilder:printcolumn:name="Company",type=string,JSONPath=`.spec.companyName`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
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
