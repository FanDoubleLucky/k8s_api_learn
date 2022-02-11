package main

import (
	"context"
	"encoding/json"
	"fmt"
	"istio.io/api/meta/v1alpha1"
	alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cv1 "k8s.io/api/core/v1"

	//bv1 "k8s.io/api/batch/v1"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/client-go/tools/clientcmd"

	istioclient "istio.io/client-go/pkg/clientset/versioned"
)

func init() {
	clientgoscheme.AddToScheme(scheme)
}

var (
	scheme = runtime.NewScheme()
	log    = ctrl.Log.WithName("collector")
)
func NewManager() ctrl.Manager {
	options := ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":0",
		LeaderElection:     false,
		Port:               9443,
		LeaderElectionID:   "study-k8s-api",
	}

	manager, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		log.Error(err, "Get Manager error")
		panic(err)
	}
	return manager
}

var reader client.Reader
var creator client.Client

func InitKubeConn(){
	m := NewManager()
	reader = m.GetAPIReader()
	creator = m.GetClient()
}

func main() {
	InitKubeConn()
	http.HandleFunc("/listPod", ListPod)
	http.HandleFunc("/createDeploy", CreateDeploy)
	port := "8111"
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err)
}

type PodResponse struct {
	Name string
	Namespace string
	Kind string
}
func ListPod(writer http.ResponseWriter, request *http.Request){
	listOption := client.ListOptions{Namespace: "ui-app"}
	var podList cv1.PodList
	if err := reader.List(context.TODO(), &podList, &listOption); err != nil {
		fmt.Println("List Err", err)
		http.Error(writer, err.Error(), 500)
		return
	}
	var result []PodResponse
	for _, pod := range podList.Items {
		var podResponse PodResponse
		podResponse.Name = pod.Name
		podResponse.Namespace = pod.Namespace
		result = append(result, podResponse)
	}
	data, _ := json.Marshal(result)
	writer.Write(data)
	return
}

func CreateDeploy(writer http.ResponseWriter, request *http.Request){
	var repliccount int32 = 2
	deploymentToCreate := v1.Deployment{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "centos-deploy",
			//GenerateName:               "",
			Namespace:                  "ui-app",
			//SelfLink:                   "",
			//UID:                        "",
			//ResourceVersion:            "",
			//Generation:                 0,
			//CreationTimestamp:          metav1.Time{},
			//DeletionTimestamp:          nil,
			//DeletionGracePeriodSeconds: nil,
			Labels: 					map[string]string{
				"app": "centos-network-tool",
			},
			//Annotations:                nil,
			//OwnerReferences:            nil,
			//Finalizers:                 nil,
			//ClusterName:                "",
			//ManagedFields:              nil,
		},
		Spec:       v1.DeploymentSpec{
			Replicas:                &repliccount,
			Selector:                &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "centos-network-tool"},
				//MatchExpressions: nil,
			},
			Template:                cv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:                       "centos",
					//GenerateName:               "",
					Namespace:                  "ui-app",
					//SelfLink:                   "",
					//UID:                        "",
					//ResourceVersion:            "",
					//Generation:                 0,
					//CreationTimestamp:          metav1.Time{},
					//DeletionTimestamp:          nil,
					//DeletionGracePeriodSeconds: nil,
					Labels: 					map[string]string{
						"app": "centos-network-tool",
						"app.soa.freewheel.tv/tiller": "api",
					},
					//Annotations:                nil,
					//OwnerReferences:            nil,
					//Finalizers:                 nil,
					//ClusterName:                "",
					//ManagedFields:              nil,
				},
				Spec:       cv1.PodSpec{
					//Volumes:                       nil,
					//InitContainers:                nil,
					Containers:                    []cv1.Container{
						{
							Name:                     "server",
							Image:                    "centos",
							Command:                  []string{"sh", "-c", "echo \"Hello, Kubernetes!\" && sleep 36000"},
							//Args:                     nil,
							//WorkingDir:               "",
							//Ports:                    nil,
							//EnvFrom:                  nil,
							//Env:                      nil,
							//Resources:                cv1.ResourceRequirements{},
							//VolumeMounts:             nil,
							//VolumeDevices:            nil,
							//LivenessProbe:            nil,
							//ReadinessProbe:           nil,
							//StartupProbe:             nil,
							//Lifecycle:                nil,
							//TerminationMessagePath:   "",
							//TerminationMessagePolicy: "",
							ImagePullPolicy:          "IfNotPresent",
							//SecurityContext:          nil,
							//Stdin:                    false,
							//StdinOnce:                false,
							TTY:                      false,
						},
					},
					//EphemeralContainers:           nil,
					//RestartPolicy:                 "",
					//TerminationGracePeriodSeconds: nil,
					//ActiveDeadlineSeconds:         nil,
					//DNSPolicy:                     "",
					//NodeSelector:                  nil,
					//ServiceAccountName:            "",
					//DeprecatedServiceAccount:      "",
					//AutomountServiceAccountToken:  nil,
					//NodeName:                      "",
					//HostNetwork:                   false,
					//HostPID:                       false,
					//HostIPC:                       false,
					//ShareProcessNamespace:         nil,
					//SecurityContext:               nil,
					//ImagePullSecrets:              nil,
					//Hostname:                      "",
					//Subdomain:                     "",
					//Affinity:                      nil,
					//SchedulerName:                 "",
					//Tolerations:                   nil,
					//HostAliases:                   nil,
					//PriorityClassName:             "",
					//Priority:                      nil,
					//DNSConfig:                     nil,
					//ReadinessGates:                nil,
					//RuntimeClassName:              nil,
					//EnableServiceLinks:            nil,
					//PreemptionPolicy:              nil,
					//Overhead:                      nil,
					//TopologySpreadConstraints:     nil,
					//SetHostnameAsFQDN:             nil,
				},
			},
			Strategy:                v1.DeploymentStrategy{
				Type:          "RollingUpdate",
				//RollingUpdate: nil,
			},
			//MinReadySeconds:         0,
			//RevisionHistoryLimit:    nil,
			//Paused:                  false,
			//ProgressDeadlineSeconds: nil,
		},
		//Status:     v1.DeploymentStatus{},
	}

	createOption := client.CreateOptions{}

	err := creator.Create(context.TODO(), &deploymentToCreate, &createOption)

	if err != nil {
		fmt.Println("Create Err", err)
		http.Error(writer, err.Error(), 500)
		return
	}
}

func test() {
	m := NewManager()
	//cli := m.GetClient()
	reader := m.GetAPIReader()
	createClient := m.GetClient()

	kubeconfig := os.Getenv("KUBECONFIG")
	restConfig, istioErr := clientcmd.BuildConfigFromFlags("", kubeconfig)

	listOption := client.ListOptions{Namespace: "ui-app"}
	var podList cv1.PodList

	istioClientInstance, istioErr := istioclient.NewForConfig(restConfig)
	vsController := istioClientInstance.NetworkingV1alpha3().VirtualServices("ui-app")
	vsList, istioErr := vsController.List(context.TODO(), metav1.ListOptions{})
	if istioErr != nil {
		fmt.Println(istioErr)
		return
	}
	for _, vs := range vsList.Items {
		fmt.Println("vs:", vs.Name)
	}

	testVS := v1alpha3.VirtualService{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "",
			APIVersion: "",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "",
			//GenerateName:               "",
			Namespace:                  "",
			//SelfLink:                   "",
			//UID:                        "",
			//ResourceVersion:            "",
			//Generation:                 0,
			//CreationTimestamp:          metav1.Time{},
			//DeletionTimestamp:          nil,
			//DeletionGracePeriodSeconds: nil,
			Labels:						map[string]string{},
			//Annotations:                nil,
			//OwnerReferences:            nil,
			//Finalizers:                 nil,
			//ClusterName:                "",
			//ManagedFields:              nil,
		},
		Spec:       alpha3.VirtualService{
			Hosts:                []string{},
			Gateways:             []string{},
			Http:                 []*alpha3.HTTPRoute{{
				Name:                 "",
				Match:                nil,
				Route:                nil,
				Redirect:             nil,
				Delegate:             nil,
				Rewrite:              nil,
				Timeout:              nil,
				Retries:              nil,
				Fault:                nil,
				//Mirror:               nil,
				//MirrorPercent:        nil,
				//MirrorPercentage:     nil,
				CorsPolicy:           nil,
				Headers:              nil,
				//XXX_NoUnkeyedLiteral: struct{}{},
				//XXX_unrecognized:     nil,
				//XXX_sizecache:        0,
			}},
			Tls:                  nil,
			Tcp:                  nil,
			//ExportTo:             nil,
			//XXX_NoUnkeyedLiteral: struct{}{},
			//XXX_unrecognized:     nil,
			//XXX_sizecache:        0,
		},
		Status:     v1alpha1.IstioStatus{
			//Conditions:           nil,
			//ValidationMessages:   nil,
			//ObservedGeneration:   0,
			//XXX_NoUnkeyedLiteral: struct{}{},
			//XXX_unrecognized:     nil,
			//XXX_sizecache:        0,
		},
	}

	vsController.Create(context.TODO(), &testVS, metav1.CreateOptions{})

	err := reader.List(context.TODO(), &podList, &listOption)

	var repliccount int32 = 2
	deploymentToCreate := v1.Deployment{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "centos-deploy",
			//GenerateName:               "",
			Namespace:                  "ui-app",
			//SelfLink:                   "",
			//UID:                        "",
			//ResourceVersion:            "",
			//Generation:                 0,
			//CreationTimestamp:          metav1.Time{},
			//DeletionTimestamp:          nil,
			//DeletionGracePeriodSeconds: nil,
			Labels: 					map[string]string{
											"app": "centos-network-tool",
										},
			//Annotations:                nil,
			//OwnerReferences:            nil,
			//Finalizers:                 nil,
			//ClusterName:                "",
			//ManagedFields:              nil,
		},
		Spec:       v1.DeploymentSpec{
			Replicas:                &repliccount,
			Selector:                &metav1.LabelSelector{
										MatchLabels: map[string]string{"app": "centos-network-tool"},
										//MatchExpressions: nil,
									 },
			Template:                cv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:                       "centos",
					//GenerateName:               "",
					Namespace:                  "ui-app",
					//SelfLink:                   "",
					//UID:                        "",
					//ResourceVersion:            "",
					//Generation:                 0,
					//CreationTimestamp:          metav1.Time{},
					//DeletionTimestamp:          nil,
					//DeletionGracePeriodSeconds: nil,
					Labels: 					map[string]string{
													"app": "centos-network-tool",
													"app.soa.freewheel.tv/tiller": "api",
												},
					//Annotations:                nil,
					//OwnerReferences:            nil,
					//Finalizers:                 nil,
					//ClusterName:                "",
					//ManagedFields:              nil,
				},
				Spec:       cv1.PodSpec{
					//Volumes:                       nil,
					//InitContainers:                nil,
					Containers:                    []cv1.Container{
						{
							Name:                     "server",
							Image:                    "centos",
							Command:                  []string{"sh", "-c", "echo \"Hello, Kubernetes!\" && sleep 36000"},
							//Args:                     nil,
							//WorkingDir:               "",
							//Ports:                    nil,
							//EnvFrom:                  nil,
							//Env:                      nil,
							//Resources:                cv1.ResourceRequirements{},
							//VolumeMounts:             nil,
							//VolumeDevices:            nil,
							//LivenessProbe:            nil,
							//ReadinessProbe:           nil,
							//StartupProbe:             nil,
							//Lifecycle:                nil,
							//TerminationMessagePath:   "",
							//TerminationMessagePolicy: "",
							ImagePullPolicy:          "IfNotPresent",
							//SecurityContext:          nil,
							//Stdin:                    false,
							//StdinOnce:                false,
							TTY:                      false,
						},
					},
					//EphemeralContainers:           nil,
					//RestartPolicy:                 "",
					//TerminationGracePeriodSeconds: nil,
					//ActiveDeadlineSeconds:         nil,
					//DNSPolicy:                     "",
					//NodeSelector:                  nil,
					//ServiceAccountName:            "",
					//DeprecatedServiceAccount:      "",
					//AutomountServiceAccountToken:  nil,
					//NodeName:                      "",
					//HostNetwork:                   false,
					//HostPID:                       false,
					//HostIPC:                       false,
					//ShareProcessNamespace:         nil,
					//SecurityContext:               nil,
					//ImagePullSecrets:              nil,
					//Hostname:                      "",
					//Subdomain:                     "",
					//Affinity:                      nil,
					//SchedulerName:                 "",
					//Tolerations:                   nil,
					//HostAliases:                   nil,
					//PriorityClassName:             "",
					//Priority:                      nil,
					//DNSConfig:                     nil,
					//ReadinessGates:                nil,
					//RuntimeClassName:              nil,
					//EnableServiceLinks:            nil,
					//PreemptionPolicy:              nil,
					//Overhead:                      nil,
					//TopologySpreadConstraints:     nil,
					//SetHostnameAsFQDN:             nil,
				},
			},
			Strategy:                v1.DeploymentStrategy{
				Type:          "RollingUpdate",
				//RollingUpdate: nil,
			},
			//MinReadySeconds:         0,
			//RevisionHistoryLimit:    nil,
			//Paused:                  false,
			//ProgressDeadlineSeconds: nil,
		},
		//Status:     v1.DeploymentStatus{},
	}

	createOption := client.CreateOptions{}

	err = createClient.Create(context.TODO(), &deploymentToCreate, &createOption)

	if err != nil {
		fmt.Println(err)
		return
	}
	for _, pod := range podList.Items {
		fmt.Println(pod.Name)
	}
}
