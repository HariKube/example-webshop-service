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
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
	"github.com/alexedwards/argon2id"
	"github.com/dlclark/regexp2"
)

// log is for logging in this package.
var registrationrequestlog = logf.Log.WithName("registrationrequest-resource")

// SetupRegistrationRequestWebhookWithManager registers the webhook for RegistrationRequest in the manager.
func SetupRegistrationRequestWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&productv1.RegistrationRequest{}).
		WithValidator(&RegistrationRequestCustomValidator{
			Client: mgr.GetClient(),
		}).
		WithDefaulter(&RegistrationRequestCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-product-webshop-harikube-info-v1-registrationrequest,mutating=true,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=registrationrequests,verbs=create,versions=v1,name=mregistrationrequest-v1.kb.io,admissionReviewVersions=v1

// RegistrationRequestCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind RegistrationRequest when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type RegistrationRequestCustomDefaulter struct {
}

var _ webhook.CustomDefaulter = &RegistrationRequestCustomDefaulter{}

var (
	passwordMatcher      = regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,64}$`, 0)
	hashingDefaultParams = &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
)

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind RegistrationRequest.
func (d *RegistrationRequestCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	registrationrequest, ok := obj.(*productv1.RegistrationRequest)

	if !ok {
		return fmt.Errorf("expected an RegistrationRequest object but got %T", obj)
	}
	registrationrequestlog.Info("Defaulting for RegistrationRequest", "name", registrationrequest.GetName())

	if ok, _ := passwordMatcher.MatchString(registrationrequest.Spec.Password); !ok {
		return fmt.Errorf("password must include upper, lower, number, special and be 8-64 chars")
	}

	hash, err := argon2id.CreateHash(registrationrequest.Spec.Password, hashingDefaultParams)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	if registrationrequest.Spec.Tenant.CompanyName == "" {
		registrationrequest.Spec.Tenant.CompanyName = fmt.Sprintf("%s %s", registrationrequest.Spec.User.FirstName, registrationrequest.Spec.User.LastName)
	}

	registrationrequest.Spec.Password = hash

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-product-webshop-harikube-info-v1-registrationrequest,mutating=false,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=registrationrequests,verbs=create,versions=v1,name=vregistrationrequest-v1.kb.io,admissionReviewVersions=v1

// RegistrationRequestCustomValidator struct is responsible for validating the RegistrationRequest resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type RegistrationRequestCustomValidator struct {
	client.Client
}

var _ webhook.CustomValidator = &RegistrationRequestCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type RegistrationRequest.
func (v *RegistrationRequestCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	registrationrequest, ok := obj.(*productv1.RegistrationRequest)
	if !ok {
		return nil, fmt.Errorf("expected a RegistrationRequest object but got %T", obj)
	}
	registrationrequestlog.Info("Validation for RegistrationRequest upon create", "name", registrationrequest.GetName())

	existnigUsers := &productv1.UserList{}
	if err := v.List(ctx, existnigUsers, client.MatchingFields{"spec.email": registrationrequest.Spec.User.Email}); err != nil {
		return nil, fmt.Errorf("failed to list existing users: %w", err)
	} else if len(existnigUsers.Items) > 0 {
		return nil, fmt.Errorf("a user with email %s already exists", registrationrequest.Spec.User.Email)
	}

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type RegistrationRequest.
func (v *RegistrationRequestCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	registrationrequest, ok := newObj.(*productv1.RegistrationRequest)
	if !ok {
		return nil, fmt.Errorf("expected a RegistrationRequest object for the newObj but got %T", newObj)
	}
	registrationrequestlog.Info("Validation for RegistrationRequest upon update", "name", registrationrequest.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type RegistrationRequest.
func (v *RegistrationRequestCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	registrationrequest, ok := obj.(*productv1.RegistrationRequest)
	if !ok {
		return nil, fmt.Errorf("expected a RegistrationRequest object but got %T", obj)
	}
	registrationrequestlog.Info("Validation for RegistrationRequest upon deletion", "name", registrationrequest.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
