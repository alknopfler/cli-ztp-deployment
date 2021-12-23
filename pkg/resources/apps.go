package resources

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type App struct {
	ClientSet *kubernetes.Clientset
	Ctx       context.Context
	Namespace string
	Name      string
}

//Constructor
func NewApp(ctx context.Context, clientSet *kubernetes.Clientset, namespace, name string) *App {
	return &App{
		ClientSet: clientSet,
		Ctx:       ctx,
		Namespace: namespace,
		Name:      name,
	}
}

//Func GetDeployments to get the list of deployment
func (a *App) GetDeployments() ([]v1.Deployment, error) {

	list, err := a.ClientSet.AppsV1().Deployments(a.Namespace).
		List(a.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

//Func GetDeployment to get a deployment with its name
func (a *App) GetDeployment() (v1.Deployment, error) {

	d, err := a.ClientSet.AppsV1().Deployments(a.Namespace).Get(a.Ctx, a.Name, metav1.GetOptions{})
	if err != nil {
		return v1.Deployment{}, err
	}
	return *d, nil
}
