package registry

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	a "github.com/containers/common/pkg/auth"
	"github.com/containers/image/v5/types"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"time"
)

//Func Login to log into the new registry
func (r *Registry) Login(ctx context.Context) error {
	args := []string{r.RegistryRoute}
	loginOpts := a.LoginOptions{
		AuthFile:      r.PullSecretTempFile,
		CertDir:       r.RegistryPathCaCert,
		Password:      r.RegistryPass,
		Username:      r.RegistryUser,
		StdinPassword: false,
		GetLoginSet:   false,
		//Verbose:                   false,
		//AcceptRepositories:        true,
		Stdin:                     os.Stdin,
		Stdout:                    os.Stdout,
		AcceptUnspecifiedRegistry: true,
	}
	sysCtx := &types.SystemContext{
		AuthFilePath:                loginOpts.AuthFile,
		DockerCertPath:              loginOpts.CertDir,
		DockerInsecureSkipTLSVerify: types.NewOptionalBool(true),
	}
	return a.Login(ctx, sysCtx, &loginOpts, args)
}

//Func UpdateTrustCA to update the trust ca in the registry
func (r *Registry) UpdateTrustCA(ctx context.Context, client *kubernetes.Clientset) error {
	res, err := client.CoreV1().Secrets("openshift-ingress").Get(ctx, "router-certs-default", metav1.GetOptions{})
	if err != nil {
		log.Printf(color.InRed("Error getting secret router-certs-defaults: %e"), err)
		return err
	}
	r.RegistryCaCertData = res.Data["tls.crt"]

	if r.Mode == config.MODE_HUB {
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-hub.crt"
	} else {
		//TODO create for more than one spoke
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-" + config.Ztp.Spokes[0].Name + ".crt"
	}

	if err := os.WriteFile(r.RegistryPathCaCert, r.RegistryCaCertData, 0644); err != nil {
		log.Printf(color.InRed("Error writing ca cert to %s: %s"), r.RegistryPathCaCert, err.Error())
		return err
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(r.RegistryPathCaCert)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", r.RegistryPathCaCert, err)
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return fmt.Errorf(color.InRed("No certs appended, using system certs only"))
	}
	log.Println(color.InGreen(">>>> [OK] Updated trust ca."))
	return nil
}

func (r *Registry) CreateCatalogSource(ctx context.Context, client *kubernetes.Clientset) error {
	log.Println(color.InYellow(">>>> Creating catalog source."))
	olmclient := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetOlmAuth()

	catalogSource := &v1alpha1.CatalogSource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.OcDisCatalog,
			Namespace: r.MarketNS,
		},
		Spec: v1alpha1.CatalogSourceSpec{
			SourceType:  v1alpha1.SourceTypeGrpc,
			Image:       r.RegistryRoute + "/" + r.RegistryOLMDestIndexNS + ":v" + config.Ztp.Config.OcOCPVersion,
			DisplayName: r.OcDisCatalog,
			Publisher:   r.OcDisCatalog,
			UpdateStrategy: &v1alpha1.UpdateStrategy{
				&v1alpha1.RegistryPoll{
					Interval: &metav1.Duration{Duration: time.Minute * 30},
				},
			},
		},
	}
	//create catalog source
	_, err := olmclient.CatalogSources(r.MarketNS).Create(ctx, catalogSource, metav1.CreateOptions{})
	if err != nil {
		log.Printf(color.InRed("[ERROR] Error creating catalog source: %s"), err.Error())
		return err
	}

	return nil
}
