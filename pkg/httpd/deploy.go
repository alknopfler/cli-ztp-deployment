package httpd

import (
	"context"
	"encoding/json"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
)

func (f *FileServer) RunDeployHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	wg.Add(4)
	go func() {
		err := f.createDeployment(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.createService(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.createRoute(ctx, *client)
		if err != nil {

		}
	}()
	go func() {
		err := f.createPVC(ctx, *client)
		if err != nil {

		}
	}()
	wg.Wait()
	return nil
}

func (f *FileServer) createDeployment(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func (f *FileServer) createRoute(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func (f *FileServer) createService(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func (f *FileServer) createPVC(ctx context.Context, client kubernetes.Clientset) error {
	defer wg.Done()

	return nil
}

func getDomainFromCluster(client dynamic.Interface, ctx context.Context) string {
	d, err := resources.NewGenericGet(ctx, client, INGRESS_CONTROLLER_GROUP, INGRESS_CONTROLLER_VERSION, INGRESS_CONTROLLER_KIND, INGRESS_CONTROLLER_NS, INGRESS_CONTROLLER_NAME, INGRESS_CONTROLLER_JQPATH).GetResourceByJq()
	if err != nil {
		log.Fatalf("[ERROR] Getting resources in GetDomainFromCluster: %e", err)
	}
	b, _ := json.Marshal(d)
	value := jsonvalue.MustUnmarshal(b)
	domain, _ := value.Get("Object", "status", "domain")

	return domain.String()
}
