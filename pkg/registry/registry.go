package registry

import (
	"context"
	"fmt"
	"github.com/alknopfler/cli-ztp-deployment/config"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
)

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
	RegistryCaCertData          []byte
	RegistryPathCaCert          string
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

func (r *Registry) getPullSecretBase(ctx context.Context, client *kubernetes.Clientset) (string, error) {
	res, err := client.CoreV1().Secrets("openshift-config").Get(ctx, "pull-secret", metav1.GetOptions{})
	if err != nil {
		fmt.Println("entra1 error")
		return "", err
	}
	fmt.Println(string(res.Data[".dockerconfigjson"]))
	return string(res.Data[".dockerconfigjson"]), nil
}

//Func to write the content of string to a temporal file
func (r *Registry) writePullSecretBaseToTempFile(string) error {
	file, err := ioutil.TempFile("/tmp", "pull-secret-temp.json")
	if err != nil {
		fmt.Println(err)
	}
	// We can choose to have these files deleted on program close
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte("hello world\n")); err != nil {
		fmt.Println(err)
	}

}
