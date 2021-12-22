package verify

import (
	"context"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func RunPreflights() error {

	fmt.Println()
	fmt.Println("Using kubeconfig: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	api := clientset.CoreV1()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initial list
	listOptions := metav1.ListOptions{LabelSelector: label, FieldSelector: field}
	pvcs, err := api.PersistentVolumes(ns).List(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}

	printPVCs(pvcs)
	fmt.Println()
	return nil
}
