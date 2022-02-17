package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"log"
	"sync"
)

var wgMirrorOLM sync.WaitGroup

func (r *Registry) RunMirrorOlm() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	//dynamicClient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuthWithGeneric()
	//ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	wgMirrorOLM.Add(2)

	//Update Trust CA if not present (tekton  use case)
	go func() error {
		if err := r.UpdateTrustCA(ctx, client); err != nil {
			log.Printf(color.InRed("[ERROR] Updating the ca for the mirror olm: %s"), err.Error())
			return err
		}
		wgMirrorOLM.Done()
		return nil
	}()

	//Create the catalog source if not present
	go func() error {
		if err := r.CreateCatalogSource(ctx); err != nil {
			log.Printf(color.InRed("[ERROR] Error creating catalog source for the mirror olm: %s"), err.Error())
			return err
		}
		wgMirrorOLM.Done()
		return nil
	}()

	wgMirrorOLM.Wait()

	return nil
}
