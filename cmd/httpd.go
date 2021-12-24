package cmd

import (
	deploy "github.com/alknopfler/cli-ztp-deployment/pkg/deploy/httpd"
	"github.com/spf13/cobra"
)

func NewHTTPD() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Deploy new File Server running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			var f deploy.FileServer
			return f.RunHttpd()
		},
	}
	//TODO add flags to customize httpd deployment
	return cmd
}
