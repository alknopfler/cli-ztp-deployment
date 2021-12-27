package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/registry"
	"github.com/spf13/cobra"
)

func NewDeployRegistry() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Deploy new File Server running on the hub cluster based on mode (hub | spoke)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r = registry.NewRegistry(mode)
			return r.RunDeployRegistry()
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of deployment for registry [hub|spoke]")
	//TODO add flag to get spoke name or ALL to deploy the registry
	return cmd
}

func NewVerifyRegistry() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Verify if registry is running on the server based on mode (hub | spoke)",
		RunE: func(cmd *cobra.Command, args []string) error {
			r = registry.NewRegistry(mode)
			r.RunVerifyRegistry()
			return nil
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of deployment for registry [hub|spoke]")
	//TODO add flag to get spoke name or ALL to verify the registry
	return cmd
}
