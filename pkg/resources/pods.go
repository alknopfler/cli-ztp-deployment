package resources

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string) ([]v1.Pod, error) {

	list, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func GetPod(clientset *kubernetes.Clientset, ctx context.Context,
	namespace string, podname string) (v1.Pod, error) {

	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podname, metav1.GetOptions{})
	if err != nil {
		return v1.Pod{}, err
	}
	return *pod, nil
}
