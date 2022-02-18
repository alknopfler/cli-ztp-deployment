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

//Run Preflight:
// - Check if the conditions are ready or not
// - Strategy: wait for all to get the error at the end in order to now where is the problem.
func (p *Preflight) RunPreflights() error {
	log.Println(color.InBold(color.InYellow(">>>> Running preflights")))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()

	var wg sync.WaitGroup
	fatalError := make(chan error)
	wgDone := make(chan bool)

	wg.Add(7)

	go func() {
		err := p.verifyNodes(*client, ctx)
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyPVS(*client, ctx)
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyClusterOperators(dynamicClient, ctx)
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyMetal3Pods(*client, ctx)
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyCommand("podman2")
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyCommand("oc")
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()
	go func() {
		err := p.verifyCommand("skopeo")
		if err != nil {
			fatalError <- err
		}
		wg.Done()
	}()

	//func to wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(wgDone)
	}()

	//Select to wait for all goroutines to finish or for a fatal error
	select {
	case err := <-fatalError:
		return err
	case <-wgDone:
		return nil
	}
}

func (p *Preflight) verifyPVS(clientset kubernetes.Clientset, ctx context.Context) error {
	pvs, err := clientset.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting the PV info: %s"), err.Error())
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
	nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting Nodes info: %s"), err.Error())
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
	co, err := resources.NewGenericList(ctx, client, CLUSTER_OPERATOR_GROUP, CLUSTER_OPERATOR_VERSION, CLUSTER_OPERATOR_RESOURCE, "", CONDITION_CO_READY).GetResourcesByJq()
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting cluster operators info: %s"), err.Error())
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
	metal, err := client.CoreV1().Pods(METAL3_NAMESPACE).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting pods about metal3: %s"), err.Error())
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
