package auth

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)


func withAuth(kubeconfig string) *kubernetes.Clientset {
	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Errorf("Error reading kubeconfig to get clientset: ",err)
		os.Exit(1)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Errorf("Error getting clientset: ",err)
		os.Exit(1)
	}
	return clientset
}