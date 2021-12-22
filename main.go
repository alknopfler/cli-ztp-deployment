package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/alknopfler/cli-ztp-deployment/cmd"
	"github.com/alknopfler/cli-ztp-deployment/config"
)

func main() {

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	command := newCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	c.AddCommand(cmd.NewVerify())
	//cmd.AddCommand(command.NewDeploy())
	//cmd.AddCommand(command.NewWaitFor())

	return c
}
