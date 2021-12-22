package verify

import (
	"context"
	"fmt"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func RunPreflights() error {

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
