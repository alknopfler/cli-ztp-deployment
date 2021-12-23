package resources

import (
	"context"
	v1 "k8s.io/api/core/v1"

	_ "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Core struct {
	Client    *kubernetes.Clientset
	Ctx       context.Context
	Namespace string
	Name      string
}

func NewCore(ctx context.Context, client *kubernetes.Clientset) *Core {
	return &Core{
		Client: client,
		Ctx:    ctx,
	}
}

func NewCoreWithParam(ctx context.Context, client *kubernetes.Clientset, namespace, name string) *Core {
	return &Core{
		Client:    client,
		Ctx:       ctx,
		Namespace: namespace,
		Name:      name,
	}
}

func (c *Core) GetNodes() (*v1.NodeList, error) {
	return c.Client.CoreV1().Nodes().List(c.Ctx, metav1.ListOptions{})
}

func (c *Core) GetPods() ([]v1.Pod, error) {

	list, err := c.Client.CoreV1().Pods(c.Namespace).List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (c *Core) GetPod() (v1.Pod, error) {

	pod, err := c.Client.CoreV1().Pods(c.Namespace).Get(c.Ctx, c.Name, metav1.GetOptions{})
	if err != nil {
		return v1.Pod{}, err
	}
	return *pod, nil
}

func (c *Core) GetPVS() (*v1.PersistentVolumeList, error) {
	return c.Client.CoreV1().PersistentVolumes().List(c.Ctx, metav1.ListOptions{})
}
