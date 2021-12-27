package main

import (
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/cmd"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

func init() {

	//Load config from file to set globally
	var err error
	config.Ztp, err = config.NewConfig()
	if err != nil {
		//There is some error in config file or creating the config ZTP.
		log.Fatal(color.InRed(err.Error()))

	}

}

func main() {
	command := newCommand()
	if err := command.Execute(); err != nil {
		log.Fatalf(color.InRed("[ERROR] %e"), err)
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

	c.AddCommand(cmd.NewVerify())
	c.AddCommand(cmd.NewDeploy())

	return c
}
