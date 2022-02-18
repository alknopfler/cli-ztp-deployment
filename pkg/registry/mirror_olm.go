package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/lib/indexer"
	"github.com/sirupsen/logrus"
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
	ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	//Get the registry route

	regName, err := r.GetRegistryRouteName(ctx, ocpclient)
	if err != nil {
		log.Printf(color.InRed("[ERROR] getting the Route Name for the registry: %s"), err.Error())
		return err
	}
	r.RegistryRoute = regName

	wgMirrorOLM.Add(3)

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

	//Login to the registry (tekton use case)
	go func() error {
		//Login to the registry to grab the authfile with the new registry credentials
		err := r.Login(ctx)
		if err != nil {
			log.Printf(color.InRed("[ERROR] login to registry: %s"), err.Error())
			return err
		}
		log.Println(color.InGreen("[INFO] login to registry successful"))
		wgMirrorOLM.Done()
		return nil
	}()

	wgMirrorOLM.Wait()

	return nil
}

func (r *Registry) mirrorOlm() error {

	logger := logrus.WithFields(logrus.Fields{"packages": r.RegistrySrcPkg})

	logger.Info("pruning the index")
	indexer.NewIndexPruner(containertools.NewContainerTool("podman", containertools.PodmanTool), logger)

	return nil
}
