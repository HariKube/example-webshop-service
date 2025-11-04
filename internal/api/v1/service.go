package v1

import (
	"context"
	"fmt"
	"net/http"
	"os"

	kaf "github.com/HariKube/kubernetes-aggregator-framework/pkg/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"

	productv1 "github.com/HariKube/example-webshop-service/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	apiServiceLog = logf.Log.WithName("api-service")
)

func New(k8sClient client.Client, scheme *runtime.Scheme, port, certPath, certFile, keyFile string) *ApiService {
	sas := ApiService{
		Client: k8sClient,
		Scheme: scheme,
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
							w.Header().Set("Content-Type", "application/json; charset=utf-8")
							w.WriteHeader(http.StatusOK)
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
	client.Client
	Scheme *runtime.Scheme
}

func (s *ApiService) Start(ctx context.Context) (err error) {
	apiServiceLog.Info("Serving api-service server", "host", "", "port", "7443")
	return s.Server.Start(ctx)
}
