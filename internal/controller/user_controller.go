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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=users/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=users/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx).WithValues("controller", "user", "name", req.NamespacedName)

	user := productv1.User{}
	if err := r.Get(ctx, req.NamespacedName, &user); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "User fetch failed")

		return ctrl.Result{}, err
	}

	if user.DeletionTimestamp != nil || !user.DeletionTimestamp.IsZero() {
		logger.Info("User deleted")

		return ctrl.Result{}, nil
	} else if user.Generation == 1 && user.Status.LastGeneration == 0 {
		logger.Info("User created")
	} else {
		if user.Status.LastGeneration == user.Generation {
			return ctrl.Result{}, nil
		}

		logger.Info("User updated")
	}

	// if len(user.Status.TenantRefs) == 0 {
	// 	displayName := user.Spec.CompanyName
	// 	if displayName == "" {
	// 		displayName = fmt.Sprintf("%s %s", user.Spec.FirstName, user.Spec.LastName)
	// 	}

	// 	tenant := productv1.Tenant{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name: user.Name,
	// 			OwnerReferences: []metav1.OwnerReference{
	// 				{
	// 					APIVersion:         user.APIVersion,
	// 					Kind:               user.Kind,
	// 					Name:               user.Name,
	// 					UID:                user.UID,
	// 					BlockOwnerDeletion: ptr.To(true),
	// 				},
	// 			},
	// 		},
	// 		Spec: productv1.TenantSpec{
	// 			DisplayName: displayName,
	// 			UserRefs: []productv1.RemoteObjectReference{
	// 				{
	// 					LocalObjectReference: corev1.LocalObjectReference{
	// 						Name: user.Name,
	// 					},
	// 					Namespace: user.Namespace,
	// 				},
	// 			},
	// 		},
	// 	}

	// 	if err := r.Create(ctx, &tenant); err != nil {
	// 		if !apierrors.IsAlreadyExists(err) {
	// 			return ctrl.Result{}, fmt.Errorf("failed to create tenant for user: %w", err)
	// 		}

	// 		return ctrl.Result{}, fmt.Errorf("a Tenant with the name '%s' already exists", user.Name)
	// 	}

	// 	logger.Info("Tenant has been created", "tenantName", tenant.Name)

	// 	patchedUser := user.DeepCopy()
	// 	patchedUser.Status.LastGeneration = user.Generation
	// 	patchedUser.Status.TenantRefs = []corev1.LocalObjectReference{
	// 		{
	// 			Name: tenant.Name,
	// 		},
	// 	}

	// 	if err := r.Status().Patch(ctx, patchedUser, client.MergeFrom(&user)); err != nil {
	// 		if apierrors.IsNotFound(err) {
	// 			return ctrl.Result{}, nil
	// 		}

	// 		logger.Error(err, "User status update failed")

	// 		return ctrl.Result{}, err
	// 	}
	// }

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager, needLeaderElection bool, maxConcurrentReconciles int) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&productv1.User{}).
		Named("user").
		WithOptions(controller.Options{
			NeedLeaderElection:      ptr.To(needLeaderElection),
			MaxConcurrentReconciles: maxConcurrentReconciles,
			RecoverPanic:            ptr.To(true),
			Logger:                  mgr.GetLogger(),
		}).
		Complete(r)
}
