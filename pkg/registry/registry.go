package registry

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	a "github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"time"
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
	RegistryRoute               string
	RegistryOCPReleaseImage     string
	RegistryOCPDestIndexNS      string
	RegistryOLMDestIndexNS      string
	RegistryOLMSourceIndex      string
	RegistrySrcPkg              string
	RegistrySrcPkgFormatted     []string
	RegistryExtraImages         []string
	OcDisCatalog                string
	OcpReleaseFull              string
	RegistryUser                string
	RegistryPass                string
	RegistrySecretName          string
	RegistryConfigMapName       string
	RegistryDeploymentName      string
	RegistryDataMountPath       string
	RegistryCertMountPath       string
	RegistryCertPath            string
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
		RegistryRoute:               "",
		RegistryConfigFile:          "config.yml",
		RegistryOCPDestIndexNS:      "ocp4/openshift4",
		RegistryOLMDestIndexNS:      "olm/redhat-operator-index",
		RegistryOCPReleaseImage:     "quay.io/openshift-release-dev/ocp-release:" + config.Ztp.Config.OcOCPTag,
		RegistryOLMSourceIndex:      "registry.redhat.io/redhat/redhat-operator-index:v" + config.Ztp.Config.OcOCPVersion,
		RegistrySrcPkg:              "kubernetes-nmstate-operator,metallb-operator,ocs-operator,local-storage-operator,advanced-cluster-management",
		RegistrySrcPkgFormatted:     []string{"kubernetes-nmstate-operator", "metallb-operator ocs-operator", "local-storage-operator", "advanced-cluster-management"},
		RegistryExtraImages:         []string{"quay.io/jparrill/registry:3", "registry.access.redhat.com/rhscl/httpd-24-rhel7:latest", "quay.io/ztpfw/ui:latest"},
		RegistryUser:                "dummy",
		RegistryPass:                "dummy123",
		RegistrySecretName:          "auth",
		RegistryConfigMapName:       "registry-conf",
		RegistryDeploymentName:      "kubeframe-registry",
		RegistryServiceName:         "kubeframe-registry",
		RegistryRouteName:           "kubeframe-registry",
		RegistryPVCName:             "data-pvc",
		RegistryDataMountPath:       "/var/lib/registry",
		RegistryCertMountPath:       "/certs",
		RegistryCertPath:            "/etc/pki/ca-trust/source/anchors",
		RegistryAutoSecretMountPath: "/auth",
		RegistryConfMountPath:       "/etc/docker/registry",
		RegistryPVMode:              "Filesystem",
		RegistryCaCertData:          []byte(""),
		RegistryPathCaCert:          "",
	}
}

//Func Login to log into the new registry
func (r *Registry) Login(ctx context.Context) error {
	args := []string{r.RegistryRoute}
	loginOpts := a.LoginOptions{
		AuthFile: r.PullSecretTempFile,
		//CertDir:       r.RegistryPathCaCert,
		CertDir:       r.RegistryCertPath,
		Password:      r.RegistryPass,
		Username:      r.RegistryUser,
		StdinPassword: false,
		GetLoginSet:   false,
		//Verbose:                   false,
		//AcceptRepositories:        true,
		Stdin:                     os.Stdin,
		Stdout:                    os.Stdout,
		AcceptUnspecifiedRegistry: true,
	}
	sysCtx := &types.SystemContext{
		AuthFilePath:                loginOpts.AuthFile,
		DockerCertPath:              loginOpts.CertDir,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	return a.Login(ctx, sysCtx, &loginOpts, args)
}

//Func UpdateTrustCA to update the trust ca in the registry
func (r *Registry) UpdateTrustCA(ctx context.Context, client *kubernetes.Clientset) error {
	res, err := client.CoreV1().Secrets("openshift-ingress").Get(ctx, "router-certs-default", metav1.GetOptions{})
	if err != nil {
		log.Printf(color.InRed("Error getting secret router-certs-defaults: %e"), err)
		return err
	}
	r.RegistryCaCertData = res.Data["tls.crt"]

	if r.Mode == config.MODE_HUB {
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-hub.crt"
	} else {
		//TODO create for more than one spoke
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-" + config.Ztp.Spokes[0].Name + ".crt"
	}

	if err := os.WriteFile(r.RegistryPathCaCert, r.RegistryCaCertData, 0644); err != nil {
		log.Printf(color.InRed("Error writing ca cert to %s: %s"), r.RegistryPathCaCert, err.Error())
		return err
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(r.RegistryPathCaCert)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", r.RegistryPathCaCert, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return fmt.Errorf(color.InRed("No certs appended, using system certs only"))
	}
	log.Println(color.InGreen(">>>> [INFO] Updated trust ca."))
	return nil
}

func (r *Registry) CreateCatalogSource(ctx context.Context) error {
	//TODO create if not exists
	log.Println(color.InYellow(">>>> Creating catalog source."))
	olmclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetOlmAuth()

	catalogSource := &v1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.OcDisCatalog,
			Namespace: r.MarketNS,
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType:  v1alpha1.SourceTypeGrpc,
			Image:       r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion,
			DisplayName: r.OcDisCatalog,
			Publisher:   r.OcDisCatalog,
			UpdateStrategy: &v1alpha1.UpdateStrategy{
				&v1alpha1.RegistryPoll{
					Interval: &metav1.Duration{Duration: time.Minute * 30},
				},
			},
		},
	}
	//create catalog source
	_, err := olmclient.CatalogSources(r.MarketNS).Create(ctx, catalogSource, metav1.CreateOptions{})
	if err != nil {
		log.Printf(color.InRed(">>>> [ERROR] Error creating catalog source: %s"), err.Error())
		return err
	}

	return nil
}

func (r *Registry) GetRegistryRouteName(ctx context.Context, client *routev1.RouteV1Client) (string, error) {
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
