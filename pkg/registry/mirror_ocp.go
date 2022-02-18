package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"sync"

	adm "github.com/openshift/oc/pkg/cli/admin/release"
	"github.com/openshift/oc/pkg/cli/image/manifest"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"time"

	"log"
	"os"
)

var wgMirrorOCP sync.WaitGroup

func (r *Registry) RunMirrorOcp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	//dynamicClient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	wgMirrorOCP.Add(2)

	//Update Trust CA if not present (tekton  use case)
	go func() error {
		if err := r.UpdateTrustCA(ctx, client); err != nil {
			log.Printf(color.InRed("[ERROR] Updating the ca for the mirror ocp: %s"), err.Error())
			return err
		}
		wgMirrorOCP.Done()
		return nil
	}()

	//Get the registry route
	go func() error {
		regName, err := r.GetRegistryRouteName(ctx, ocpclient)
		if err != nil {
			log.Printf(color.InRed("[ERROR] getting the Route Name for the registry: %s"), err.Error())
			return err
		}
		r.RegistryRoute = regName
		wgMirrorOCP.Done()
		return nil
	}()

	wgVerifyRegistry.Wait()

	//Login to the registry to grab the authfile with the new registry credentials
	err := r.Login(ctx)
	if err != nil {
		log.Printf(color.InRed("[ERROR] login to registry: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen("[INFO] login to registry successful"))

	//Mirror ocp with a retry strategic to avoid errors
	err = resources.Retry(4, 1*time.Minute, func() (err error) {
		return r.mirrorOcp()
	})
	if err != nil {
		log.Printf(color.InRed("[ERROR] mirroring the OCP image: %s"), err.Error())
		return err
	}
	log.Println(color.InGreen("[INFO] mirroring the OCP image successful"))
	return nil
}

func (r *Registry) mirrorOcp() error {
	opt := adm.MirrorOptions{
		IOStreams: genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		SecurityOptions: manifest.SecurityOptions{
			RegistryConfig: r.PullSecretTempFile,
		},
		ParallelOptions: manifest.ParallelOptions{},
		From:            r.RegistryOCPReleaseImage,
		To:              r.RegistryRoute + "/" + r.RegistryOCPDestIndexNS,
		ToRelease:       r.RegistryRoute + "/" + r.RegistryOCPDestIndexNS + ":" + config.Ztp.Config.OcOCPTag,
		SkipRelease:     false,
		DryRun:          false,
		ImageStream:     nil,
		TargetFn:        nil,
	}
	opt.Validate()

	return opt.Run()
}
