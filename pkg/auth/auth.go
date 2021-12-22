package auth

import (
	"fmt"
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Set(kubeconfig string) *kubernetes.Clientset {
	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Errorf("Error reading kubeconfig to get clientset: ", err)
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Error getting clientset: ", err)
		os.Exit(1)
	}
	return clientset
}

func SetWithDynamic(kubeconfig string) dynamic.Interface {

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Errorf("Error reading kubeconfig to get clientset: ", err)
		os.Exit(1)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Error getting clientset: ", err)
		os.Exit(1)
	}
	return client
}
