package resources

import (
	"context"
	v1 "k8s.io/api/core/v1"

	_ "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPVS(clientset *kubernetes.Clientset, ctx context.Context) (*v1.PersistentVolumeList, error) {
	return clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
}
