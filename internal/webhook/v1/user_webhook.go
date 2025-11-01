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

	"github.com/alexedwards/argon2id"
	"github.com/dlclark/regexp2"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// log is for logging in this package.
var userlog = logf.Log.WithName("user-resource")

// SetupUserWebhookWithManager registers the webhook for User in the manager.
func SetupUserWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&productv1.User{}).
		WithValidator(&UserCustomValidator{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).
		WithDefaulter(&UserCustomDefaulter{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
		}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-product-webshop-harikube-info-v1-user,mutating=true,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=users;users/status,verbs=create;update,versions=v1,name=muser-v1.kb.io,admissionReviewVersions=v1

// UserCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind User when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type UserCustomDefaulter struct {
	client.Client
	Scheme *runtime.Scheme
}

var _ webhook.CustomDefaulter = &UserCustomDefaulter{}

var defaultParams = &argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind User.
func (d *UserCustomDefaulter) Default(ctx context.Context, obj runtime.Object) error {
	user, ok := obj.(*productv1.User)

	if !ok {
		return fmt.Errorf("expected an User object but got %T", obj)
	}

	userlog := userlog.WithValues("user", user.Name)
	userlog.Info("Defaulting for User")

	if !controllerutil.ContainsFinalizer(user, "foregroundDeletion") {
		controllerutil.AddFinalizer(user, "foregroundDeletion")
	}

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-product-webshop-harikube-info-v1-user,mutating=false,failurePolicy=fail,sideEffects=None,groups=product.webshop.harikube.info,resources=users,verbs=create;update;delete,versions=v1,name=vuser-v1.kb.io,admissionReviewVersions=v1

// UserCustomValidator struct is responsible for validating the User resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type UserCustomValidator struct {
	client.Client
	Scheme *runtime.Scheme
}

var _ webhook.CustomValidator = &UserCustomValidator{}

var re = regexp2.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,64}$`, 0)

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type User.
func (v *UserCustomValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	user, ok := obj.(*productv1.User)
	if !ok {
		return nil, fmt.Errorf("expected a User object but got %T", obj)
	}

	if plainPassword, ok := user.Annotations["example-webshop.harikube.info/set-password"]; !ok {
		return nil, fmt.Errorf("missing required annotation 'example-webshop.harikube.info/set-password' to set the user password")
	} else if ok, _ := re.MatchString(plainPassword); !ok {
		return nil, fmt.Errorf("password must include upper, lower, number, special and be 8-64 chars")
	}

	tenant := productv1.Tenant{}
	if err := v.Get(ctx, types.NamespacedName{Name: user.Name}, &tenant); err != nil {
		if !apierrors.IsNotFound(err) {
			return nil, fmt.Errorf("failed to fetch tenant for user: %w", err)
		}
	} else {
		return nil, fmt.Errorf("a Tenant with the name '%s' already exists", user.Name)
	}

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
