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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

// TenantReconciler reconciles a Tenant object
type TenantReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=tenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=product.webshop.harikube.info,resources=tenants/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Tenant object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx).WithValues("controller", "tenant", "name", req.NamespacedName)

	tenant := productv1.Tenant{}
	if err := r.Get(ctx, req.NamespacedName, &tenant); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "Tenant fetch failed")
		return ctrl.Result{}, err
	}
	tenant.GetObjectKind().SetGroupVersionKind(productv1.GroupVersion.WithKind("Tenant"))

	if tenant.DeletionTimestamp != nil || !tenant.DeletionTimestamp.IsZero() {
		logger.Info("Tenant deleted")

		namespace := corev1.Namespace{}
		if err := r.Get(ctx, types.NamespacedName{
			Name: string(tenant.UID),
		}, &namespace); err != nil {
			if !apierrors.IsNotFound(err) {
				logger.Error(err, "Namespace fetch failed", "namespaceName", namespace.Name)
				return ctrl.Result{}, err
			}
		}
		if controllerutil.ContainsFinalizer(&namespace, "product.webshop.harikube.info/tenant") {
			controllerutil.RemoveFinalizer(&namespace, "product.webshop.harikube.info/tenant")
			if err := r.Update(ctx, &namespace); err != nil {
				if apierrors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}

				logger.Error(err, "Tenant finalizer removal failed", "namespaceName", namespace.Name)
				return ctrl.Result{}, err
			}
		}

		if controllerutil.ContainsFinalizer(&tenant, "product.webshop.harikube.info/tenant") {
			controllerutil.RemoveFinalizer(&tenant, "product.webshop.harikube.info/tenant")
			if err := r.Update(ctx, &tenant); err != nil {
				if apierrors.IsNotFound(err) {
					return ctrl.Result{}, nil
				}

				logger.Error(err, "Tenant finalizer removal failed")
				return ctrl.Result{}, err
			}
		}

		return ctrl.Result{}, nil
	} else if tenant.Generation == 1 && tenant.Status.LastGeneration == 0 {
		logger.Info("Tenant created")

		namespace := corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: tenant.Name,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion:         tenant.APIVersion,
						Kind:               tenant.Kind,
						Name:               tenant.Name,
						UID:                tenant.UID,
						BlockOwnerDeletion: ptr.To(true),
					},
				},
				Finalizers: []string{
					"product.webshop.harikube.info/tenant",
				},
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
	} else {
		if tenant.Status.LastGeneration == tenant.Generation {
			return ctrl.Result{}, nil
		}

		logger.Info("User updated")
	}

	patchedTenant := tenant.DeepCopy()
	patchedTenant.Status.LastGeneration = tenant.Generation
	if err := r.Status().Patch(ctx, patchedTenant, client.MergeFrom(&tenant)); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "Tenant status update failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TenantReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&productv1.Tenant{}).
		Named("tenant").
		Complete(r)
}
