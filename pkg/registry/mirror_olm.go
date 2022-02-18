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

	var wg sync.WaitGroup
	wg.Add(3)
	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	//Update Trust CA if not present (tekton  use case)
	go func() {
		if err := r.UpdateTrustCA(ctx, client); err != nil {
			log.Printf(color.InRed("[ERROR] Updating the ca for the mirror olm: %s"), err.Error())
			fatalErrors <- err
		}
		wg.Done()
	}()

	//Create the catalog source if not present
	go func() {
		if err := r.CreateCatalogSource(ctx); err != nil {
			log.Printf(color.InRed("[ERROR] Error creating catalog source for the mirror olm: %s"), err.Error())
			fatalErrors <- err
		}
		log.Println(color.InGreen("[INFO] Created the catalog source successfully"))
		wg.Done()
	}()

	//Login to the registry (tekton use case)
	go func() {
		//Login to the registry to grab the authfile with the new registry credentials
		err := r.Login(ctx)
		if err != nil {
			log.Printf(color.InRed("[ERROR] login to registry: %s"), err.Error())
			fatalErrors <- err
		}
		log.Println(color.InGreen("[INFO] Login to registry successful"))
		wg.Done()
	}()

	go func() {
		wg.Wait()
		close(wgDone)
	}()

	// Wait until either WaitGroup is done or an error is received through the channel
	select {
	case <-wgDone:
		// carry on
		break
	case err := <-fatalErrors:
		close(fatalErrors)
		return err
	}

	errMirrorOlm := r.PruneCatalog()
	if errMirrorOlm != nil {
		log.Printf(color.InRed("[ERROR] Error pruning the olm catalog: %s"), errMirrorOlm.Error())
		return errMirrorOlm
	}

	return nil
}

func (r *Registry) PruneCatalog() error {

	//TODO For to parallelize the mirroring of the olm with goroutines
	logger := logrus.WithFields(logrus.Fields{"packages": r.RegistrySrcPkg})

	logger.Info(color.InYellow("[INFO] >>>> Pruning the index"))
	indexPruner := indexer.NewIndexPruner(containertools.NewContainerTool("podman", containertools.PodmanTool), logger)

	request := indexer.PruneFromIndexRequest{
		Generate:  true,
		FromIndex: r.RegistryOLMSourceIndex,
		Packages:  r.RegistrySrcPkgFormatted,
		Tag:       r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion,
	}

	err := indexPruner.PruneFromIndex(request)
	if err != nil {
		return err
	}
	return nil
}
