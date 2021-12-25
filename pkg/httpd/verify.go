package httpd

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
)

func (f *FileServer) RunVerifyHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).Set()
	wg.Add(4)
	go func() {
		res, err := f.verifyDeployment(ctx, *client)
		if err != nil {
			log.Fatal("[ERROR] verifyDeployment: ", err)
		}
		log.Println("Verify Deployment httpd: ", res)
	}()
	go func() {
		res, err := f.verifyService(ctx, *client)
		if err != nil {
			log.Fatal("[ERROR] verifyService: ", err)
		}
		log.Println("Verify Service httpd: ", res)
	}()
	go func() {
		res, err := f.verifyRoute(ctx, *client)
		if err != nil {
			log.Fatal("[ERROR] verifyRoute: ", err)
		}
		log.Println("Verify Route httpd: ", res)
	}()
	go func() {
		res, err := f.verifyPVC(ctx, *client)
		if err != nil {
			log.Fatal("[ERROR] verifyPVC: ", err)
		}
		log.Println("Verify PVC httpd: ", res)
	}()

	return nil
}

func (f *FileServer) verifyDeployment(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	defer wg.Done()
	//TODO verify if exist (skip flow) and if not create.
	//

	return true, nil
}

func (f *FileServer) verifyRoute(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	defer wg.Done()

	return true, nil
}

func (f *FileServer) verifyService(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	defer wg.Done()

	return true, nil
}

func (f *FileServer) verifyPVC(ctx context.Context, client kubernetes.Clientset) (bool, error) {
	defer wg.Done()
	pvc, err := client.CoreV1().PersistentVolumeClaims(HTTPD_NAMESPACE).Get(ctx, HTTPD_PVC_NAME, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("[ERROR] error getting persistent volume: %e", err)
		return false, err
	}
	if pvc.Status.Phase != "Bound" {
		log.Fatalf("[ERROR] pvc is not bound: %s", pvc.Status.Phase)
		return false, err
	}
	return true, nil
}
