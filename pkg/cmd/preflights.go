package cmd

import (
	"fmt"
	"os"

	config "github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/spf13/cobra"
)

func NewVerify() *cobra.Command {
	err := config.NewConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:   "preflights",
		Short: "Run Preflight checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunVerify()
		},
	}

	return cmd
}

func RunVerify() error {

	return nil
}
