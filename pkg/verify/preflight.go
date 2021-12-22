package verify

import (
	"context"
	"fmt"
	"log"

	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// initial list
	listOptions := metav1.ListOptions{LabelSelector: label, FieldSelector: field}
	client := auth.DynamicWithAuth(config.Ztp.Config.KubeconfigHUB)
	pvcs, err := client.
	if err != nil {
		log.Fatal(err)
	}

	printPVCs(pvcs)
	fmt.Println()
	return nil
}
