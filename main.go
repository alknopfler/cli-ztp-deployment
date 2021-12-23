package main

import (
	"github.com/alknopfler/cli-ztp-deployment/cmd"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
)

func main() {

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf("[ERROR] %e", err)
	}
}

func newCommand() *cobra.Command {
	c := &cobra.Command{
		Use:   "ztpcli",
		Short: "Ztp is a command line to deploy ztp openshift clusters",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	c.AddCommand(cmd.NewVerify())
	//cmd.AddCommand(command.NewDeploy())
	//cmd.AddCommand(command.NewWaitFor())

	return c
}
