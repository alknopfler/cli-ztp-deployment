package cmd

import (
	"github.com/spf13/cobra"
)

func NewVerify() *cobra.Command {
	c := &cobra.Command{
		Use:          "verify",
		Short:        "Commands to verify things",
		SilenceUsage: true,
	}

	c.AddCommand(NewPreflights())
	c.AddCommand(NewVerifyHTTPD())

	return c
}
