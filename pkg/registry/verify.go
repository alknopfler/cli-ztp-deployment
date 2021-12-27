package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wgVerifyRegistry sync.WaitGroup

func (r *Registry) RunVerifyRegistry() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()

	wgVerifyRegistry.Add(4)
	go func() {
		found, err := r.verifyNamespace(ctx, client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Namespace " + config.Ztp.Config.RegistryNamespace + " for registry not found"))
		} else {
			log.Println(color.InGreen("[OK] NameSpace " + config.Ztp.Config.RegistryNamespace + " for registry found"))
		}
		wgVerifyRegistry.Done()
	}()
}

func (r *Registry) verifyNamespace(ctx context.Context, client *kubernetes.Clientset) (bool, error) {
	_, err := client.CoreV1().Namespaces().Get(ctx, config.Ztp.Config.RegistryNamespace, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}
