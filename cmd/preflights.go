package cmd

import (
	"github.com/spf13/cobra"
)

func NewPreflights() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "preflights",
		Short: "Run Preflight checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunPreflights()
		},
	}

	return cmd
}

func RunPreflights() error {
	return nil
}
