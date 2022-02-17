package cmd

import (
	"errors"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/registry"
	"github.com/spf13/cobra"
)

func NewMirrorOlm() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "olm",
		Short: "Mirroring the OLM operators to the registry deployed based on mode (hub or spoke)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mode != config.MODE_HUB && mode != config.MODE_SPOKE {
				return errors.New(color.InRed("mode must be either hub or spoke"))
			}
			r = registry.NewRegistry(mode)
			return r.RunMirrorOlm()
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of working with the mirroring [hub|spoke]")
	//TODO add flag to get spoke name or ALL to deploy the registry
	return cmd
}

func NewVerifyMirrorOlm() *cobra.Command {
	var r *registry.Registry
	var mode string
	cmd := &cobra.Command{
		Use:   "mirror-olm",
		Short: "Verify if the OLM operators are successful based on mode (hub or spoke)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if mode != config.MODE_HUB && mode != config.MODE_SPOKE {
				return errors.New(color.InRed("mode must be either hub or spoke"))
			}
			r = registry.NewRegistry(mode)
			r.RunVerifyMirrorOlm()
			return nil
		},
	}
	flags := cmd.Flags()
	// Read the config flag directly into the struct, so it's immediately available.
	flags.StringVar(&mode, "mode", "", "Mode of working with the mirroring [hub|spoke]")
	//TODO add flag to get spoke name or ALL to verify the registry
	return cmd
}
