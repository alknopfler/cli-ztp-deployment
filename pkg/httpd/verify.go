package httpd

import (
	"context"
	"errors"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wgVerifyHTTP sync.WaitGroup

//Run Preflight:
// - Check if the conditions are ready or not to be deployed or even verify if the resource is ready
// - Strategy: wait for all to get the error at the end in order to now where is the problem.
func (f *FileServer) RunVerifyHttpd() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	routeClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

	wgVerifyHTTP.Add(4)
	go func() {
		found, err := f.verifyDeployment(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Deployment httpd not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Deployment httpd found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Deployment httpd found and ready"))
		}
		wgVerifyHTTP.Done()
	}()
	go func() {
		found, err := f.verifyService(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Service httpd not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Service httpd found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Service httpd found and ready"))
		}
		wgVerifyHTTP.Done()
	}()
	go func() {
		found, err := f.verifyRoute(ctx, *routeClient)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] Route httpd not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] Route httpd found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] Route httpd found and ready"))
		}
		wgVerifyHTTP.Done()
	}()
	go func() {
		found, err := f.verifyPVC(ctx, *client)
		if !found && err != nil {
			log.Println(color.InRed("[ERROR] PVC httpd not found"))
		} else if found && err != nil {
			log.Println(color.InRed("[ERROR] PVC httpd found but not ready"))
		} else {
			log.Println(color.InGreen("[OK] PVC httpd found and ready"))
		}
		wgVerifyHTTP.Done()
	}()
	wgVerifyHTTP.Wait()
}

func (f *FileServer) verifyDeployment(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	deployment, err := client.AppsV1().Deployments(HTTPD_NAMESPACE).Get(ctx, "nginx", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if deployment.Status.AvailableReplicas < 1 {
		return true, errors.New(color.InRed("[ERROR] Deployment is not ready, replicas available insufficent"))
	}
	return true, nil
}

func (f *FileServer) verifyRoute(ctx context.Context, client routev1.RouteV1Client) (found bool, err error) {
	route, err := client.Routes(HTTPD_NAMESPACE).Get(ctx, "httpd-server-route", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if route.Status.Ingress[0].Host != "httpd-server"+f.Domain {
		return true, errors.New(color.InRed("[ERROR] verifying Route httpd: Route is not ready"))
	}
	return true, nil
}

func (f *FileServer) verifyService(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	service, err := client.CoreV1().Services(HTTPD_NAMESPACE).Get(ctx, "httpd-server-service", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if service.Spec.Ports[0].Port != f.Port {
		return true, errors.New(color.InRed("[ERROR] verifying Service httpd: Service port is not ready"))
	}
	return true, nil
}

func (f *FileServer) verifyPVC(ctx context.Context, client kubernetes.Clientset) (found bool, err error) {
	pvc, err := client.CoreV1().PersistentVolumeClaims(HTTPD_NAMESPACE).Get(ctx, "httpd-pv-claim", metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if pvc.Status.Phase != "Bound" {
		return true, errors.New(color.InRed("[ERROR] verifying PVC httpd: PVC is not bound"))
	}
	return true, nil
}

//TODO change all string to constant
