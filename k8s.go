package main

import (

	// appsAPI "k8s.io/api/apps/v1"
	"log"

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
