package httpd

import (
	"context"
	"fmt"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/dynamic"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).SetWithDynamic()
	return &FileServer{
		mountPath:  DEFAULT_MOUNT_PATH,
		size:       DEFAULT_SIZE,
		domain:     getDomainFromCluster(client, ctx),
		port:       DEFAULT_PORT,
		targetPort: DEFAULT_TARGETPORT,
	}
}

func RunHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).SetWithDynamic()

	getDomainFromCluster(client, ctx)
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

func getDomainFromCluster(client dynamic.Interface, ctx context.Context) string {

	domain, err := resources.NewGenericList(ctx, client, INGRESS_CONTROLLER_GROUP, INGRESS_CONTROLLER_VERSION, INGRESS_CONTROLLER_VERSION, INGRESS_CONTROLLER_NS, ".").GetResourcesDynamically()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(domain)

	return ""
}
