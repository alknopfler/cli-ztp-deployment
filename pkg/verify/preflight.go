package verify

import (
	"cli-ztp-deployment/pkg/resources"
	"context"
	"fmt"
	"log"

	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
)

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := auth.Set(config.Ztp.Config.KubeconfigHUB)
	pvs, err := resources.GetPVS(c, ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(pvs)
	return nil
}
