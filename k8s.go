package main

import (
	"fmt"
	"log"
	"time"

	appsAPI "k8s.io/api/apps/v1"
	metaAPI "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
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
func GetDeployment(host, ns string) (*appsAPI.Deployment, error) {
	deploys, err := k8sClient.AppsV1().Deployments(ns).List(
		metaAPI.ListOptions{
			LabelSelector: fmt.Sprintf("host=%s", host),
		},
	)
	if err != nil {
		return nil, Error{Code: K8sError, Err: err}
	}

	count := len(deploys.Items)
	switch count {
	case 1:
		return &deploys.Items[0], nil
	case 0:
		return nil, Error{Code: NoDeploymentError, Err: fmt.Errorf("No deployment with host '%s' found", host)}
	default: // >1
		return nil, Error{Code: TooManyDeploymentsError, Err: fmt.Errorf("expected exactly 1 Deployment with host label = '%s', found %d", host, count)}
	}
}

// RestartPods restarts the Deployment's Pods
//
// This is achived by patching the Deployment and adding/updating the
// restartedAt/restartReason annotations of its Pods' template.
//
// These annotations are also set in the Deployment.
func RestartPods(d *appsAPI.Deployment, reason string) error {
	annotations := fmt.Sprintf(`
	    "annotations": {
			"restartedAt": "%s",
			"restartReason": "%s"
		}`,
		time.Now().UTC().Format(time.RFC3339),
		reason,
	)

	// the annotations are set both for the Deployment and for the pod
	// template (which will ultimately trigger the restart)
	patch := fmt.Sprintf(`{
			"metadata": {
				%s
			},
			"spec": {
				"template": {
					"metadata": {
						%s
					}
				}
			}
		}`,
		annotations,
		annotations,
	)

	_, err := k8sClient.AppsV1().Deployments(d.Namespace).Patch(
		d.Name,
		types.StrategicMergePatchType,
		[]byte(patch),
	)
	if err != nil {
		return Error{Code: K8sError, Err: err}
	}
	return nil
}
