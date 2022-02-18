package preflight

import (
	"context"
	"errors"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"os/exec"

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

//Run Preflight:
// - Check if the conditions are ready or not
// - Strategy: wait for all to get the error at the end in order to now where is the problem.
func (p *Preflight) RunPreflights() error {
	log.Println(color.InBold(color.InYellow(">>>> Running preflights")))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()

	wg.Add(7)
	var errNodes, errPVS, errCO, errMetal3, errCommands error
	go func() {
		errNodes = p.verifyNodes(*client, ctx)
	}()
	go func() {
		errPVS = p.verifyPVS(*client, ctx)
	}()
	go func() {
		errCO = p.verifyClusterOperators(dynamicClient, ctx)
	}()
	go func() {
		errMetal3 = p.verifyMetal3Pods(*client, ctx)
	}()
	go func() {
		errCommands = p.verifyCommand("podman")
	}()
	go func() {
		errCommands = p.verifyCommand("oc")
	}()
	go func() {
		errCommands = p.verifyCommand("skopeo")
	}()
	wg.Wait()

	if errNodes != nil || errPVS != nil || errCO != nil || errMetal3 != nil || errCommands != nil {
		return errors.New(color.InRed("[ERROR] Preflights failed"))
	}
	return nil
}

func (p *Preflight) verifyPVS(clientset kubernetes.Clientset, ctx context.Context) error {
	defer wg.Done()
	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting the PV info: %e"), err)
		return err
	}

	if len(pvs.Items) < 3 {
		log.Println(color.InRed("[ERROR] PV insufficients..."))
		return errors.New("[ERROR] PV insufficients...")
	}
	log.Println(color.InGreen(">>>>[OK] Pvs validated"))
	return nil
}

func (p *Preflight) verifyNodes(clientset kubernetes.Clientset, ctx context.Context) error {
	defer wg.Done()
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting Nodes info: %e"), err)
		return err
	}

	if len(nodes.Items) < 3 {
		log.Println(color.InRed("[ERROR] Nodes insufficient."))
		return errors.New("[ERROR] Nodes insufficient.")
	}
	log.Println(color.InGreen(">>>>[OK] Nodes validated"))
	return nil
}

func (p *Preflight) verifyClusterOperators(client dynamic.Interface, ctx context.Context) error {
	defer wg.Done()
	co, err := resources.NewGenericList(ctx, client, CLUSTER_OPERATOR_GROUP, CLUSTER_OPERATOR_VERSION, CLUSTER_OPERATOR_RESOURCE, "", CONDITION_CO_READY).GetResourcesByJq()
	if err != nil {
		log.Println(color.InRed("[ERROR] Error getting cluster operators info: %e"), err)
		return err
	}

	if len(co) > 0 {
		log.Println(color.InRed("[ERROR] Cluster operators are not available...Exiting"))
		return errors.New("[ERROR] Cluster operators are not available...Exiting")
	}
	log.Println(color.InGreen(">>>>[OK] Cluster Operators validated"))
	return nil
}

func (p *Preflight) verifyMetal3Pods(client kubernetes.Clientset, ctx context.Context) error {
	defer wg.Done()
	metal, err := client.CoreV1().Pods(METAL3_NAMESPACE).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Println(color.InRed("[ERROR] Error getting pods about metal3: %e"), err)
		return err
	}

	if len(metal.Items) < 1 {
		log.Println(color.InRed("[ERROR] Metal3 pods insufficient...Exiting"))
		return errors.New("[ERROR] Metal3 pods insufficient...Exiting")
	}
	log.Println(color.InGreen(">>>>[OK] Metal3 pods validated"))
	return nil
}

func (p *Preflight) verifyCommand(command string) error {
	defer wg.Done()
	return isCommandAvailable(command)
}

func isCommandAvailable(name string) error {
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		log.Println(color.InRed("[ERROR] '" + name + " is not installed...Exiting"))
		return err
	}
	log.Println(color.InGreen(">>>>[OK] Command '" + name + "' installed and validated"))
	return nil
}
