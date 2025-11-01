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

// EmailTemplateSpec defines the desired state of EmailTemplate.
type EmailTemplateSpec struct {
	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the template.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the template.
	Description string `json:"description,omitempty"`

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

// EmailTemplateStatus defines the observed state of EmailTemplate.
type EmailTemplateStatus struct {
	LastGeneration int64 `json:"lastGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"

// EmailTemplate is the Schema for the emailtemplates API.
type EmailTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailTemplateSpec   `json:"spec,omitempty"`
	Status EmailTemplateStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailTemplateList contains a list of EmailTemplate.
type EmailTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmailTemplate{}, &EmailTemplateList{})
}
