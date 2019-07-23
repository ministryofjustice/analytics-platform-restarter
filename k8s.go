package main

import (
	"fmt"
	"log"

	appsAPI "k8s.io/api/apps/v1"
	metaAPI "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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
func GetDeployment(host string) (*appsAPI.Deployment, error) {
	deps, err := k8sClient.AppsV1().Deployments("").List(
		metaAPI.ListOptions{
			LabelSelector: fmt.Sprintf("host=%s", host),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed listing deployments: %s", err)
	}

	count := len(deps.Items)
	switch count {
	case 0:
		return nil, nil
	case 1:
		return &deps.Items[0], nil
	default:
		return nil, fmt.Errorf("expected 1 or no deployments with host label, found %d", count)
	}
}
