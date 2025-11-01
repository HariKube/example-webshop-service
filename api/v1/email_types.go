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

// EmailSpec defines the desired state of Email.
type EmailSpec struct {
	// +kubebuilder:validation:Required
	// FromAddress represents the email address of the recipient.
	ToAddress string `json:"toAddress"`

	// +kubebuilder:validation:Required
	// FromName represents the name of the sender.
	FromName string `json:"fromName"`

	// +kubebuilder:validation:Required
	// FromAddress represents the email address of the sender.
	FromAddress string `json:"fromAddress"`

	// +kubebuilder:validation:Required
	// Subject represents the subject of the email.
	Subject string `json:"subject"`

	// +kubebuilder:validation:Required
	// Body represents the body content of the email.
	Body string `json:"body"`
}

// EmailStatus defines the observed state of Email.
type EmailStatus struct {
	LastGeneration int64       `json:"lastGeneration,omitempty"`
	ErrorMessage   string      `json:"errorMessage,omitempty"`
	ErrorTimestamp metav1.Time `json:"errorTimestamp,omitempty"`
	SentTimestamp  metav1.Time `json:"sentTimestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Email",type="string",JSONPath=".spec.toAddress"
// +kubebuilder:printcolumn:name="Sent",type="date",JSONPath=".status.sentTimestamp"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage"
// +kubebuilder:selectablefield:JSONPath=".spec.toAddress"

// Email is the Schema for the emails API.
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailList contains a list of Email.
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
