package httpd

import (
	"context"
	"github.com/TwiN/go-color"
	"k8s.io/apimachinery/pkg/api/resource"
	"sync"

	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	apiroutev1 "github.com/openshift/api/route/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
)

var wgDeployHTTPD sync.WaitGroup

//Run deploy HTTPD:
// - Launch all resources to be created at the same time
// - Strategy: First fail (if any resource fails to be created)
func (f *FileServer) RunDeployHttpd() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

	wgDeployHTTPD.Add(4)
	go func() error {
		err := f.createDeployment(ctx, *client)
		wgDeployHTTPD.Done()
		if err != nil {
			log.Fatalf("Error creating deployment: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := f.createService(ctx, *client)
		wgDeployHTTPD.Done()
		if err != nil {
			log.Fatalf("Error creating service: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := f.createRoute(ctx, *ocpclient, dynamicClient)
		wgDeployHTTPD.Done()
		if err != nil {
			log.Fatalf("Error creating route: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := f.createPVC(ctx, *client)
		wgDeployHTTPD.Done()
		if err != nil {
			log.Fatalf("Error creating PVC: %v", err)
			return err
		}
		return nil
	}()
	wgDeployHTTPD.Wait()
	return nil
}

//Func to create deployment
func (f *FileServer) createDeployment(ctx context.Context, client kubernetes.Clientset) error {
	if found, err := f.verifyDeployment(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Creating Deployment HTTPD")))
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: HTTPD_DEPLOYMENT_NAME,
				Labels: map[string]string{
					"app": HTTPD_DEPLOYMENT_NAME,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": HTTPD_DEPLOYMENT_NAME,
					},
				},
				Template: apiv1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": HTTPD_DEPLOYMENT_NAME,
						},
					},
					Spec: apiv1.PodSpec{
						Volumes: []apiv1.Volume{
							{
								Name: HTTPD_VOLUME_NAME,
								VolumeSource: apiv1.VolumeSource{
									PersistentVolumeClaim: &apiv1.PersistentVolumeClaimVolumeSource{
										ClaimName: HTTPD_PVC_NAME,
									},
								},
							},
						},
						Containers: []apiv1.Container{
							{
								Name:  "httpd",
								Image: DEFAULT_IMAGE_NAME,
								Ports: []apiv1.ContainerPort{
									{
										ContainerPort: 8080,
									},
								},
								VolumeMounts: []apiv1.VolumeMount{
									{
										Name:      HTTPD_VOLUME_NAME,
										MountPath: DEFAULT_MOUNT_PATH,
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
			log.Printf(color.InRed("Error creating deployment: %e"), err)
			return err
		}
		err = resources.WaitForDeployment(ctx, res, &client)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for deployment: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created deployment %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> Deployment HTTPD already exists. Skipping creation."))
	return nil
}

//Func to create a Route
func (f *FileServer) createRoute(ctx context.Context, client routev1.RouteV1Client, dynamicclient dynamic.Interface) error {
	if found, err := f.verifyRoute(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Creating route HTTPD")))
		route := apiroutev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"app": "httpd",
				},
				Name:      "httpd-server-route",
				Namespace: HTTPD_NAMESPACE,
			},
			Spec: apiroutev1.RouteSpec{
				Host: "httpd-server" + GetDomainFromCluster(dynamicclient, ctx),
				Port: &apiroutev1.RoutePort{
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 8080,
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
			log.Printf(color.InRed("Error creating route: %e"), err)
			return err
		}
		err = resources.WaitForRoute(ctx, &client, res)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for route: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created route %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> Route for HTTPD already exists. Skipping creation."))
	return nil
}

//Func to create a Service
func (f *FileServer) createService(ctx context.Context, client kubernetes.Clientset) error {
	if found, err := f.verifyService(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Creating Service HTTPD")))
		serviceSpec := &coreV1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: "httpd-server-service",
			},
			Spec: coreV1.ServiceSpec{
				Selector: map[string]string{
					"app": "httpd",
				},
				Type: "ClusterIP",
				Ports: []coreV1.ServicePort{
					{
						Port:     DEFAULT_PORT,
						Protocol: coreV1.ProtocolTCP,
						TargetPort: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: DEFAULT_TARGETPORT,
						},
					},
				},
			},
		}

		res, err := client.CoreV1().Services(HTTPD_NAMESPACE).Create(ctx, serviceSpec, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating Service: %e"), err)
			return err
		}
		err = resources.WaitForService(ctx, &client, res)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for Service: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created Service %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> Service  HTTPD already exists. Skipping creation."))
	return nil
}

func (f *FileServer) createPVC(ctx context.Context, client kubernetes.Clientset) error {
	if found, err := f.verifyPVC(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Creating PVC HTTPD")))
		pvc := &apiv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name: HTTPD_PVC_NAME,
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse(DEFAULT_SIZE),
					},
				},
			},
		}
		res, err := client.CoreV1().PersistentVolumeClaims(HTTPD_NAMESPACE).Create(ctx, pvc, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating Pvc: %e"), err)
			return err
		}
		err = resources.WaitForPVC(ctx, &client, res)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for Pvc: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created Pvc %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> PVC for  HTTPD already exists. Skipping creation."))
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

//TODO change all string to constant
