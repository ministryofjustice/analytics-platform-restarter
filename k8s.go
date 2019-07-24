package main

import (
	"errors"
	"fmt"
	"log"

	appsAPI "k8s.io/api/apps/v1"
	metaAPI "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ErrK8s struct {
	Err error
}

func (e ErrK8s) Error() string {
	return e.Err.Error()
}

type ErrTooManyDeployments struct {
}

// KubernetesClient constructs a new Kubernetes client
func KubernetesClient(path string) k8s.Interface {
	config := loadConfig(path)
	client, err := k8s.NewForConfig(config)
	if err != nil {
		logger.Fatalf("💥 failed to create new client from config: %s", err)
	}

	return client
}

func loadConfig(path string) *rest.Config {
	config, err := rest.InClusterConfig()
	if err == nil {
		return config
	}
	config, err = clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		log.Fatalf("💥 failed to load k8s config from cluster and from kube config (fallback): %s", err)
		return nil
	}

	return config
}

// GetDeployment returns the deployment for the host
func GetDeployment(host string) (deploy *appsAPI.Deployment, err error) {
	deploys, err := k8sClient.AppsV1().Deployments("").List(
		metaAPI.ListOptions{
			LabelSelector: fmt.Sprintf("host=%s", host),
		},
	)
	if err != nil {
		return nil, Error{Type: K8sError, Err: err}
	}

	count := len(deploys.Items)
	switch count {
	case 1:
		return &deploys.Items[0], nil
	case 0:
		return nil, Error{Type: NoDeploymentError, Err: fmt.Errorf("No deployment with host '%s' found", host)}
	default: // >1
		return nil, Error{Type: TooManyDeploymentsError, Err: fmt.Errorf("expected exactly 1 Deployment with host label = '%s', found %d", host, count)}
	}
}

func SetRestartAnnotations(deploy *appsAPI.Deployment) error {
	return errors.New("TODO: Implement me")
}
