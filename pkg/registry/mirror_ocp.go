package registry

import (
	"context"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	a "github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"log"
	"os"
)

func (r *Registry) RunMirrorOcp() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	//client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	//dynamicClient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetRouteAuth()

	regName, err := r.getRegistryRouteName(ctx, ocpclient)

	if err != nil {
		log.Printf(color.InRed("[ERROR] getting the Route Name for the registry: %e"), err)
		return err
	}
	args := []string{regName}
	loginOpts := a.LoginOptions{
		AuthFile:                  r.PullSecretTempFile,
		CertDir:                   r.RegistryPathCaCert,
		Password:                  r.RegistryPass,
		Username:                  r.RegistryUser,
		StdinPassword:             false,
		GetLoginSet:               false,
		Verbose:                   false,
		AcceptRepositories:        true,
		Stdin:                     os.Stdin,
		Stdout:                    os.Stdout,
		AcceptUnspecifiedRegistry: true,
	}
	sysCtx := &types.SystemContext{
		AuthFilePath:                loginOpts.AuthFile,
		DockerCertPath:              loginOpts.CertDir,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	err = a.Login(ctx, sysCtx, &loginOpts, args)
	if err != nil {
		log.Printf(color.InRed("[ERROR] Logging in to the registry: %e"), err)
		return err
	}
	fmt.Println("[INFO] Logging in to the registry")
	return nil
}
