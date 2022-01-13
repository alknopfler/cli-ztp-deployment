package registry

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//type FileServer
type Registry struct {
	Mode                        string
	PullSecretTempFile          string
	PullSecretName              string
	PullSecretNS                string
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
	RegistryCaCertData          []byte
	RegistryPathCaCert          string
}

//Constructor NewFileServer
func NewRegistry(mode string) *Registry {
	return &Registry{
		Mode:                        mode,
		PullSecretTempFile:          "/tmp/pull-secret-temp.json",
		PullSecretName:              "pull-secret",
		PullSecretNS:                "openshift-config",
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
		RegistryCaCertData:          []byte(""),
		RegistryPathCaCert:          "",
	}
}

func (r *Registry) getRegistryRouteName(ctx context.Context, client *routev1.RouteV1Client) (string, error) {
	route, err := client.Routes(r.RegistryNS).Get(ctx, r.RegistryRouteName, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return route.Status.Ingress[0].Host, nil
}

func (r *Registry) GetPullSecretBase() string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	res, err := client.CoreV1().Secrets(r.PullSecretNS).Get(ctx, r.PullSecretName, metav1.GetOptions{})
	if err != nil {
		return ""
	}
	return string(res.Data[".dockerconfigjson"])
}

//Func to write the content of string to a temporal file
func (r *Registry) WritePullSecretBaseToTempFile(data string) error {
	err := ioutil.WriteFile(r.PullSecretTempFile, []byte(data), 0644)
	if err != nil {
		return err
	}
	// Defer done in the cmd cobra command in order to be available during the cmd execution and remove after program closed
	//defer os.Remove("/tmp/pull-secret-temp.json")
	return nil
}
