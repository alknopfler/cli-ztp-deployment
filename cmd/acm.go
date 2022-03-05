package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/acm"
	"github.com/spf13/cobra"
)

func NewDeployACM() *cobra.Command {
	var a *acm.ACM
	cmd := &cobra.Command{
		Use:   "acm",
		Short: "Deploy ACM (Advanced Cluster Management) ",
		RunE: func(cmd *cobra.Command, args []string) error {
			a = acm.NewACMDefault()
			return a.RunDeployACM()
		},
	}

	return cmd
}

func NewVerifyACM() *cobra.Command {
	var a *acm.ACM
	cmd := &cobra.Command{
		Use:   "acm",
		Short: "Verify if ACM is running",
		RunE: func(cmd *cobra.Command, args []string) error {
			a = acm.NewACMDefault()
			a.RunVerifyACM()
			return nil
		},
	}

	return cmd
}
