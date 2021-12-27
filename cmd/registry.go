package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/registry"
	"github.com/spf13/cobra"
)

func NewDeployRegistry() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Deploy new File Server running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			r = registry.NewRegistry(mode)
			return r.RunDeployRegistry()
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of deployment for registry [hub|spoke]")
	return cmd
}

func NewVerifyRegistry() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Deploy new File Server running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			r = registry.NewRegistry(mode)
			r.RunVerifyHttpd()
			return nil
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of deployment for registry [hub|spoke]")
	return cmd
}
