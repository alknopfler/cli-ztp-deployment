package resources

import (
	"context"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"time"
)

//func to wait for all replicas in a deployment to be ready
func WaitForDeployment(ctx context.Context, deployment *appsv1.Deployment, client *kubernetes.Clientset) error {
	err := wait.Poll(time.Second, time.Minute, func() (bool, error) {
		deployment, err := client.AppsV1().Deployments(deployment.Namespace).Get(ctx, deployment.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if deployment.Status.ReadyReplicas == deployment.Status.Replicas {
			return true, nil
		}
		return false, nil
	})
	return err
}
