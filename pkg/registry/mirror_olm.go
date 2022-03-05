package registry

import (
	"bytes"
	"context"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	adm "github.com/openshift/oc/pkg/cli/admin/catalog"
	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/lib/indexer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
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

	//Prune catalog
	err = resources.Retry(4, 1*time.Minute, func() (err error) {
		return r.pruneCatalog()
	})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error Pruning the OLM the catalog: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen(">>>> [INFO] Prune the OLM catalog done successful"))

	//push the catalog pruned to registry
	err = resources.Retry(4, 1*time.Minute, func() (err error) {
		return r.pushCatalog()
	})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error Pushing the index to registry: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen(">>>> [INFO] Push index to registry  successfully"))

	//adm catalog to registry
	err = resources.Retry(4, 1*time.Minute, func() (err error) {
		return r.mirrorCatalog()
	})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error Mirroring the catalog to registry: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen(">>>> [INFO] Mirror the catalog to registry done successfully"))

	//skopeo the extra images
	err = r.copyExtraImages()
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error copy extra images to registry: %s"), err.Error())
		return err
	}

	log.Println(color.InGreen(">>>> [OK] Mirror OLM completely finished"))
	return nil

}

func (r *Registry) pruneCatalog() error {

	//TODO For to parallelize the mirroring of the olm with goroutines
	logger := logrus.WithFields(logrus.Fields{"packages": r.RegistrySrcPkg})

	logger.Info(color.InYellow(" >>>> [INFO] Pruning the index"))
	indexPruner := indexer.NewIndexPruner(containertools.NewContainerTool("podman", containertools.PodmanTool), logger)

	request := indexer.PruneFromIndexRequest{
		Generate:  true,
		FromIndex: r.RegistryOLMSourceIndex,
		Packages:  r.RegistrySrcPkgFormatted,
		Tag:       r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion,
	}

	err := indexPruner.PruneFromIndex(request)
	if err != nil {
		logger.Error(color.InRed(" >>>> [ERROR] Error pruning the index: %s"), err.Error())
		return err
	}
	logger.Info(color.InGreen(" >>>> [INFO] Pruning the index done successful"))
	return nil
}

func (r *Registry) pushCatalog() error {
	//TODO solve issue with podman https://github.com/alknopfler/cli-ztp-deployment/issues/1#issue-1147187533
	//Workarround using exec to push the catalog
	log.Println(color.InYellow(" >>>> [INFO] Doing: podman push " + r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion + " --authfile " + r.PullSecretTempFile))
	cmd := exec.Command("podman", "push", r.RegistryRoute+"/"+r.RegistryOLMDestIndexNS+":v"+config.Ztp.Config.OcOCPVersion, "--authfile", r.PullSecretTempFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	log.Println(color.InGreen("[INFO] Podman push done successfully"))
	return nil
}

func (r *Registry) mirrorCatalog() error {
	log.Println(color.InYellow("[INFO] Creating new Mirror Catalog"))
	o := adm.NewMirrorCatalogOptions(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	o.SecurityOptions.RegistryConfig = r.PullSecretTempFile
	o.SecurityOptions.SkipVerification = true
	o.ParallelOptions.MaxPerRegistry = 100
	o.ManifestDir = ""
	o.IndexPath = ""
	o.DryRun = false
	o.ManifestOnly = false
	o.FileDir = ""
	o.FromFileDir = ""
	o.MaxPathComponents = 2
	o.IcspScope = "repository"
	o.MaxICSPSize = 250000
	cmd := cobra.Command{}

	log.Println(color.InYellow("[INFO] Generating options source and destionation for the mirror catalog"))
	source := r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion
	dest := r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS
	err := o.Complete(&cmd, []string{source, dest})
	if err != nil {
		log.Println(color.InRed("[ERROR] Error generating options source and destination for the mirror catalog"))
		return err
	}
	log.Println(color.InYellow("[INFO] Doing the mirror catalog"))
	err = o.Run()
	if err != nil {
		log.Println(color.InRed("[ERROR] Error doing the mirror catalog"))
		return err
	}
	log.Println(color.InGreen("[INFO] Mirror Catalog done successfully"))
	return nil
}

func (r *Registry) copyExtraImages() error {
	log.Println(color.InYellow("[INFO] Copying extra images to the registry"))

	var wg sync.WaitGroup
	wg.Add(3)
	fatalErrors := make(chan error)
	wgDone := make(chan bool)

	for _, image := range r.RegistryExtraImages {
		log.Println(color.InYellow("[INFO] Copying extra image: " + image))
		go func() {
			if err := r.copyImage(image); err != nil {
				log.Printf(color.InRed("[ERROR] Updating the ca for the mirror olm: %s"), err.Error())
				fatalErrors <- err
			}
			wg.Done()
		}()
	}
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

	log.Println(color.InGreen("[INFO] Copying extra images to the registry done successfully"))
	return nil
}

func (r *Registry) copyImage(image string) error {
	//skopeo copy docker://${image} docker://${DESTINATION_REGISTRY}/${image#*/} --all --authfile ${PULL_SECRET}
	cmd := exec.Command("skopeo", "copy", "docker://"+image, "docker://"+r.RegistryRoute+"/"+strings.Join(strings.Split(image, "/")[1:], "/"), "--all", "--authfile", r.PullSecretTempFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%s\nerr:\n%s\n", outStr, errStr)
	return nil
}
