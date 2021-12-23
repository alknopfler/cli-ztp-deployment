package resources

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetDeployments(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string) ([]v1.Deployment, error) {

	list, err := clientset.AppsV1().Deployments(namespace).
		List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func GetDeployment(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string, name string) (v1.Deployment, error) {

	d, err := clientset.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return v1.Deployment{}, err
	}
	return *d, nil
}
