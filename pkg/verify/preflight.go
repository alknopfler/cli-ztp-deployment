package verify

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"log"
)

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := auth.Set(config.Ztp.Config.KubeconfigHUB)

	pvs, err := resources.GetPVS(c, ctx)
	if err != nil {
		log.Fatal(err)
	}

	if len(pvs.Items) < 3 {
		log.Fatal("Error PV insufficients...Exiting", err)
	}
	log.Println("Pvs validated")
	return nil
}
