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
	authorizationv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// +kubebuilder:rbac:groups="",resources=secrets;serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

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

	rules := []authorizationv1.PolicyRule{}
	for kind, verbs := range map[string][]string{
		"orders":         {"get", "list", "watch"},
		"registrytokens": {"get", "list", "watch", "create", "delete"},
		"tenants":        {"get", "list", "watch", "update", "patch"},
	} {
		rules = append(rules, authorizationv1.PolicyRule{
			APIGroups: []string{"product.webshop.harikube.info"},
			Resources: []string{kind},
			Verbs:     verbs,
		})
	}
	rules = append(rules, authorizationv1.PolicyRule{
		APIGroups:     []string{"product.webshop.harikube.info"},
		Resources:     []string{"users"},
		ResourceNames: []string{user.Name},
		Verbs:         []string{"get", "list", "watch", "update", "patch"},
	})

	role := authorizationv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      string(user.UID),
			Namespace: user.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: user.APIVersion,
					Kind:       user.Kind,
					Name:       user.Name,
					UID:        user.UID,
				},
			},
		},
		Rules: rules,
	}
	if err := r.Create(ctx, &role); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "Role creation failed")
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, types.NamespacedName{
			Name:      role.Name,
			Namespace: role.Namespace,
		}, &role); err != nil {
			logger.Error(err, "Role fetch failed")
			return ctrl.Result{}, err
		}

		role.Rules = rules
		if err := r.Update(ctx, &role); err != nil {
			logger.Error(err, "Role update failed")
			return ctrl.Result{}, err
		}

		logger.Info("Role has been updated", "roleName", role.Name)
	} else {
		logger.Info("Role has been created", "roleName", role.Name)
	}

	roleBinding := authorizationv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      string(user.UID),
			Namespace: user.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: user.APIVersion,
					Kind:       user.Kind,
					Name:       user.Name,
					UID:        user.UID,
				},
			},
		},
		Subjects: []authorizationv1.Subject{
			{
				Kind:      authorizationv1.ServiceAccountKind,
				Name:      string(user.UID),
				Namespace: user.Namespace,
			},
		},
		RoleRef: authorizationv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role",
			Name:     string(user.UID),
		},
	}
	if err := r.Create(ctx, &roleBinding); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "RoleBinding creation failed")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("RoleBinding has been created", "roleBindingName", roleBinding.Name)
	}

	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      string(user.UID),
			Namespace: user.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: user.APIVersion,
					Kind:       user.Kind,
					Name:       user.Name,
					UID:        user.UID,
				},
			},
		},
	}
	if err := r.Create(ctx, &serviceAccount); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			logger.Error(err, "ServiceAccount creation failed")
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("ServiceAccount has been created", "serviceAccountName", serviceAccount.Name)
	}

	if hash, ok := user.Annotations["product.webshop.harikube.info/password"]; ok {
		password := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      user.Name,
				Namespace: user.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					{
						APIVersion: user.APIVersion,
						Kind:       user.Kind,
						Name:       user.Name,
						UID:        user.UID,
					},
				},
			},
			StringData: map[string]string{
				"hash": hash,
			},
		}
		if err := r.Create(ctx, &password); err != nil {
			if !apierrors.IsAlreadyExists(err) {
				logger.Error(err, "Secret creation failed")
				return ctrl.Result{}, err
			}

			if err := r.Get(ctx, types.NamespacedName{
				Name:      password.Name,
				Namespace: password.Namespace,
			}, &password); err != nil {
				logger.Error(err, "Secret fetch failed")
				return ctrl.Result{}, err
			}

			password.StringData = map[string]string{
				"hash": hash,
			}
			if err := r.Update(ctx, &password); err != nil {
				logger.Error(err, "Secret update failed")
				return ctrl.Result{}, err
			}

			logger.Info("Secret has been updated", "secretName", password.Name)
		} else {
			logger.Info("Secret has been created", "secretName", password.Name)
		}

		delete(user.Annotations, "product.webshop.harikube.info/password")
		user.Status.PasswordRef = &corev1.LocalObjectReference{
			Name: password.Name,
		}
	}

	patchedUser := user.DeepCopy()
	patchedUser.Status.LastGeneration = user.Generation
	if err := r.Status().Patch(ctx, patchedUser, client.MergeFrom(&user)); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		logger.Error(err, "User status update failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&productv1.User{}).
		Named("user").
		Complete(r)
}
