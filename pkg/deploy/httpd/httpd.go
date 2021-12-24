package httpd

import (
	"context"
	"fmt"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"log"
	"sync"
)

const (
	DEFAULT_MOUNT_PATH         = "/usr/share/nginx/html"
	DEFAULT_SIZE               = "5Gi"
	DEFAULT_PORT               = 80
	DEFAULT_TARGETPORT         = 8080
	INGRESS_CONTROLLER_GROUP   = "operator.openshift.io"
	INGRESS_CONTROLLER_VERSION = "v1"
	INGRESS_CONTROLLER_KIND    = "IngressController"
	INGRESS_CONTROLLER_JQPATH  = ".status.domain"
	INGRESS_CONTROLLER_NS      = "openshift-ingress-operator"
	INGRESS_CONTROLLER_NAME    = "default"
)

//type FileServer
type FileServer struct {
	mountPath  string
	size       string
	domain     string
	port       int
	targetPort int
}

var wg sync.WaitGroup

//Constructor NewFileServer
func NewFileServer(mountPath, size, domain string, port int, targetPort int) *FileServer {
	return &FileServer{
		mountPath:  mountPath,
		size:       size,
		domain:     domain,
		port:       port,
		targetPort: targetPort,
	}
}

func NewFileServerDefault() *FileServer {

	return &FileServer{
		mountPath:  DEFAULT_MOUNT_PATH,
		size:       DEFAULT_SIZE,
		domain:     getDomainFromCluster(),
		port:       DEFAULT_PORT,
		targetPort: DEFAULT_TARGETPORT,
	}
}

func RunHttpd() error {
	getDomainFromCluster()
	return nil
}

func createDeployment() error {

	return nil
}

func createRoute() error {

	return nil
}

func createService() error {

	return nil
}

func createPVC() error {

	return nil
}

func getDomainFromCluster() string {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).SetWithDynamic()
	domain, err := resources.NewGenericGet(ctx, dynamicClient, INGRESS_CONTROLLER_GROUP, INGRESS_CONTROLLER_VERSION, INGRESS_CONTROLLER_KIND, INGRESS_CONTROLLER_NS, INGRESS_CONTROLLER_NAME, INGRESS_CONTROLLER_JQPATH).
		GetResourceByJq()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(domain)

	return ""
}
