package verify

import (
	"context"
	"fmt"
	"log"

	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	resource "github.com/alknopfler/cli-ztp-deployment/pkg/resources"
)

func RunPreflights() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := auth.SetWithDynamic(config.Ztp.Config.KubeconfigHUB)
	pvcs, err := resource.GetResourcesDynamically(c, ctx)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()
	return nil
}
