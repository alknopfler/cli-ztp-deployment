package cmd

import (
	"github.com/alknopfler/cli-ztp-deployment/pkg/httpd"
	"github.com/spf13/cobra"
)

func NewDeployHTTPD() *cobra.Command {
	var f *httpd.FileServer
	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Deploy new File Server running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			f = httpd.NewFileServerDefault()
			return f.RunDeployHttpd()
		},
	}
	//TODO set flags to customize the file server
	return cmd
}

func NewVerifyHTTPD() *cobra.Command {
	var f *httpd.FileServer
	cmd := &cobra.Command{
		Use:   "httpd",
		Short: "Verify if File Server is running on the hub cluster ",
		RunE: func(cmd *cobra.Command, args []string) error {
			f = httpd.NewFileServerDefault()
			f.RunVerifyHttpd()
			return nil
		},
	}
	//TODO set flags to customize the file server
	return cmd
}
