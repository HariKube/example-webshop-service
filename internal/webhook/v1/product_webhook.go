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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// log is for logging in this package.
var productlog = logf.Log.WithName("product-resource")

// SetupProductWebhookWithManager registers the webhook for Product in the manager.
func SetupProductWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&productv1.Product{}).
		WithValidator(&ProductCustomValidator{}).
		WithDefaulter(&ProductCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-product-webshop-harikube-info-v1-product,mutating=true,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=products,verbs=create;update,versions=v1,name=mproduct-v1.kb.io,admissionReviewVersions=v1

// ProductCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Product when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ProductCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ProductCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Product.
func (d *ProductCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	product, ok := obj.(*productv1.Product)

	if !ok {
		return fmt.Errorf("expected an Product object but got %T", obj)
	}
	productlog.Info("Defaulting for Product", "name", product.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-product-webshop-harikube-info-v1-product,mutating=false,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=products,verbs=create;update,versions=v1,name=vproduct-v1.kb.io,admissionReviewVersions=v1

// ProductCustomValidator struct is responsible for validating the Product resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ProductCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ProductCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Product.
func (v *ProductCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	product, ok := obj.(*productv1.Product)
	if !ok {
		return nil, fmt.Errorf("expected a Product object but got %T", obj)
	}
	productlog.Info("Validation for Product upon creation", "name", product.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Product.
func (v *ProductCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	product, ok := newObj.(*productv1.Product)
	if !ok {
		return nil, fmt.Errorf("expected a Product object for the newObj but got %T", newObj)
	}
	productlog.Info("Validation for Product upon update", "name", product.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Product.
func (v *ProductCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	product, ok := obj.(*productv1.Product)
	if !ok {
		return nil, fmt.Errorf("expected a Product object but got %T", obj)
	}
	productlog.Info("Validation for Product upon deletion", "name", product.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
