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

// CouponSpec defines the desired state of Coupon.
type CouponSpec struct {
	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the coupon.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the coupon.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// value represents the value of the coupon.
	Value int64 `json:"value"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=price;percent
	// CouponType represents the type of the coupon.
	CouponType string `json:"couponType"`

	// +kubebuilder:validation:Optional
	// ExpireTimestamp represents the expiration date of the coupon.
	ExpireTimestamp metav1.Time `json:"expireTimestamp,omitempty"`
}

// CouponStatus defines the observed state of Coupon.
type CouponStatus struct {
	LastGeneration int64 `json:"lastGeneration,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Value",type="integer",JSONPath=".spec.value"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.couponType"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"
// +kubebuilder:selectablefield:JSONPath=".spec.couponType"

// Coupon is the Schema for the coupons API.
type Coupon struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CouponSpec   `json:"spec,omitempty"`
	Status CouponStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CouponList contains a list of Coupon.
type CouponList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Coupon `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Coupon{}, &CouponList{})
}
