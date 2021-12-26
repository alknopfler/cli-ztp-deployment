package httpd

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"sync"
)

const (
	DEFAULT_MOUNT_PATH         = "/usr/share/nginx/html"
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
	HTTPD_PVC_NAME             = "httpd-pv-claim"
	HTTPD_DEPLOYMENT_NAME      = "nginx"
)

var wg sync.WaitGroup

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
		Domain:     resources.GetDomainFromCluster(client, ctx),
		Port:       DEFAULT_PORT,
		TargetPort: DEFAULT_TARGETPORT,
	}
}
