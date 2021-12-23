package verify

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wg sync.WaitGroup

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.Set(config.Ztp.Config.KubeconfigHUB)
	dynamicClient := auth.SetWithDynamic(config.Ztp.Config.KubeconfigHUB)

	wg.Add(4)
	go verifyNodes(*client, ctx)
	go verifyPVS(*client, ctx)
	go verifyClusterOperators(dynamicClient, ctx)
	go verifyMetal3Pods(*client, ctx)
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

func verifyClusterOperators(client dynamic.Interface, ctx context.Context) {
	defer wg.Done()
	co, err := resources.GetResourcesByJq(client, ctx, "config.openshift.io", "v1", "clusteroperators", "", ".status.conditions[] | select (.type == \"Available\" and .status == \"False\")")
	if err != nil {
		log.Fatal(err)
	}

	if len(co) > 0 {
		log.Fatal("[ERROR] Cluster operators are not available...Exiting")
	}
	log.Println(">>>> co validated")
}

func verifyMetal3Pods(client kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	metal, err := resources.GetPods(&client, ctx, "openshift-machine-api")
	if err != nil {
		log.Fatal(err)
	}

	if len(metal) < 1 {
		log.Fatal("[ERROR] Metal3 pods insufficient...Exiting")
	}
	log.Println(">>>> Metal3 pods validated")
}
