package k8s_utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	k8s_structs "github.com/chujieyang/commonops/ops/forms/k8s"
	"github.com/chujieyang/commonops/ops/models"
	"github.com/chujieyang/commonops/ops/prometheus"
	"github.com/ghodss/yaml"
	apiAppsV1 "k8s.io/api/apps/v1"
	apiCoreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	appsV1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	coreV1 "k8s.io/client-go/kubernetes/typed/core/v1"
	extensionsV1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sClientSet struct {
	Client *kubernetes.Clientset
}

func InitK8sClient(clusterId string) (client K8sClientSet, err error) {
	filePath := models.GetKubeConfigPathByClusterId(clusterId)
	kubeConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}
	kubeRestConf, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		return
	}
	client.Client, err = kubernetes.NewForConfig(kubeRestConf)
	if err != nil {
		return
	}
	return
}

func (clientSet K8sClientSet) GetK8sCoreV1() (v1Interface coreV1.CoreV1Interface) {
	v1Interface = clientSet.Client.CoreV1()
	return
}

func (clientSet K8sClientSet) GetK8sAppsV1() (v1Interface appsV1.AppsV1Interface) {
	v1Interface = clientSet.Client.AppsV1()
	return
}

func (clientSet K8sClientSet) GetK8sExtensionsV1beta1() (v1Interface extensionsV1beta1.ExtensionsV1beta1Interface) {
	v1Interface = clientSet.Client.ExtensionsV1beta1()
	return
}

// 获取 namespaces
func (clientSet K8sClientSet) GetNamespaces() (namespaceList *apiCoreV1.NamespaceList, err error) {
	namespaceList, err = clientSet.GetK8sCoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 新建 namespaces
func (clientSet K8sClientSet) CreateNamespaces(name string) (namespace *apiCoreV1.Namespace, err error) {
	namespace, err = clientSet.GetK8sCoreV1().Namespaces().Create(&apiCoreV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 删除 namespaces
func (clientSet K8sClientSet) DeleteNamespaces(name string) (err error) {
	err = clientSet.GetK8sCoreV1().Namespaces().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 获取 pod
func (clientSet K8sClientSet) GetPods(namespace string) (podList *apiCoreV1.PodList, err error) {
	podList, err = clientSet.GetK8sCoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取 pod 中容器的日志
func (clientSet K8sClientSet) GetPodContainerLogs(namespace string, podName string, containerName string) (logData string, err error) {
	tailLine := int64(200)
	logBytes, err := clientSet.GetK8sCoreV1().Pods(namespace).GetLogs(podName, &apiCoreV1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLine,
	}).Do().Raw()
	if err != nil {
		log.Println(err.Error())
	}
	logData = string(logBytes)
	return
}

// 获取 nodes
func (clientSet K8sClientSet) GetNodes() (nodeList *apiCoreV1.NodeList, err error) {
	nodeList, err = clientSet.GetK8sCoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取 Deployment
func (clientSet K8sClientSet) GetDeployments(namespace string) (podList *apiAppsV1.DeploymentList, err error) {
	podList, err = clientSet.GetK8sAppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取 ReplicationController
func (clientSet K8sClientSet) GetReplicationControllers(namespace string) (podList *apiCoreV1.ReplicationControllerList, err error) {
	podList, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取 ReplicaSets
func (clientSet K8sClientSet) GetReplicaSets(namespace string) (podList *apiAppsV1.ReplicaSetList, err error) {
	podList, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取 Services
func (clientSet K8sClientSet) GetServices(namespace string) (podList *apiCoreV1.ServiceList, err error) {
	podList, err = clientSet.GetK8sCoreV1().Services(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
	}
	return
}

// 获取配置字典
func (clientSet K8sClientSet) GetConfigDict(namespace string) (configMapList *apiCoreV1.ConfigMapList, err error) {
	configMapList, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 删除配置字典
func (clientSet K8sClientSet) DeleteConfigDict(namespace string, resName string) (err error) {
	err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Delete(resName, &metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 获取保密字典
func (clientSet K8sClientSet) GetSecretDict(namespace string) (secretList *apiCoreV1.SecretList, err error) {
	secretList, err = clientSet.GetK8sCoreV1().Secrets(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// 删除保密字典
func (clientSet K8sClientSet) DeleteSecretDict(namespace string, resName string) (err error) {
	err = clientSet.GetK8sCoreV1().Secrets(namespace).Delete(resName, &metav1.DeleteOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func (clientSet K8sClientSet) GetComponentStatus() (statusList *apiCoreV1.ComponentStatusList, err error) {
	statusList, err = clientSet.GetK8sCoreV1().ComponentStatuses().List(metav1.ListOptions{})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func RemoteCommandContainerExec(clusterId string, namespace string, podName string,
	containerName string, handler *k8s_structs.StreamHandler) {
	filePath := models.GetKubeConfigPathByClusterId(clusterId)
	kubeConfig, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
	}
	kubeRestConf, err := clientcmd.RESTConfigFromKubeConfig(kubeConfig)
	if err != nil {
		fmt.Println(err.Error())
	}
	var clientSet K8sClientSet
	clientSet.Client, err = kubernetes.NewForConfig(kubeRestConf)
	if err != nil {
		fmt.Println(err.Error())
	}

	sshReq := clientSet.GetK8sCoreV1().RESTClient().
		Post().
		Resource("pods").Namespace(namespace).Name(podName).
		SubResource("exec").
		VersionedParams(&apiCoreV1.PodExecOptions{
			Container: containerName,
			Command:   []string{"/bin/sh"},
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	// 创建到容器的连接

	executor, err := remotecommand.NewSPDYExecutor(kubeRestConf, "POST", sshReq.URL())
	if err != nil {
		log.Println(err)
	}
	if err := executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		log.Println(err)
	}
}

//应用 yaml文件创建
func (clientSet K8sClientSet) ApplyYaml(namespace string, yamlContent string) (resource interface{}, err error) {
	jsonData, err := yaml.YAMLToJSON([]byte(yamlContent))
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println(err)
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		fmt.Println(err)
		return
	}
	resourceType := mapData["kind"]
	fmt.Println(resourceType)
	switch resourceType {
	case "ReplicaSet":
		var resCreator apiAppsV1.ReplicaSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "ReplicationController":
		var resCreator apiCoreV1.ReplicationController
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "Deployment":
		var resCreator apiAppsV1.Deployment
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().Deployments(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "Service":
		var resCreator apiCoreV1.Service
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Services(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "StatefulSet":
		var resCreator apiAppsV1.StatefulSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().StatefulSets(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "DaemonSet":
		var resCreator apiAppsV1.DaemonSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().DaemonSets(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "ConfigMap":
		var resCreator apiCoreV1.ConfigMap
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	case "Secret":
		var resCreator apiCoreV1.Secret
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			fmt.Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Secrets(namespace).Create(&resCreator)
		if err != nil {
			log.Println(err.Error())
			return
		}
		break
	default:
		err = errors.New(fmt.Sprintf("目前暂不支持该类型的资源：%s", resourceType))
		break
	}
	return
}

func (clientSet K8sClientSet) GetYamlFile(namespace string, resType string, resName string) (resource interface{}, err error) {
	switch resType {
	case "rc":
		resource, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Get(resName, metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		return
	case "rs":
		resource, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).Get(resName, metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		return
	case "pod":
		resource, err = clientSet.GetK8sCoreV1().Pods(namespace).Get(resName, metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		return
	case "deployment":
		resource, err = clientSet.GetK8sExtensionsV1beta1().Deployments(namespace).Get(resName, metav1.GetOptions{})
		if err != nil {
			log.Println(err.Error())
		}
		return
	}

	return
}

func (clientSet K8sClientSet) ScaleResource(namespace string, resType string, resName string, replicaCount int32) (err error) {
	switch resType {
	case "rc":
		scale, err1 := clientSet.GetK8sCoreV1().ReplicationControllers(namespace).GetScale(resName, metav1.GetOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).UpdateScale(resName, scale)
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "rs":
		scale, err1 := clientSet.GetK8sAppsV1().ReplicaSets(namespace).GetScale(resName, metav1.GetOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sAppsV1().ReplicaSets(namespace).UpdateScale(resName, scale)
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "deployment":
		scale, err1 := clientSet.GetK8sExtensionsV1beta1().Deployments(namespace).GetScale(resName, metav1.GetOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sExtensionsV1beta1().Deployments(namespace).UpdateScale(resName, scale)
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	default:
		err = errors.New("不支持的资源类型")
		break
	}
	return
}

func (clientSet K8sClientSet) DeleteResource(namespace string, resType string, resName string) (err error) {
	switch resType {
	case "rc":
		err1 := clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Delete(resName, &metav1.DeleteOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "rs":
		err1 := clientSet.GetK8sAppsV1().ReplicaSets(namespace).Delete(resName, &metav1.DeleteOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "pod":
		err1 := clientSet.GetK8sCoreV1().Pods(namespace).Delete(resName, &metav1.DeleteOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "deployment":
		err1 := clientSet.GetK8sExtensionsV1beta1().Deployments(namespace).Delete(resName, &metav1.DeleteOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	case "service":
		err1 := clientSet.GetK8sCoreV1().Services(namespace).Delete(resName, &metav1.DeleteOptions{})
		if err1 != nil {
			log.Println(err1.Error())
			return
		}
		break
	default:
		err = errors.New("不支持的资源类型")
		break
	}
	return
}

var controller cache.Controller
var store cache.Store

func handlePodAdd(obj interface{}) {
	node := obj.(*apiCoreV1.Pod)
	fmt.Println("informer handle pod add ")
	fmt.Println(fmt.Printf("Node [%s] is added; checking resources...", node.Name))
}

func handlePodUpdate(old, current interface{}) {
	nodeInterface, exists, err := store.GetByKey("test")
	if exists && err == nil {
		fmt.Printf("Found the node [%v] in cache", nodeInterface)
	}
	node := current.(*apiCoreV1.Pod)
	fmt.Println("informer handle update pod")
	data, _ := json.Marshal(node)
	fmt.Println(string(data))
}

func (clientSet K8sClientSet) PodsWatcher() {
	watchList := cache.NewListWatchFromClient(
		clientSet.GetK8sCoreV1().RESTClient(),
		"pods",
		"test",
		fields.Everything())

	store, controller = cache.NewInformer(
		watchList,
		&apiCoreV1.Pod{},
		time.Second*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    handlePodAdd,
			UpdateFunc: handlePodUpdate,
		},
	)

	stop := make(chan struct{})
	go controller.Run(stop)
}

// k8s 监控相关

// 获取 Prometheus 监控 Metricss
func (clientSet K8sClientSet) GetNodesMetrics(query string, start int, end int, step string) (data string, err error) {
	data, err = prometheus.PrometheusQuery(query, start, end, step)
	return
}
