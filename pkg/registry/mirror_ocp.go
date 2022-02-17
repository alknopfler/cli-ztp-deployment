package registry

import (
	"context"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	a "github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"github.com/openshift/oc/pkg/cli/image/manifest"

	adm "github.com/openshift/oc/pkg/cli/admin/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"log"
	"os"
)

func (r *Registry) RunMirrorOcp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	//client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode))sie.GetAuth()
	//dynamicClient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	regName, err := r.getRegistryRouteName(ctx, ocpclient)
	if err != nil {
		log.Printf(color.InRed("[ERROR] getting the Route Name for the registry: %e"), err)
		return err
	}
	r.RegistryRoute = regName
	if r.login(ctx) != nil {
		log.Printf(color.InRed("[ERROR] login to registry: %e"), err)
		return err
	}
	log.Println(color.InGreen("[INFO] login to registry successful"))

	if r.mirrorOcp(ctx) != nil {
		log.Printf(color.InRed("[ERROR] mirroring the OCP image: %e"), err)
		return err
	}
	log.Println(color.InGreen("[INFO] mirroring the OCP image successful"))
	return nil
}

func (r *Registry) login(ctx context.Context) error {
	args := []string{r.RegistryRoute}
	loginOpts := a.LoginOptions{
		AuthFile:      r.PullSecretTempFile,
		CertDir:       r.RegistryPathCaCert,
		Password:      r.RegistryPass,
		Username:      r.RegistryUser,
		StdinPassword: false,
		GetLoginSet:   false,
		//Verbose:                   false,
		//AcceptRepositories:        true,
		Stdin:                     os.Stdin,
		Stdout:                    os.Stdout,
		AcceptUnspecifiedRegistry: true,
	}
	sysCtx := &types.SystemContext{
		AuthFilePath:                loginOpts.AuthFile,
		DockerCertPath:              loginOpts.CertDir,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	return a.Login(ctx, sysCtx, &loginOpts, args)

}

func (r *Registry) mirrorOcp(ctx context.Context) error {
	opt := adm.MirrorOptions{
		IOStreams:       genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr},
		SecurityOptions: manifest.SecurityOptions{},
		ParallelOptions: manifest.ParallelOptions{},
		From:            r.RegistryOCPReleaseImage,
		To:              r.RegistryRoute + "/" + r.RegistryOCPDestIndexNS,
		ToRelease:       r.RegistryRoute + "/" + r.RegistryOCPDestIndexNS + ":" + config.Ztp.Config.OcOCPTag,
		SkipRelease:     false,
		DryRun:          false,
		ImageStream:     nil,
		TargetFn:        nil,
	}
	return opt.Run()
}
