package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	kaf "github.com/HariKube/kubernetes-aggregator-framework/pkg/framework"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	yaml "sigs.k8s.io/yaml"
)

var (
	apiServiceLog = logf.Log.WithName("api-service")
)

func New(dynamicClient dynamic.Interface, scheme *runtime.Scheme, port, certPath, certFile, keyFile string) *ApiService {
	sas := ApiService{
		DynamicClient: dynamicClient,
		Scheme:        scheme,
		Server: *kaf.NewServer(kaf.ServerConfig{
			Port:     port,
			CertFile: fmt.Sprintf("%s%c%s", certPath, os.PathSeparator, certFile),
			KeyFile:  fmt.Sprintf("%s%c%s", certPath, os.PathSeparator, keyFile),
			Group:    "api." + productv1.GroupVersion.Group,
			Version:  productv1.GroupVersion.Version,
			APIKinds: []kaf.APIKind{
				{
					ApiResource: metav1.APIResource{
						Name:  "registrations",
						Verbs: []string{"create"},
					},
					CustomResource: &kaf.CustomResource{
						CreateHandler: func(namespace, name string, w http.ResponseWriter, r *http.Request) {
							body, err := io.ReadAll(r.Body)
							if err != nil {
								http.Error(w, "failed to read request body: "+err.Error(), http.StatusBadRequest)
								return
							}

							var req productv1.RegistrationRequest
							ct := strings.ToLower(r.Header.Get("Content-Type"))
							switch {
							case strings.Contains(ct, "json"):
								if err := json.Unmarshal(body, &req); err != nil {
									http.Error(w, "failed to decode json request: "+err.Error(), http.StatusBadRequest)
									return
								}
							case strings.Contains(ct, "yaml"), strings.Contains(ct, "x-yaml"):
								if err := yaml.Unmarshal(body, &req); err != nil {
									http.Error(w, "failed to decode yaml request: "+err.Error(), http.StatusBadRequest)
									return
								}
							default:
								http.Error(w, "invalid content type: "+ct, http.StatusBadRequest)
								return
							}

							req.GetObjectKind().SetGroupVersionKind(productv1.GroupVersion.WithKind("RegistrationRequest"))

							objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&req)
							if err != nil {
								http.Error(w, "failed to convert object to unstructured: "+err.Error(), http.StatusInternalServerError)
								return
							}
							u := &unstructured.Unstructured{Object: objMap}
							u.SetGroupVersionKind(productv1.GroupVersion.WithKind("RegistrationRequest"))

							gvr := schema.GroupVersionResource{Group: productv1.GroupVersion.Group, Version: productv1.GroupVersion.Version, Resource: "registrationrequests"}

							_, err = dynamicClient.Resource(gvr).Namespace("default").Create(r.Context(), u, metav1.CreateOptions{})
							if err != nil {
								if apierrors.IsAlreadyExists(err) {
									http.Error(w, "resource already exists", http.StatusConflict)
									return
								}

								http.Error(w, "failed to create resource: "+err.Error(), http.StatusInternalServerError)
								return
							}

							listOpts := metav1.ListOptions{FieldSelector: "metadata.name=" + req.GetName(), Watch: true}
							wch, err := dynamicClient.Resource(gvr).Namespace("default").Watch(r.Context(), listOpts)
							if err != nil {
								http.Error(w, "failed to watch registration flow", http.StatusInternalServerError)
								return
							}
							defer wch.Stop()

							for ev := range wch.ResultChan() {
								if ev.Object == nil {
									continue
								}
								if ev.Type == watch.Deleted {
									w.WriteHeader(http.StatusCreated)
									if fl, ok := w.(http.Flusher); ok {
										fl.Flush()
									}
									return
								}
							}
						},
					},
				},
				{
					ApiResource: metav1.APIResource{
						Name:  "users",
						Verbs: []string{"get"},
					},
					RawEndpoints: map[string]http.HandlerFunc{
						"/login": func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusOK)
						},
					},
				},
			},
		}),
	}

	return &sas
}

type ApiService struct {
	kaf.Server
	DynamicClient dynamic.Interface
	Scheme        *runtime.Scheme
}

func (s *ApiService) Start(ctx context.Context) (err error) {
	apiServiceLog.Info("Serving api-service server", "host", "", "port", "7443")
	return s.Server.Start(ctx)
}
