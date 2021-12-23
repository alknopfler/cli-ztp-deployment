package verify

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wg sync.WaitGroup

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.Set(config.Ztp.Config.KubeconfigHUB)

	wg.Add(2)
	go verifyNodes(*client, ctx)
	go verifyPVS(*client, ctx)
	wg.Wait()
	return nil
}

func verifyPVS(clientset kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	pvs, err := resources.GetPVS(&clientset, ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(pvs.Items) < 3 {
		log.Fatal("[ERROR] PV insufficients...Exiting")
	}
	log.Println(">>>> Pvs validated")
}

func verifyNodes(clientset kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	nodes, err := resources.GetNodes(&clientset, ctx)
	if err != nil {
		log.Fatal("[ERROR] ", err)
	}

	if len(nodes.Items) < 3 {
		log.Fatal("[ERROR] nodes insufficient...Exiting")
	}
	log.Println(">>>> Nodes validated")
}
