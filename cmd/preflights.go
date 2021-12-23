package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/verify"
	"github.com/spf13/cobra"
)

func NewPreflights() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "preflights",
		Short: "Run Preflight checks to validate the future deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			return verify.RunPreflights()
		},
	}

	return cmd
}
