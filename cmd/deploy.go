package cmd

import (
	"github.com/spf13/cobra"
)

func NewDeploy() *cobra.Command {
	c := &cobra.Command{
		Use:          "deploy",
		Short:        "Commands to deploy things",
		SilenceUsage: true,
	}

	c.AddCommand(NewHTTPD())

	return c
}
