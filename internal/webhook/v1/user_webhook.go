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

// nolint:unused
// log is for logging in this package.
var userlog = logf.Log.WithName("user-resource")

// SetupUserWebhookWithManager registers the webhook for User in the manager.
func SetupUserWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&productv1.User{}).
		WithValidator(&UserCustomValidator{}).
		WithDefaulter(&UserCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-product-webshop-harikube-info-v1-user,mutating=true,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=users,verbs=create;update,versions=v1,name=muser-v1.kb.io,admissionReviewVersions=v1

// UserCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind User when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type UserCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &UserCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind User.
func (d *UserCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	user, ok := obj.(*productv1.User)

	if !ok {
		return fmt.Errorf("expected an User object but got %T", obj)
	}
	userlog.Info("Defaulting for User", "name", user.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-product-webshop-harikube-info-v1-user,mutating=false,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=users,verbs=create;update,versions=v1,name=vuser-v1.kb.io,admissionReviewVersions=v1

// UserCustomValidator struct is responsible for validating the User resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type UserCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &UserCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	user, ok := obj.(*productv1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object but got %T", obj)
	}
	userlog.Info("Validation for User upon creation", "name", user.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	user, ok := newObj.(*productv1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object for the newObj but got %T", newObj)
	}
	userlog.Info("Validation for User upon update", "name", user.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	user, ok := obj.(*productv1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object but got %T", obj)
	}
	userlog.Info("Validation for User upon deletion", "name", user.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
