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

// OrderSpec defines the desired state of Order.
type OrderSpec struct {
	// +kubebuilder:validation:Required
	// User represents the order user information.
	User UserSpec `json:"user"`

	// +kubebuilder:validation:Optional
	// BillingUser represents the billing user information.
	BillingUser UserSpec `json:"billingUser,omitempty"`

	// +kubebuilder:validation:Required
	// Products represents the list of products within this order.
	Products []OrderProduct `json:"products"`

	// +kubebuilder:validation:Optional
	// Coupon represents an optional coupon for the order.
	Coupon Coupon `json:"couponCode,omitempty"`

	// +kubebuilder:validation:Required
	// OrderTimestamp represents the date when the order was placed.
	OrderTimestamp metav1.Time `json:"orderTimestamp"`
}

// OrderProduct represents a product within an order.
type OrderProduct struct {
	// +kubebuilder:validation:Required
	// Product represents the product.
	Product Product `json:"product"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default:=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// Quantity represents the quantity of the product ordered.
	Quantity int32 `json:"quantity"`

	// +kubebuilder:validation:Optional
	// Addons represents the list of addons associated with this product.
	Addons []Addon `json:"addons,omitempty"`
}

// OrderStatus defines the observed state of Order.
type OrderStatus struct {
	LastGeneration int64                        `json:"lastGeneration,omitempty"`
	ErrorMessage   string                       `json:"errorMessage,omitempty"`
	ErrorTimestamp metav1.Time                  `json:"errorTimestamp,omitempty"`
	TotalPrice     int64                        `json:"totalPrice,omitempty"`
	PaymentRef     *corev1.LocalObjectReference `json:"paymentRef,omitempty"`
	Licences       []Licence                    `json:"licences,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Date",type="date",JSONPath=".spec.orderTimestamp"
// +kubebuilder:printcolumn:name="User",type="string",JSONPath=".spec.user.spec.email"
// +kubebuilder:printcolumn:name="Error",type="string",JSONPath=".status.errorMessage"
// +kubebuilder:printcolumn:name="Price",type="number",JSONPath=".status.totalPrice"

// Order is the Schema for the orders API.
type Order struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrderSpec   `json:"spec,omitempty"`
	Status OrderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OrderList contains a list of Order.
type OrderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Order `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Order{}, &OrderList{})
}
