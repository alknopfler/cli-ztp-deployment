package httpd

import (
	"sync"
)

const (
	DEFAULT_MOUNT_PATH = "/usr/share/nginx/html"
	DEFAULT_SIZE       = "5Gi"
	DEFAULT_PORT       = 80
	DEFAULT_TARGETPORT = 8080
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

func RunHttpd() error {

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
