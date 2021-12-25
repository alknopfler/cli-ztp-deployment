package cmd

import (
	deploy "github.com/alknopfler/cli-ztp-deployment/pkg/httpd"
	"github.com/spf13/cobra"
)

func NewHTTPD() *cobra.Command {
	var f *deploy.FileServer
	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Deploy new File Server running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			f = deploy.NewFileServerDefault()
			//f.RunDeployHttpd()
			//TODO f.RunVerifyHttpd()
			return f.RunDeployHttpd()
		},
	}
	//TODO set flags to customize the file server
	return cmd
}
