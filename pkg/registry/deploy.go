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
)

func (r *Registry) RunDeployRegistry() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()

	//Step 1 - Create the namespace for the registry
	err := r.createNamespace(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating Namespace for the registry: %v"), err)
		return err
	}
	//Step 2 - Create the secret and config map for the registry
	err = r.createSecret(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating secret and config map for the registry: %v"), err)
		return err
	}
	err = r.createConfigMap(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating secret and config map for the registry: %v"), err)
		return err
	}

	return nil
}

func (r *Registry) createNamespace(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyNamespace(ctx, client); !found && err != nil {
		log.Printf(color.InYellow("Namespace %s not found, Creating it..."), r.RegistryNS)
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistryNS,
			},
		}
		_, err := client.CoreV1().Namespaces().Create(ctx, nsName, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating namespace: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Namespace for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Namespace for the registry already exists. Skipping creation"))
	return nil
}

//Func to create the secret for the registry
func (r *Registry) createSecret(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifySecret(ctx, client); !found && err != nil {
		log.Printf(color.InYellow("Secret and Config Map for the registry not found, Creating it..."))
		//create secret
		secret := &coreV1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistrySecretName,
			},
			Data: map[string][]byte{
				"htpasswd": []byte(r.RegistrySecretHash),
			},
		}

		_, err := client.CoreV1().Secrets(r.RegistryNS).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating secret: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Secret for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Secret for the registry already exists. Skipping creation"))
	return nil
}

//Func to create the config map for the registry
func (r *Registry) createConfigMap(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyConfigMap(ctx, client); !found && err != nil {
		log.Printf(color.InYellow("Secret and Config Map for the registry not found, Creating it..."))
		//create config map
		configMap := &coreV1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "registry-config",
			},
			Data: map[string]string{
				"config.yaml": `version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
compatibility:
  schema1:
    enabled: true`,
			},
		}
		_, err := client.CoreV1().ConfigMaps(r.RegistryNS).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating config map: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Config Map for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Config Map for the registry already exists. Skipping creation"))
	return nil
}
