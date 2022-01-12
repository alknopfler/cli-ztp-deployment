package cmd

import (
	"github.com/spf13/cobra"
)

func NewMirror() *cobra.Command {
	c := &cobra.Command{
		Use:          "mirror",
		Short:        "Commands to mirroring",
		SilenceUsage: true,
	}

	c.AddCommand(NewMirrorOcp())
	c.AddCommand(NewMirrorOlm())

	return c
}
