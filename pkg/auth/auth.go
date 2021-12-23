package auth

import (
	"fmt"
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

func (z *ZTPAuth) Set() *kubernetes.Clientset {
	fmt.Println(">>>> Using kubeconfig: ", z.KubeConfig)
	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf("[ERROR] error reading kubeconfig to get clientset: %e", err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("[ERROR] error getting clientset: %e", err)
	}
	return clientset
}

func (z *ZTPAuth) SetWithDynamic() dynamic.Interface {

	config, err := clientcmd.BuildConfigFromFlags("", z.KubeConfig)
	if err != nil {
		log.Fatalf("[ERROR] error reading kubeconfig to get clientset: %e", err)

	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("[ERROR] error getting clientset: %e", err)

	}
	return client
}
