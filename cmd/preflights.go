package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/preflight"
	"github.com/spf13/cobra"
)

func NewPreflights() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "preflights",
		Short: "Run Preflight checks to validate the future deployments",
		RunE: func(cmd *cobra.Command, args []string) error {
			var p preflight.Preflight
			return p.RunPreflights()
		},
	}
	return cmd
}
