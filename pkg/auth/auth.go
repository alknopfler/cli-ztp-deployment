package auth

import (
	"fmt"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func Set(kubeconfig string) *kubernetes.Clientset {
	fmt.Println(">>>> Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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

func SetWithDynamic(kubeconfig string) dynamic.Interface {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("[ERROR] error reading kubeconfig to get clientset: %e", err)

	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("[ERROR] error getting clientset: %e", err)

	}
	return client
}
