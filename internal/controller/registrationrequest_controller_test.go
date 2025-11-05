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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
)

var _ = Describe("RegistrationRequest Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		registrationrequest := &productv1.RegistrationRequest{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind RegistrationRequest")
			err := k8sClient.Get(ctx, typeNamespacedName, registrationrequest)
			if err != nil && apierrors.IsNotFound(err) {
				registrationrequest = &productv1.RegistrationRequest{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: productv1.RegistrationRequestSpec{
						User: productv1.UserSpec{
							FirstName: "First",
							LastName:  "Last",
							Email:     "email@harikube.info",
						},
						Password: "Passwd123!",
						Tenant: productv1.TenantSpec{
							Country:    "HU",
							City:       "Budapest",
							Address:    "Address 123",
							PostalCode: "1234",
						},
					},
				}
				Expect(k8sClient.Create(ctx, registrationrequest)).To(Succeed())
			}

			namespace := &corev1.Namespace{}

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name: string(registrationrequest.UID),
			}, namespace)
			if err != nil && apierrors.IsNotFound(err) {
				namespace = &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: string(registrationrequest.UID),
					},
				}
				Expect(k8sClient.Create(ctx, namespace)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &productv1.RegistrationRequest{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(apierrors.IsNotFound(err)).To(BeTrue())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &RegistrationRequestReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
