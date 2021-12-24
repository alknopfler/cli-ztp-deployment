package httpd

import (
	"context"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"k8s.io/client-go/kubernetes"
)

func (f *FileServer) RunVerifyHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).Set()
	wg.Add(4)
	go func() {
		err := f.verifyDeployment(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.verifyService(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.verifyRoute(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.verifyPVC(ctx, *client)
		if err != nil {

		}
	}()

	return nil
}

func (f *FileServer) verifyDeployment(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()
	//TODO verify if exist (skip flow) and if not create.
	//

	return nil
}

func (f *FileServer) verifyRoute(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func (f *FileServer) verifyService(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func (f *FileServer) verifyPVC(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}
