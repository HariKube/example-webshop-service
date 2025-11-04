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

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// RegistrationRequestReconciler reconciles a RegistrationRequest object
type RegistrationRequestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=registrationrequests,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=registrationrequests/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=registrationrequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RegistrationRequest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the registrationrequest.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *RegistrationRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx).WithValues("controller", "registrationrequest", "name", req.NamespacedName)

	request := productv1.RegistrationRequest{}
	if err := r.Get(ctx, req.NamespacedName, &request); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "RegistrationRequest fetch failed")

		return ctrl.Result{}, err
	}

	if request.DeletionTimestamp != nil || !request.DeletionTimestamp.IsZero() {
		logger.Info("RegistrationRequest deleted")

		return ctrl.Result{}, nil
	}

	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: string(request.UID),
		},
	}
	if err := r.Create(ctx, &namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "Namespace creation failed")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Namespace has been created", "namespaceName", namespace.Name)
	}
	if err := r.Get(ctx, types.NamespacedName{
		Name: namespace.Name,
	}, &namespace); err != nil {
		logger.Error(err, "Namespace fetch failed")
		return ctrl.Result{}, err
	}

	tenant := productv1.Tenant{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: namespace.Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion:         namespace.APIVersion,
					Kind:               namespace.Kind,
					Name:               namespace.Name,
					UID:                namespace.UID,
					BlockOwnerDeletion: ptr.To(true),
				},
			},
		},
		Spec: *request.Spec.Tenant.DeepCopy(),
	}
	if err := r.Create(ctx, &tenant); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "Tenant creation failed")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("Tenant has been created", "tenantName", tenant.Name)
	}

	user := productv1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      request.Name,
			Namespace: namespace.Name,
			Annotations: map[string]string{
				"product.webshop.harikube.info/password": request.Spec.Password,
			},
		},
		Spec: *request.Spec.User.DeepCopy(),
	}
	if err := r.Create(ctx, &user); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "User creation failed")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("User has been created", "userName", user.Name)
	}

	if err := r.Delete(ctx, &request); err != nil {
		logger.Error(err, "RegistrationRequest deletion failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RegistrationRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&productv1.RegistrationRequest{}).
		Named("registrationrequest").
		Complete(r)
}
