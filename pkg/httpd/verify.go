package httpd

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"sync"
)

var wgVerifyHTTP sync.WaitGroup

func (f *FileServer) RunVerifyHttpd() error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	routeClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()
	wgVerifyHTTP.Add(4)
	go func() {
		res, err := f.verifyDeployment(ctx, *client)
		wgVerifyHTTP.Done()
		if err != nil {
			log.Fatal("[ERROR] verifyDeployment: ", err)
		}
		log.Println("Verify Deployment httpd: ", res)
	}()
	go func() {
		res, err := f.verifyService(ctx, *client)
		wgVerifyHTTP.Done()
		if err != nil {
			log.Fatal("[ERROR] verifyService: ", err)
		}
		log.Println("Verify Service httpd: ", res)
	}()
	go func() {
		res, err := f.verifyRoute(ctx, *routeClient)
		wgVerifyHTTP.Done()
		if err != nil {
			log.Fatal("[ERROR] verifyRoute: ", err)
		}
		log.Println("Verify Route httpd: ", res)
	}()
	go func() {
		res, err := f.verifyPVC(ctx, *client)
		wgVerifyHTTP.Done()
		if err != nil {
			log.Fatal("[ERROR] verifyPVC: ", err)
		}
		log.Println("Verify PVC httpd: ", res)
	}()
	wgVerifyHTTP.Wait()
	return nil
}

func (f *FileServer) verifyDeployment(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	deployment, err := client.AppsV1().Deployments(HTTPD_NAMESPACE).Get(ctx, "nginx", metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] verifying Deployment httpd: %s", err)
		return false, err
	}
	if deployment.Status.AvailableReplicas != 1 {
		log.Printf("[ERROR] verifying Deployment httpd: %s", "Deployment is not ready")
		return false, nil
	}
	return true, nil
}

func (f *FileServer) verifyRoute(ctx context.Context, client routev1.RouteV1Client) (bool, error) {
	route, err := client.Routes(HTTPD_NAMESPACE).Get(ctx, "nginx", metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] verifying Route httpd: %s", err)
		return false, err
	}
	if route.Status.Ingress[0].Host != "httpd-server"+f.Domain {
		log.Printf("[ERROR] verifying Route httpd: %s", "Route is not ready")
		return false, nil
	}
	return true, nil
}

func (f *FileServer) verifyService(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	service, err := client.CoreV1().Services(HTTPD_NAMESPACE).Get(ctx, "nginx", metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] error getting Service: %e", err)
		return false, nil
	}
	if service.Spec.Ports[0].Port != f.Port {
		log.Printf("[ERROR] verifying Service httpd: %s", "Service is not ready")
		return false, nil
	}

	return true, nil
}

func (f *FileServer) verifyPVC(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	pvc, err := client.CoreV1().PersistentVolumeClaims(HTTPD_NAMESPACE).Get(ctx, HTTPD_PVC_NAME, metav1.GetOptions{})
	if err != nil {
		log.Printf("[ERROR] error getting persistent volume: %e", err)
		return false, err
	}
	if pvc.Status.Phase != "Bound" {
		log.Printf("[ERROR] verifying PVC httpd: %s", "PVC is not ready")
		return false, err
	}
	return true, nil
}
