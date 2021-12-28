package registry

import "github.com/alknopfler/cli-ztp-deployment/config"

//type FileServer
type Registry struct {
	Mode                        string
	KubeframeNS                 string
	RegistryNS                  string
	MarketNS                    string
	RegistryConfigFile          string
	RegistrySrcPkg              string
	RegistrySrcPkgFormatted     string
	RegistryExtraImages         string
	OcDisCatalog                string
	OcpReleaseFull              string
	RegistryUser                string
	RegistryPass                string
	RegistrySecretHash          string
	RegistrySecretName          string
	RegistryConfigMapName       string
	RegistryDeploymentName      string
	RegistryDataMountPath       string
	RegistryCertMountPath       string
	RegistryAutoSecretMountPath string
	RegistryConfMountPath       string
	RegistryServiceName         string
	RegistryRouteName           string
	RegistryPVCName             string
	RegistryPVMode              string
}

//Constructor NewFileServer
func NewRegistry(mode string) *Registry {
	return &Registry{
		Mode:                        mode,
		KubeframeNS:                 "kubeframe",
		MarketNS:                    "openshift-marketplace",
		OcDisCatalog:                "kubeframe-catalog",
		OcpReleaseFull:              config.Ztp.Config.OcOCPVersion + ".0",
		RegistryNS:                  "kubeframe-registry",
		RegistryConfigFile:          "config.yml",
		RegistrySrcPkg:              "kubernetes-nmstate-operator,metallb-operator,ocs-operator,local-storage-operator,advanced-cluster-management",
		RegistrySrcPkgFormatted:     "kubernetes-nmstate-operator metallb-operator ocs-operator local-storage-operator advanced-cluster-management",
		RegistryExtraImages:         "quay.io/jparrill/registry:2",
		RegistryUser:                "dummy",
		RegistryPass:                "dummy",
		RegistrySecretHash:          "dummy:$2y$05$VYlWo5DJrfSddVPrGWREwuuy8K.UgMoPoH2pSQpxPxwSiHrWbMa22",
		RegistrySecretName:          "auth",
		RegistryConfigMapName:       "registry-conf",
		RegistryDeploymentName:      "kubeframe-registry",
		RegistryServiceName:         "kubeframe-registry",
		RegistryRouteName:           "kubeframe-registry",
		RegistryPVCName:             "data-pvc",
		RegistryDataMountPath:       "/var/lib/registry",
		RegistryCertMountPath:       "/certs",
		RegistryAutoSecretMountPath: "/auth",
		RegistryConfMountPath:       "/etc/docker/registry",
		RegistryPVMode:              "Filesystem",
	}
}
