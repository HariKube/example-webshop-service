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

// PaymentSpec defines the desired state of Payment.
type PaymentSpec struct {
	// +kubebuilder:validation:Required
	// Price represents the price of the payment in cents.
	Price int64 `json:"price"`
}

// PaymentStatus defines the observed state of Payment.
type PaymentStatus struct {
	LastGeneration   int64       `json:"lastGeneration,omitempty"`
	ErrorMessage     string      `json:"errorMessage,omitempty"`
	ErrorTimestamp   metav1.Time `json:"errorTimestamp,omitempty"`
	PaymentTimestamp metav1.Time `json:"paymentTimestamp,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Price",type="number",JSONPath=".spec.price"
// +kubebuilder:printcolumn:name="Date",type="date",JSONPath=".status.paymentTimestamp"

// Payment is the Schema for the payments API.
type Payment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PaymentSpec   `json:"spec,omitempty"`
	Status PaymentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PaymentList contains a list of Payment.
type PaymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Payment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Payment{}, &PaymentList{})
}
