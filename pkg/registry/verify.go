package registry

import (
	"context"
	"errors"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wgVerifyRegistry sync.WaitGroup

func (r *Registry) RunVerifyRegistry() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	routeClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()

	wgVerifyRegistry.Add(9)
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
	go func() {
		found, err := r.verifyDeployment(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Deployment Registry not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Deployment Registry found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Deployment Registry found and ready"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		found, err := r.verifyService(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Service Registry not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Service Registry found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Service Registry found and ready"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		found, err := r.verifyRoute(ctx, *routeClient)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Route Registry not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Route Registry found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Route Registry found and ready"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		found, err := r.verifyPVC(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] PVC Registry not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] PVC Registry found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] PVC Registry found and ready"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		err := r.verifyMachineConfig(ctx, dynamicClient)
		if err != nil {
			log.Println(color.InRed("[ERROR] Machine Config not applied "))
		} else {
			log.Println(color.InGreen("[OK] Machine Config found and applied ready"))
		}
		wgVerifyRegistry.Done()
	}()
	go func() {
		err := r.verifyMCP(ctx, dynamicClient)
		if err != nil {
			log.Println(color.InRed("[ERROR] Machine Config and MCP not updated "))
		} else {
			log.Println(color.InGreen("[OK] Machine Config found and MCP ready"))
		}
		wgVerifyRegistry.Done()
	}()
	wgVerifyRegistry.Wait()
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
	_, err := client.CoreV1().ConfigMaps(r.RegistryNS).Get(ctx, r.RegistryConfigMapName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Registry) verifyDeployment(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	deployment, err := client.AppsV1().Deployments(r.RegistryNS).Get(ctx, r.RegistryDeploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if deployment.Status.AvailableReplicas < 1 {
		return true, errors.New(color.InRed("[ERROR] Deployment is not ready, replicas available insufficent"))
	}
	return true, nil
}

func (r *Registry) verifyRoute(ctx context.Context, client routev1.RouteV1Client) (found bool, err error) {
	_, err = client.Routes(r.RegistryNS).Get(ctx, r.RegistryRouteName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Registry) verifyService(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	service, err := client.CoreV1().Services(r.RegistryNS).Get(ctx, r.RegistryServiceName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if service.Spec.Ports[0].Port != 443 {
		return true, errors.New(color.InRed("[ERROR] verifying Service Registry: Service port is not ready"))
	}
	return true, nil
}

func (r *Registry) verifyPVC(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	pvc, err := client.CoreV1().PersistentVolumeClaims(r.RegistryNS).Get(ctx, r.RegistryPVCName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if pvc.Status.Phase != "Bound" {
		return true, errors.New(color.InRed("[ERROR] verifying PVC Registry: PVC is not bound"))
	}
	return true, nil
}

func (r *Registry) verifyMCP(ctx context.Context, client dynamic.Interface) error {
	mcp, err := resources.NewGenericGet(ctx, client, "machineconfiguration.openshift.io", "v1", "machineconfigpools", "", "master", ".status.conditions[] | select (.type == \"Updated\")").GetResourcesByJq()
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting MachineConfig info: %e"), err.Error())
		return err
	}

	if len(mcp) == 0 {
		log.Println(color.InRed("[ERROR] MCP Query returned no results"))
		return errors.New("[ERROR] MCP Query returned no results")
	}

	log.Println(color.InGreen("[OK] MCP Query returned results"))
	return nil
}

func (r *Registry) verifyMachineConfig(ctx context.Context, client dynamic.Interface) error {
	mcp, err := resources.NewGenericGet(ctx, client, "machineconfiguration.openshift.io", "v1", "machineconfig", "", "update-localregistry-ca-certs", "").GetResourceDynamically()
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error getting MachineConfig info: %e"), err.Error())
		return err
	}

	if mcp.Object == nil {
		log.Println(color.InRed("[ERROR] MachineConfig not found or applied"))
		return errors.New("[ERROR] MachineConfig not found or applied")
	}
	return nil
}

func (r *Registry) RunVerifyMirrorOcp() {

	/*ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	routeClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

	*/
}

func (r *Registry) RunVerifyMirrorOlm() {

	/*ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	routeClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

	*/
}
