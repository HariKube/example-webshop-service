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

func New(dynamicClient dynamic.Interface, scheme *runtime.Scheme, port, certPath, certFile, keyFile, namespace string) *ApiService {
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
							log := apiServiceLog.WithValues("handler", "registrations", "method", r.Method, "path", r.URL.Path)
							log.Info("Received registration request")

							body, err := io.ReadAll(r.Body)
							if err != nil {
								log.Error(err, "Failed to read request body")
								http.Error(w, "failed to read request body: "+err.Error(), http.StatusBadRequest)
								return
							}
							log.Info("Request body read successfully", "size", len(body))

							var req productv1.RegistrationRequest

							ct := strings.ToLower(r.Header.Get("Content-Type"))
							log.Info("Processing request", "contentType", ct)
							switch {
							case strings.Contains(ct, "json") || strings.HasPrefix(strings.TrimSpace(string(body)), "{"):
								if err := json.Unmarshal(body, &req); err != nil {
									log.Error(err, "Failed to decode JSON request")
									http.Error(w, "failed to decode json request: "+err.Error(), http.StatusBadRequest)
									return
								}
								log.Info("Successfully decoded JSON request")
							case strings.Contains(ct, "yaml"), strings.Contains(ct, "x-yaml"):
								fallthrough
							default:
								if err := yaml.Unmarshal(body, &req); err != nil {
									log.Error(err, "Failed to decode YAML request")
									http.Error(w, "failed to decode yaml request: "+err.Error(), http.StatusBadRequest)
									return
								}
								log.Info("Successfully decoded YAML request")
							}

							req.GetObjectKind().SetGroupVersionKind(productv1.GroupVersion.WithKind("RegistrationRequest"))
							log.Info("RegistrationRequest details", "name", req.GetName(), "email", req.Spec.User.Email)

							uriParts := strings.Split(r.RequestURI, "/")
							req.Name = uriParts[len(uriParts)-1]

							objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&req)
							if err != nil {
								log.Error(err, "Failed to convert object to unstructured")
								http.Error(w, "failed to convert object to unstructured: "+err.Error(), http.StatusInternalServerError)
								return
							}

							gvr := schema.GroupVersionResource{Group: productv1.GroupVersion.Group, Version: productv1.GroupVersion.Version, Resource: "registrationrequests"}
							log.Info("Creating RegistrationRequest", "gvr", gvr.String(), "name", req.GetName())

							_, err = dynamicClient.Resource(gvr).Namespace(namespace).Create(r.Context(), &unstructured.Unstructured{Object: objMap}, metav1.CreateOptions{})
							if err != nil {
								if apierrors.IsAlreadyExists(err) {
									log.Info("RegistrationRequest already exists", "name", req.GetName())
									http.Error(w, "resource already exists", http.StatusConflict)
									return
								}

								log.Error(err, "Failed to create RegistrationRequest", "name", req.GetName())
								http.Error(w, "failed to create resource: "+err.Error(), http.StatusInternalServerError)
								return
							}
							log.Info("RegistrationRequest created successfully", "name", req.GetName())

							listOpts := metav1.ListOptions{FieldSelector: "metadata.name=" + req.GetName(), Watch: true}
							log.Info("Starting watch for RegistrationRequest deletion", "name", req.GetName())
							wch, err := dynamicClient.Resource(gvr).Namespace(namespace).Watch(r.Context(), listOpts)
							if err != nil {
								log.Error(err, "Failed to watch registration flow", "name", req.GetName())
								http.Error(w, "failed to watch registration flow", http.StatusInternalServerError)
								return
							}
							defer wch.Stop()

							for ev := range wch.ResultChan() {
								if ev.Object == nil {
									continue
								}

								log.V(1).Info("Received watch event", "type", ev.Type, "name", req.GetName())

								if ev.Type == watch.Deleted {
									log.Info("RegistrationRequest completed (deleted), registration flow finished", "name", req.GetName())
									w.WriteHeader(http.StatusCreated)
									if fl, ok := w.(http.Flusher); ok {
										fl.Flush()
									}
									return
								}
							}
							log.Info("Watch channel closed", "name", req.GetName())
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
							log := apiServiceLog.WithValues("handler", "login", "method", r.Method, "path", r.URL.Path)
							log.Info("Login endpoint called")
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusOK)
							log.Info("Login response sent", "status", http.StatusOK)
						},
						"/verify": func(w http.ResponseWriter, r *http.Request) {
							log := apiServiceLog.WithValues("handler", "verify", "method", r.Method, "path", r.URL.Path)
							log.Info("Verify endpoint called")
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusOK)
							log.Info("Verify response sent", "status", http.StatusOK)
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
