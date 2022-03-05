package acm

import "github.com/alknopfler/cli-ztp-deployment/config"

type ACM struct {
	Name           string
	Source         string
	DefaultChannel string
	Csv            string
	Namespace      string
}

func NewACMDefault() *ACM {
	return &ACM{
		Name:           "advanced-cluster-management",
		Source:         "redhat-operators",
		DefaultChannel: "release-" + config.Ztp.Config.OcACMVersion,
		Csv:            "",
		Namespace:      "open-cluster-management",
	}
}
