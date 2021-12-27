package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wgDeployRegistry sync.WaitGroup

func (r *Registry) RunDeployRegistry() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()

	wgDeployRegistry.Add(4)
	go func() error {
		err := r.createNamespace(ctx, client)
		wgDeployRegistry.Done()
		if err != nil {
			log.Fatalf("Error creating deployment: %v", err)
			return err
		}
		return nil
	}()
	return nil
}

func (r *Registry) createNamespace(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyNamespace(ctx, client); !found && err == nil {
		log.Printf(color.InYellow("Namespace %s not found, creating it..."), config.Ztp.Config.RegistryNamespace)
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: config.Ztp.Config.RegistryNamespace,
			},
		}
		_, err := client.CoreV1().Namespaces().Create(ctx, nsName, metav1.CreateOptions{})
		if err != nil {
			log.Fatalf("Error creating namespace: %v", err)
			return err
		}
	}
	return nil
}
