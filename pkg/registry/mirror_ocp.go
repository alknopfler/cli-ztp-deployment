package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

func (r *Registry) RunMirrorOcp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	regName, err := r.getRegistryRouteName(ctx, ocpclient)
	if err != nil {
		log.Printf(color.InRed("[ERROR] getting the Route Name for the registry: %e"), err)
		return err
	}

	//TODO get pull secret into a temporal file

	return nil
}

func (r *Registry) getRegistryRouteName(ctx context.Context, client *routev1.RouteV1Client) (string, error) {
	ro, err := client.Routes(r.RegistryNS).Get(ctx, r.RegistryRouteName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return ro.Status.Ingress[0].Host, nil
}
