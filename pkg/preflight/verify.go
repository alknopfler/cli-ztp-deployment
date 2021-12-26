package preflight

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"log"

	"sync"
)

const (
	CLUSTER_OPERATOR_GROUP    = "config.openshift.io"
	CLUSTER_OPERATOR_VERSION  = "v1"
	CLUSTER_OPERATOR_RESOURCE = "clusteroperators"
	CONDITION_CO_READY        = ".status.conditions[] | select (.type == \"Available\" and .status == \"False\")"
	METAL3_NAMESPACE          = "openshift-machine-api"
)

type Preflight struct{}

var wg sync.WaitGroup

func (p *Preflight) RunPreflights() error {
	log.Println(">>>> Running preflights")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()

	wg.Add(4)
	go p.verifyNodes(*client, ctx)
	go p.verifyPVS(*client, ctx)
	go p.verifyClusterOperators(dynamicClient, ctx)
	go p.verifyMetal3Pods(*client, ctx)
	wg.Wait()
	return nil
}

func (p *Preflight) verifyPVS(clientset kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("[ERROR] Error getting the PV info: %e", err)
	}

	if len(pvs.Items) < 3 {
		log.Println("[ERROR] PV insufficients...")
	}
	log.Println(">>>>[OK] Pvs validated")
}

func (p *Preflight) verifyNodes(clientset kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("[ERROR] Error getting Nodes info: %e", err)
	}

	if len(nodes.Items) < 3 {
		log.Println("[ERROR] Nodes insufficient.")
	}
	log.Println(">>>>[OK] Nodes validated")
}

func (p *Preflight) verifyClusterOperators(client dynamic.Interface, ctx context.Context) {
	defer wg.Done()
	co, err := resources.NewGenericList(ctx, client, CLUSTER_OPERATOR_GROUP, CLUSTER_OPERATOR_VERSION, CLUSTER_OPERATOR_RESOURCE, "", CONDITION_CO_READY).GetResourcesByJq()
	if err != nil {
		log.Fatal(err)
	}

	if len(co) > 0 {
		log.Fatal("[ERROR] Cluster operators are not available...Exiting")
	}
	log.Println(">>>> co validated")
}

func (p *Preflight) verifyMetal3Pods(client kubernetes.Clientset, ctx context.Context) {
	defer wg.Done()
	metal, err := client.CoreV1().Pods(METAL3_NAMESPACE).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	if len(metal.Items) < 1 {
		log.Fatal("[ERROR] Metal3 pods insufficient...Exiting")
	}
	log.Println(">>>> Metal3 pods validated")
}
