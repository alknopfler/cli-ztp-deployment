package httpd

import (
	"context"
	"encoding/json"
	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	apiroutev1 "github.com/openshift/api/route/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
)

func (f *FileServer) RunDeployHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

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
		err := f.createRoute(ctx, *ocpclient, dynamicClient)
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
	if _, err := f.verifyDeployment(ctx, client); err != nil {
		log.Println(">>>> Creating deployment HTTPD")
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nginx",
				Labels: map[string]string{
					"app": "nginx",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(2),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: apiv1.PodSpec{
						Volumes: []apiv1.Volume{
							{
								Name: "httpd-pv-storage",
								VolumeSource: apiv1.VolumeSource{
									PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
										ClaimName: "httpd-pv-claim",
									},
								},
							},
						},
						Containers: []apiv1.Container{
							{
								Name:  "nginx",
								Image: "quay.io/openshift-scale/nginx:latest",
								Ports: []apiv1.ContainerPort{
									{
										ContainerPort: 8080,
									},
								},
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name:      "httpd-pv-storage",
										MountPath: "/usr/share/nginx/html",
									},
								},
							},
						},
					},
				},
			},
		}
		res, err := client.AppsV1().Deployments(HTTPD_NAMESPACE).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Error creating deployment: %e", err)
			return err
		}
		err = resources.WaitForDeployment(ctx, res, &client)
		if err != nil {
			log.Printf("[ERROR] waiting for deployment: %s", err)
			return err
		}
		log.Printf(">>>> Created deployment %s\n", res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(">>>> Deployment HTTPD already exists. Skipping creation.")
	return nil
}

func (f *FileServer) createRoute(ctx context.Context, client routev1.RouteV1Client, dynamicclient dynamic.Interface) error {
	defer wg.Done()
	if _, err := f.verifyRoute(ctx, client); err != nil {
		log.Println(">>>> Creating route HTTPD")
		route := apiroutev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": "nginx",
				},
				Name:      "httpd-server-route",
				Namespace: HTTPD_NAMESPACE,
			},
			Spec: apiroutev1.RouteSpec{
				Host: "httpd-server" + getDomainFromCluster(dynamicclient, ctx),
				Port: &apiroutev1.RoutePort{
					TargetPort: intstr.IntOrString{
						Type:   DEFAULT_TARGETPORT,
						IntVal: DEFAULT_TARGETPORT,
						StrVal: "8080",
					},
				},
				To: apiroutev1.RouteTargetReference{
					Kind:   "Service",
					Name:   "httpd-server-service",
					Weight: nil,
				},
				WildcardPolicy: "None",
			},
		}
		res, err := client.Routes(HTTPD_NAMESPACE).Create(ctx, &route, metav1.CreateOptions{})
		if err != nil {
			log.Printf("Error creating route: %e", err)
			return err
		}
		err = resources.WaitForRoute(ctx, &client, res)
		if err != nil {
			log.Printf("[ERROR] waiting for route: %s", err)
			return err
		}
		log.Printf(">>>> Created route %s\n", res.GetObjectMeta().GetName())
		return nil
	}

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

func int32Ptr(i int32) *int32 { return &i }
