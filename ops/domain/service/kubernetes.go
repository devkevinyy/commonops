package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chujieyang/commonops/ops/opslog"
	"k8s.io/api/extensions/v1beta1"
	"time"

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
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/remotecommand"
)

type K8sClientSet struct {
	Client *kubernetes.Clientset
}

func NewKubernetesService(apiServer, token string) (client K8sClientSet, err error) {
	config := &rest.Config{
		Host:        apiServer,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	client.Client, err = kubernetes.NewForConfig(config)
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
	namespaceList, err = clientSet.GetK8sCoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 新建 namespaces
func (clientSet K8sClientSet) CreateNamespaces(name string) (namespace *apiCoreV1.Namespace, err error) {
	namespace, err = clientSet.GetK8sCoreV1().Namespaces().Create(context.TODO(), &apiCoreV1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}, metav1.CreateOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 删除 namespaces
func (clientSet K8sClientSet) DeleteNamespaces(name string) (err error) {
	err = clientSet.GetK8sCoreV1().Namespaces().Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 获取 pod
func (clientSet K8sClientSet) GetPods(namespace string) (podList *apiCoreV1.PodList, err error) {
	podList, err = clientSet.GetK8sCoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取 pod 中容器的日志
func (clientSet K8sClientSet) GetPodContainerLogs(namespace string, podName string, containerName string) (logData string, err error) {
	tailLine := int64(200)
	logBytes, err := clientSet.GetK8sCoreV1().Pods(namespace).GetLogs(podName, &apiCoreV1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLine,
	}).Do(context.TODO()).Raw()
	if err != nil {
		opslog.Error().Println(err)
	}
	logData = string(logBytes)
	return
}

// 获取 nodes
func (clientSet K8sClientSet) GetNodes() (nodeList *apiCoreV1.NodeList, err error) {
	nodeList, err = clientSet.GetK8sCoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取 Deployment
func (clientSet K8sClientSet) GetDeployments(namespace string) (podList *apiAppsV1.DeploymentList, err error) {
	podList, err = clientSet.GetK8sAppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 更新 Deployment 镜像
func (clientSet K8sClientSet) UpdateDeploymentWithImage(namespace string, deployName string, imageName string) (err error) {
	deploy, err := clientSet.GetK8sAppsV1().Deployments(namespace).Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	deploy.Spec.Template.Spec.Containers[0].Image = imageName
	_, err = clientSet.GetK8sAppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 重启 Deployment
func (clientSet K8sClientSet) RestartDeployment(namespace string, deployName string) (err error) {
	deploy, err := clientSet.GetK8sAppsV1().Deployments(namespace).Get(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	envList := deploy.Spec.Template.Spec.Containers[0].DeepCopy().Env
	hasRestartEnv := false
	if len(envList) > 0 {
		var newEnvList []apiCoreV1.EnvVar
		for _, envInfo := range envList {
			if envInfo.Name == "restartMark" {
				hasRestartEnv = true
				envInfo.Value = time.Now().String()
			}
			newEnvList = append(newEnvList, envInfo)
		}
		deploy.Spec.Template.Spec.Containers[0].Env = newEnvList
	}
	if len(envList) == 0 || hasRestartEnv == false {
		deploy.Spec.Template.Spec.Containers[0].Env = append(deploy.Spec.Template.Spec.Containers[0].Env, apiCoreV1.EnvVar{
			Name:  "restartMark",
			Value: time.Now().String(),
		})
	}
	_, err = clientSet.GetK8sAppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取 ReplicationController
func (clientSet K8sClientSet) GetReplicationControllers(namespace string) (podList *apiCoreV1.ReplicationControllerList, err error) {
	podList, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取 ReplicaSets
func (clientSet K8sClientSet) GetReplicaSets(namespace string) (podList *apiAppsV1.ReplicaSetList, err error) {
	podList, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取 Services
func (clientSet K8sClientSet) GetServices(namespace string) (podList *apiCoreV1.ServiceList, err error) {
	podList, err = clientSet.GetK8sCoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

// 获取配置字典
func (clientSet K8sClientSet) GetConfigDict(namespace string) (configMapList *apiCoreV1.ConfigMapList, err error) {
	configMapList, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 删除配置字典
func (clientSet K8sClientSet) DeleteConfigDict(namespace string, resName string) (err error) {
	err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 获取保密字典
func (clientSet K8sClientSet) GetSecretDict(namespace string) (secretList *apiCoreV1.SecretList, err error) {
	secretList, err = clientSet.GetK8sCoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 删除保密字典
func (clientSet K8sClientSet) DeleteSecretDict(namespace string, resName string) (err error) {
	err = clientSet.GetK8sCoreV1().Secrets(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

// 获取 ingress
func (clientSet K8sClientSet) GetIngress(namespace string) (ingressList *v1beta1.IngressList, err error) {
	ingressList, err = clientSet.GetK8sExtensionsV1beta1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
	}
	return
}

func (clientSet K8sClientSet) GetComponentStatus() (statusList *apiCoreV1.ComponentStatusList, err error) {
	statusList, err = clientSet.GetK8sCoreV1().ComponentStatuses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	return
}

func RemoteCommandContainerExec(apiServer, token string, namespace string, podName string,
	containerName string, handler *StreamHandler) {
	config := &rest.Config{
		Host:        apiServer,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		opslog.Error().Println(err)
	}

	sshReq := clientSet.CoreV1().RESTClient().
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

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", sshReq.URL())
	if err != nil {
		opslog.Error().Println(err)
	}
	if err := executor.Stream(remotecommand.StreamOptions{
		Stdin:             handler,
		Stdout:            handler,
		Stderr:            handler,
		TerminalSizeQueue: handler,
		Tty:               true,
	}); err != nil {
		opslog.Error().Println(err)
	}
}

//应用 yaml文件创建
func (clientSet K8sClientSet) ApplyYaml(namespace string, yamlContent string) (resource interface{}, err error) {
	jsonData, err := yaml.YAMLToJSON([]byte(yamlContent))
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	resourceType := mapData["kind"]
	switch resourceType {
	case "ReplicaSet":
		var resCreator apiAppsV1.ReplicaSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "ReplicationController":
		var resCreator apiCoreV1.ReplicationController
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "Deployment":
		var resCreator apiAppsV1.Deployment
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().Deployments(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "Service":
		var resCreator apiCoreV1.Service
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Services(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "Ingress":
		var resCreator v1beta1.Ingress
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sExtensionsV1beta1().Ingresses(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "StatefulSet":
		var resCreator apiAppsV1.StatefulSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().StatefulSets(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "DaemonSet":
		var resCreator apiAppsV1.DaemonSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().DaemonSets(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "ConfigMap":
		var resCreator apiCoreV1.ConfigMap
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "Secret":
		var resCreator apiCoreV1.Secret
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Secrets(namespace).Create(context.TODO(), &resCreator, metav1.CreateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	default:
		err = errors.New(fmt.Sprintf("目前暂不支持该类型的资源：%s", resourceType))
		break
	}
	return
}

//应用 yaml文件更新
func (clientSet K8sClientSet) UpdateYaml(namespace string, resType string, yamlContent string) (resource interface{}, err error) {
	jsonData, err := yaml.YAMLToJSON([]byte(yamlContent))
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	var mapData map[string]interface{}
	err = json.Unmarshal(jsonData, &mapData)
	if err != nil {
		opslog.Error().Println(err)
		return
	}
	switch resType {
	case "rs":
		var resCreator apiAppsV1.ReplicaSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "rc":
		var resCreator apiCoreV1.ReplicationController
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "deployment":
		var resCreator apiAppsV1.Deployment
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().Deployments(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "service":
		var resCreator apiCoreV1.Service
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Services(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "ingress":
		var resCreator v1beta1.Ingress
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sExtensionsV1beta1().Ingresses(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "statefulset":
		var resCreator apiAppsV1.StatefulSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().StatefulSets(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "daemonset":
		var resCreator apiAppsV1.DaemonSet
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sAppsV1().DaemonSets(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "configmap":
		var resCreator apiCoreV1.ConfigMap
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	case "secret":
		var resCreator apiCoreV1.Secret
		err = yaml.Unmarshal(jsonData, &resCreator)
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		resource, err = clientSet.GetK8sCoreV1().Secrets(namespace).Update(context.TODO(), &resCreator, metav1.UpdateOptions{})
		if err != nil {
			opslog.Error().Println(err)
			return
		}
		break
	default:
		err = errors.New(fmt.Sprintf("目前暂不支持该类型的资源：%s", resType))
		break
	}
	return
}

func (clientSet K8sClientSet) GetYamlFile(namespace string, resType string, resName string) (resource interface{}, err error) {
	switch resType {
	case "rc":
		resource, err = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "rs":
		resource, err = clientSet.GetK8sAppsV1().ReplicaSets(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "pod":
		resource, err = clientSet.GetK8sCoreV1().Pods(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "deployment":
		resource, err = clientSet.GetK8sAppsV1().Deployments(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "service":
		resource, err = clientSet.GetK8sCoreV1().Services(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "configmap":
		resource, err = clientSet.GetK8sCoreV1().ConfigMaps(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "secret":
		resource, err = clientSet.GetK8sCoreV1().Secrets(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	case "ingress":
		resource, err = clientSet.GetK8sExtensionsV1beta1().Ingresses(namespace).Get(context.TODO(), resName, metav1.GetOptions{})
		if err != nil {
			opslog.Error().Println(err)
		}
		return
	default:
		err = errors.New("暂不支持该资源的YAML文件查看")
	}
	return
}

func (clientSet K8sClientSet) ScaleResource(namespace string, resType string, resName string, replicaCount int32) (err error) {
	switch resType {
	case "rc":
		scale, err1 := clientSet.GetK8sCoreV1().ReplicationControllers(namespace).GetScale(context.TODO(), resName, metav1.GetOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sCoreV1().ReplicationControllers(namespace).UpdateScale(context.TODO(), resName, scale, metav1.UpdateOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
			return
		}
		break
	case "rs":
		scale, err1 := clientSet.GetK8sAppsV1().ReplicaSets(namespace).GetScale(context.TODO(), resName, metav1.GetOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sAppsV1().ReplicaSets(namespace).UpdateScale(context.TODO(), resName, scale, metav1.UpdateOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
			return
		}
		break
	case "deployment":
		scale, err1 := clientSet.GetK8sAppsV1().Deployments(namespace).GetScale(context.TODO(), resName, metav1.GetOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
			return
		}
		scale.Spec.Replicas = replicaCount
		scale, err1 = clientSet.GetK8sAppsV1().Deployments(namespace).UpdateScale(context.TODO(), resName, scale, metav1.UpdateOptions{})
		if err1 != nil {
			err = err1
			opslog.Error().Println(err1)
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
		err1 := clientSet.GetK8sCoreV1().ReplicationControllers(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
			return
		}
		break
	case "rs":
		err1 := clientSet.GetK8sAppsV1().ReplicaSets(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
			return
		}
		break
	case "pod":
		err1 := clientSet.GetK8sCoreV1().Pods(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
			return
		}
		break
	case "deployment":
		err1 := clientSet.GetK8sAppsV1().Deployments(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
			return
		}
		break
	case "service":
		err1 := clientSet.GetK8sCoreV1().Services(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
			return
		}
		break
	case "ingress":
		err1 := clientSet.GetK8sExtensionsV1beta1().Ingresses(namespace).Delete(context.TODO(), resName, metav1.DeleteOptions{})
		if err1 != nil {
			opslog.Error().Println(err1)
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
	opslog.Info().Println("informer handle pod add ")
	opslog.Info().Println(fmt.Printf("Node [%s] is added; checking resources...", node.Name))
}

func handlePodUpdate(old, current interface{}) {
	nodeInterface, exists, err := store.GetByKey("test")
	if exists && err == nil {
		fmt.Printf("Found the node [%v] in cache", nodeInterface)
	}
	node := current.(*apiCoreV1.Pod)
	opslog.Info().Println("informer handle update pod")
	data, _ := json.Marshal(node)
	opslog.Info().Println(string(data))
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
