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

	wgVerifyRegistry.Add(3)
	go func() {
		found, err := r.verifyNamespace(ctx, client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Namespace " + r.RegistryNS + " for registry not found"))
		} else {
			log.Println(color.InGreen("[OK] NameSpace " + r.RegistryNS + " for registry found"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		found, err := r.verifySecret(ctx, client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Secret " + r.RegistrySecretName + " for registry not found"))
		} else {
			log.Println(color.InGreen("[OK] Secret " + r.RegistrySecretName + " for registry found"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		found, err := r.verifyConfigMap(ctx, client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] ConfigMap " + "registry-conf" + " for registry not found"))
		} else {
			log.Println(color.InGreen("[OK] Configmap " + "registry-conf" + " for registry found"))
		}
		wgVerifyRegistry.Done()
	}()
}

func (r *Registry) verifyNamespace(ctx context.Context, client *kubernetes.Clientset) (bool, error) {
	_, err := client.CoreV1().Namespaces().Get(ctx, r.RegistryNS, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

//func to verify if the secret exists
func (r *Registry) verifySecret(ctx context.Context, client *kubernetes.Clientset) (bool, error) {
	_, err := client.CoreV1().Secrets(r.RegistryNS).Get(ctx, r.RegistrySecretName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

//func to verify if the configmap exists
func (r *Registry) verifyConfigMap(ctx context.Context, client *kubernetes.Clientset) (bool, error) {
	_, err := client.CoreV1().ConfigMaps(r.RegistryNS).Get(ctx, "registry-conf", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}
