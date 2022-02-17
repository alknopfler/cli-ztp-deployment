package httpd

import (
	"context"
	"encoding/json"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/dynamic"
	"log"
)

const (
	DEFAULT_MOUNT_PATH         = "/var/www/html"
	DEFAULT_IMAGE_NAME         = "registry.access.redhat.com/rhscl/httpd-24-rhel7:latest"
	DEFAULT_SIZE               = "5Gi"
	DEFAULT_PORT               = 80
	DEFAULT_TARGETPORT         = 8080
	INGRESS_CONTROLLER_GROUP   = "operator.openshift.io"
	INGRESS_CONTROLLER_VERSION = "v1"
	INGRESS_CONTROLLER_KIND    = "ingresscontrollers"
	INGRESS_CONTROLLER_JQPATH  = ".status.domain"
	INGRESS_CONTROLLER_NS      = "openshift-ingress-operator"
	INGRESS_CONTROLLER_NAME    = "default"
	HTTPD_NAMESPACE            = "default"
	HTTPD_VOLUME_NAME          = "httpd-pv-storage"
	HTTPD_PVC_NAME             = "httpd-pv-claim"
	HTTPD_DEPLOYMENT_NAME      = "nginx"
)

//type FileServer
type FileServer struct {
	MountPath  string
	Size       string
	Domain     string
	Port       int32
	TargetPort int32
}

//Constructor NewFileServer
func NewFileServerDefault() *FileServer {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()

	return &FileServer{
		MountPath:  DEFAULT_MOUNT_PATH,
		Size:       DEFAULT_SIZE,
		Domain:     GetDomainFromCluster(client, ctx),
		Port:       DEFAULT_PORT,
		TargetPort: DEFAULT_TARGETPORT,
	}
}

func GetDomainFromCluster(client dynamic.Interface, ctx context.Context) string {
	d, err := resources.NewGenericGet(ctx, client, INGRESS_CONTROLLER_GROUP, INGRESS_CONTROLLER_VERSION, INGRESS_CONTROLLER_KIND, INGRESS_CONTROLLER_NS, INGRESS_CONTROLLER_NAME, INGRESS_CONTROLLER_JQPATH).GetResourceByJq()
	if err != nil {
		log.Printf(color.InRed("[ERROR] Getting resources in GetDomainFromCluster: "), err)
		return "[ERROR GETTING DOMAIN]"
	}
	b, _ := json.Marshal(d)
	value := jsonvalue.MustUnmarshal(b)
	domain, _ := value.Get("Object", "status", "domain")

	return domain.String()
}
