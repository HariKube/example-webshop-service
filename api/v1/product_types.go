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

// ProductSpec defines the desired state of Product.
type ProductSpec struct {
	// +kubebuilder:validation:Required
	// DisplayName represents the human friendly name of the addon.
	DisplayName string `json:"displayName"`

	// +kubebuilder:validation:Optional
	// Description represents a brief description of the addon.
	Description string `json:"description,omitempty"`

	// +kubebuilder:validation:Required
	// Price represents the price of the product in cents.
	Price int64 `json:"price"`

	// +kubebuilder:validation:Optional
	// Addons represents a list of addons associated with the product.
	Addons []Addon `json:"addons,omitempty"`
}

// ProductStatus defines the observed state of Product.
type ProductStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="Price",type="string",JSONPath=".spec.price"
// +kubebuilder:selectablefield:JSONPath=".spec.displayName"

// Product is the Schema for the products API.
type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductSpec   `json:"spec,omitempty"`
	Status ProductStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductList contains a list of Product.
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Product `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Product{}, &ProductList{})
}
