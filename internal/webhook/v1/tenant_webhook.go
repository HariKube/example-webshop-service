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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// log is for logging in this package.
var tenantlog = logf.Log.WithName("tenant-resource")

// SetupTenantWebhookWithManager registers the webhook for Tenant in the manager.
func SetupTenantWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&productv1.Tenant{}).
		WithValidator(&TenantCustomValidator{}).
		WithDefaulter(&TenantCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-product-webshop-harikube-info-v1-tenant,mutating=true,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=tenants,verbs=create;update,versions=v1,name=mtenant-v1.kb.io,admissionReviewVersions=v1

// TenantCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Tenant when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type TenantCustomDefaulter struct {
}

var _ webhook.CustomDefaulter = &TenantCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Tenant.
func (d *TenantCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	tenant, ok := obj.(*productv1.Tenant)

	if !ok {
		return fmt.Errorf("expected an Tenant object but got %T", obj)
	}
	tenantlog.Info("Defaulting for Tenant", "name", tenant.GetName())

	if !controllerutil.ContainsFinalizer(tenant, "product.webshop.harikube.info/tenant") {
		controllerutil.AddFinalizer(tenant, "product.webshop.harikube.info/tenant")
		tenantlog.Info("Added finalizer for Tenant", "name", tenant.GetName())
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-product-webshop-harikube-info-v1-tenant,mutating=false,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=tenants,verbs=create;update,versions=v1,name=vtenant-v1.kb.io,admissionReviewVersions=v1

// TenantCustomValidator struct is responsible for validating the Tenant resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type TenantCustomValidator struct {
}

var _ webhook.CustomValidator = &TenantCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Tenant.
func (v *TenantCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	tenant, ok := obj.(*productv1.Tenant)
	if !ok {
		return nil, fmt.Errorf("expected a Tenant object but got %T", obj)
	}
	tenantlog.Info("Validation for Tenant upon create", "name", tenant.GetName())

	if len(tenant.OwnerReferences) == 0 {
		return nil, fmt.Errorf("tenant %s must have an owner reference", tenant.GetName())
	} else if tenant.OwnerReferences[0].Kind != "Namespace" {
		return nil, fmt.Errorf("tenant %s must be owned by a Namespace", tenant.GetName())
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Tenant.
func (v *TenantCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	tenant, ok := newObj.(*productv1.Tenant)
	if !ok {
		return nil, fmt.Errorf("expected a Tenant object for the newObj but got %T", newObj)
	}
	tenantlog.Info("Validation for Tenant upon update", "name", tenant.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Tenant.
func (v *TenantCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	tenant, ok := obj.(*productv1.Tenant)
	if !ok {
		return nil, fmt.Errorf("expected a Tenant object but got %T", obj)
	}
	tenantlog.Info("Validation for Tenant upon deletion", "name", tenant.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
