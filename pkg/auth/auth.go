package auth

import (
	"github.com/TwiN/go-color"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type ZTPAuth struct {
	KubeConfig string
}

func NewZTPAuth(kubeconfig string) *ZTPAuth {
	return &ZTPAuth{
		KubeConfig: kubeconfig,
	}
}

func (z *ZTPAuth) GetProjectAuth() *projectv1.ProjectV1Client {
	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error reading kubeconfig to get clientset: %e"), err)
	}

	// create the clientset
	clientset, err := projectv1.NewForConfig(config)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error getting project ocp clientset: %e"), err)
	}
	return clientset
}

func (z *ZTPAuth) GetRouteAuth() *routev1.RouteV1Client {
	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error reading kubeconfig to get clientset: %e"), err)
	}

	// create the clientset
	clientset, err := routev1.NewForConfig(config)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error getting route ocp clientset: %e"), err)
	}
	return clientset
}

func (z *ZTPAuth) GetAuth() *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error reading kubeconfig to get clientset: %e"), err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error getting clientset: %e"), err)
	}
	return clientset
}

func (z *ZTPAuth) GetAuthWithGeneric() dynamic.Interface {

	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error reading kubeconfig to get clientset: %e"), err)

	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf(color.InRed("[ERROR] error getting clientset: %e"), err)

	}
	return client
}
