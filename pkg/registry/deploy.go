package registry

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/alknopfler/cli-ztp-deployment/config"
	"github.com/alknopfler/cli-ztp-deployment/pkg/auth"
	"github.com/alknopfler/cli-ztp-deployment/pkg/resources"
	apiroutev1 "github.com/openshift/api/route/v1"
	routev1 "github.com/openshift/client-go/route/clientset/versioned/typed/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"sync"
)

var wgDeployRegistry sync.WaitGroup

func (r *Registry) RunDeployRegistry() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//get client from kubeconfig extracted based on Mode (HUB or SPOKE)
	client := auth.NewZTPAuth(config.GetKubeconfigFromMode(r.Mode)).GetAuth()
	dynamicClient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetAuthWithGeneric()
	ocpclient := auth.NewZTPAuth(config.Ztp.Config.KubeconfigHUB).GetRouteAuth()

	//Step 1 - Create the namespace for the registry
	err := r.createNamespace(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating Namespace for the registry: %v"), err)
		return err
	}
	//Step 2 - Create the secret and config map for the registry
	err = r.createSecret(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating secret and config map for the registry: %v"), err)
		return err
	}
	err = r.createConfigMap(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error creating secret and config map for the registry: %v"), err)
		return err
	}
	// Step 3 - Create the rest of the manifests for the registry. We'll use goroutines to do this
	wgDeployRegistry.Add(4)
	go func() error {
		err := r.createDeployment(ctx, client)
		wgDeployRegistry.Done()
		if err != nil {
			log.Fatalf("Error creating deployment: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := r.createService(ctx, client)
		wgDeployRegistry.Done()
		if err != nil {
			log.Fatalf("Error creating service: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := r.createRoute(ctx, *ocpclient, dynamicClient)
		wgDeployRegistry.Done()
		if err != nil {
			log.Fatalf("Error creating route: %v", err)
			return err
		}
		return nil
	}()
	go func() error {
		err := r.createPVC(ctx, *client)
		wgDeployRegistry.Done()
		if err != nil {
			log.Fatalf("Error creating PVC: %v", err)
			return err
		}
		return nil
	}()
	wgDeployRegistry.Wait()

	err = r.updateTrustCA(ctx, client)
	if err != nil {
		log.Printf(color.InRed("Error updating the system CA with the new registry cert to be trusted: %v"), err)
		return err
	}
	return nil
}

func (r *Registry) createNamespace(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyNamespace(ctx, client); !found && err != nil {
		log.Printf(color.InBold(color.InYellow("Namespace %s not found, Creating it...")), r.RegistryNS)
		nsName := &coreV1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistryNS,
			},
		}
		_, err := client.CoreV1().Namespaces().Create(ctx, nsName, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating namespace: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Namespace for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Namespace for the registry already exists. Skipping creation"))
	return nil
}

//Func to create the secret for the registry
func (r *Registry) createSecret(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifySecret(ctx, client); !found && err != nil {
		log.Printf(color.InBold(color.InYellow("Secret for the registry not found, Creating it...")))
		//create secret
		secret := &coreV1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistrySecretName,
			},
			Data: map[string][]byte{
				"htpasswd": []byte(r.RegistrySecretHash),
			},
		}

		_, err := client.CoreV1().Secrets(r.RegistryNS).Create(ctx, secret, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating secret: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Secret for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Secret for the registry already exists. Skipping creation"))
	return nil
}

//Func to create the config map for the registry
func (r *Registry) createConfigMap(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyConfigMap(ctx, client); !found && err != nil {
		log.Printf(color.InBold(color.InYellow("Config Map for the registry not found, Creating it...")))
		//create config map
		configMap := &coreV1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistryConfigMapName,
			},
			Data: map[string]string{
				"config.yml": `version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
compatibility:
  schema1:
    enabled: true

`,
			},
		}
		_, err := client.CoreV1().ConfigMaps(r.RegistryNS).Create(ctx, configMap, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating config map: %v"), err)
			return err
		}
		log.Println(color.InGreen(">>>> Config Map for the registry created successfully"))
		return nil
	}
	log.Println(color.InGreen(">>>> Config Map for the registry already exists. Skipping creation"))
	return nil
}

//Func createDeployment to create the deployment for the registry
func (r *Registry) createDeployment(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyDeployment(ctx, *client); !found && err != nil {
		log.Printf(color.InBold(color.InYellow("Deployment for the registry not found, Creating it...")))
		//create deployment
		deployment := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistryDeploymentName,
				Labels: map[string]string{
					"name": r.RegistryDeploymentName,
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: int32Ptr(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"name": r.RegistryDeploymentName,
					},
				},
				Template: coreV1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"name": r.RegistryDeploymentName,
						},
					},
					Spec: coreV1.PodSpec{
						Containers: []coreV1.Container{
							{
								Name:  r.RegistryDeploymentName,
								Image: r.RegistryExtraImages,
								Ports: []coreV1.ContainerPort{
									{
										Name:          "registry",
										ContainerPort: 5000,
										Protocol:      coreV1.ProtocolTCP,
									},
								},
								VolumeMounts: []coreV1.VolumeMount{
									{
										Name:      "data",
										MountPath: r.RegistryDataMountPath,
									},
									{
										Name:      "certs-secret",
										MountPath: r.RegistryCertMountPath,
										ReadOnly:  true,
									},
									{
										Name:      "auth-secret",
										MountPath: r.RegistryAutoSecretMountPath,
										ReadOnly:  true,
									},
									{
										Name:      "registry-conf",
										MountPath: r.RegistryConfMountPath,
										ReadOnly:  true,
										//SubPath:   r.RegistryConfigFile,
									},
								},
								Env: []coreV1.EnvVar{
									{
										Name:  "REGISTRY_AUTH",
										Value: "htpasswd",
									},
									{
										Name:  "REGISTRY_AUTH_HTPASSWD_REALM",
										Value: "Registry",
									},
									{
										Name:  "REGISTRY_AUTH_HTPASSWD_PATH",
										Value: "/auth/htpasswd",
									},
									{
										Name:  "REGISTRY_HTTP_TLS_CERTIFICATE",
										Value: "/certs/tls.crt",
									},
									{
										Name:  "REGISTRY_HTTP_TLS_KEY",
										Value: "/certs/tls.key",
									},
									{
										Name:  "REGISTRY_HTTP_SECRET",
										Value: "ALongRandomSecretForRegistry",
									},
								},
							},
						},
						Volumes: []coreV1.Volume{
							{
								Name: "data",
								VolumeSource: coreV1.VolumeSource{
									PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
										ClaimName: r.RegistryPVCName},
								},
							},
							{
								Name: "certs-secret",
								VolumeSource: coreV1.VolumeSource{
									Secret: &coreV1.SecretVolumeSource{
										SecretName: "kubeframe-registry-tls",
									},
								},
							},
							{
								Name: "auth-secret",
								VolumeSource: coreV1.VolumeSource{
									Secret: &coreV1.SecretVolumeSource{
										SecretName: "auth",
									},
								},
							},
							{
								Name: "registry-conf",
								VolumeSource: coreV1.VolumeSource{
									ConfigMap: &coreV1.ConfigMapVolumeSource{
										LocalObjectReference: coreV1.LocalObjectReference{
											Name: r.RegistryConfigMapName},
									},
								},
							},
						},
					},
				},
			},
		}
		res, err := client.AppsV1().Deployments(r.RegistryNS).Create(ctx, deployment, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating deployment: %e"), err)
			return err
		}
		err = resources.WaitForDeployment(ctx, res, client)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for deployment: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created deployment %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> Deployment Registry already exists. Skipping creation."))
	return nil
}

//Func createService to create the service for the registry
func (r *Registry) createService(ctx context.Context, client *kubernetes.Clientset) error {
	if found, err := r.verifyService(ctx, *client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Service for the registry not found. Creating Service Registry")))
		serviceSpec := &coreV1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: r.RegistryServiceName,
				Labels: map[string]string{
					"name": r.RegistryServiceName,
				},
				Annotations: map[string]string{
					"service.beta.openshift.io/serving-cert-secret-name": "kubeframe-registry-tls",
				},
			},
			Spec: coreV1.ServiceSpec{
				Selector: map[string]string{
					"name": r.RegistryServiceName,
				},
				Type: "ClusterIP",
				Ports: []coreV1.ServicePort{
					{
						Name:     "registry",
						Port:     443,
						Protocol: coreV1.ProtocolTCP,
						TargetPort: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: 5000,
						},
					},
				},
				SessionAffinity: "None",
			},
		}

		res, err := client.CoreV1().Services(r.RegistryNS).Create(ctx, serviceSpec, metav1.CreateOptions{})
		if err != nil {
			log.Printf(color.InRed("Error creating Service: %e"), err)
			return err
		}
		err = resources.WaitForService(ctx, client, res)
		if err != nil {
			log.Printf(color.InRed("[ERROR] waiting for Service: %s"), err)
			return err
		}
		log.Printf(color.InGreen(">>>> Created Service %s\n"), res.GetObjectMeta().GetName())
		return nil
	}
	// Already created and return nil
	log.Printf(color.InGreen(">>>> Service Registry already exists. Skipping creation."))
	return nil
}

//Func createRoute to create a route for the registry
func (r *Registry) createRoute(ctx context.Context, client routev1.RouteV1Client, dynamicclient dynamic.Interface) error {
	if found, err := r.verifyRoute(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Route not found. Creating the Registry route")))
		route := apiroutev1.Route{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"name": r.RegistryRouteName,
				},
				Name: r.RegistryRouteName,
			},
			Spec: apiroutev1.RouteSpec{
				Port: &apiroutev1.RoutePort{
					TargetPort: intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "registry",
					},
				},
				To: apiroutev1.RouteTargetReference{
					Name: r.RegistryRouteName,
				},
			},
		}
		res, err := client.Routes(r.RegistryNS).Create(ctx, &route, metav1.CreateOptions{})
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
	log.Printf(color.InGreen(">>>> Route for Registry already exists. Skipping creation."))
	return nil
}

func (r *Registry) createPVC(ctx context.Context, client kubernetes.Clientset) error {
	if found, err := r.verifyPVC(ctx, client); err != nil && !found {
		log.Println(color.InBold(color.InYellow(">>>> Creating PVC Registry")))
		fs := coreV1.PersistentVolumeFilesystem
		pvc := &apiv1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      r.RegistryPVCName,
				Namespace: r.RegistryNS,
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceName(apiv1.ResourceStorage): resource.MustParse("100Gi"),
					},
				},
				VolumeMode: &fs,
			},
		}
		res, err := client.CoreV1().PersistentVolumeClaims(r.RegistryNS).Create(ctx, pvc, metav1.CreateOptions{})
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
	log.Printf(color.InGreen(">>>> PVC for registry already exists. Skipping creation."))
	return nil
}

//Func updateTrustCA to update the trust ca in the registry
func (r *Registry) updateTrustCA(ctx context.Context, client *kubernetes.Clientset) error {
	res, err := client.CoreV1().Secrets("openshift-ingress").Get(ctx, "router-certs-default", metav1.GetOptions{})
	if err != nil {
		log.Printf(color.InRed("Error getting secret router-certs-defaults: %e"), err)
		return err
	}
	r.RegistryCaCertData = res.Data["tls.crt"]

	var pathCaCert string
	if r.Mode == "Hub" {
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-hub.crt"
	} else {
		//TODO create for more than one spoke
		r.RegistryPathCaCert = "/etc/pki/ca-trust/source/anchors/internal-registry-" + config.Ztp.Spokes[0].Name + ".crt"
	}
	if err := os.WriteFile(pathCaCert, r.RegistryCaCertData, 0644); err != nil {
		log.Printf(color.InRed("Error writing ca cert to %s: %e"), pathCaCert, err)
		return err
	}

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if ok := rootCAs.AppendCertsFromPEM(r.RegistryCaCertData); !ok {
		return fmt.Errorf(color.InRed("No certs appended, using system certs only"))
	}
	log.Println(color.InGreen(">>>> [OK] Updated trust ca."))
	return nil
}

func (r *Registry) createMachineConfig(ctx context.Context, client *auth.ZTPAuth) error {
	machineConfigGVR := schema.GroupVersionResource{
		Group:    "machineconfiguration.openshift.io",
		Version:  "v1",
		Resource: "MachineConfig",
	}

	machineConfigSpec := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "MachineConfig",
			"apiVersion": "operators.coreos.com/v1",
			"metadata": map[string]interface{}{
				"name": "update-localregistry-ca-certs",
				"labels": map[string]interface{}{
					"machineconfiguration.openshift.io/role": "master",
				},
			},
			"spec": map[string]interface{}{
				"config": map[string]interface{}{
					"ignition": map[string]string{
						"version": "3.1.0",
					},
					"storage": map[string]interface{}{
						"files": []map[string]interface{}{
							{
								"path": r.RegistryPathCaCert,
								"mode": "0493",
								"contents": map[string]interface{}{
									"source": "data:text/plain;charset=us-ascii;base64," + string(r.RegistryCaCertData[:]),
								},
							},
						},
					},
				},
			},
		},
	}
	res, err := client.GetAuthWithGeneric().Resource(machineConfigGVR).Namespace(r.RegistryNS).Create(ctx, machineConfigSpec, metav1.CreateOptions{})
	if err != nil {
		log.Printf(color.InRed("Error creating MachineConfig: %e"), err)
		return err
	}
	log.Printf(color.InGreen(">>>> Created MachineConfig %s\n"), res.GetName())
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
