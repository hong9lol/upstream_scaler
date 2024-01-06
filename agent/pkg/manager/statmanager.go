package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/config"
	database "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db"
	entity "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db/entity"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type StatCollector struct {
	kubeCfg    clientcmd.ClientConfig
	restCfg    *rest.Config
	coreClient *kubernetes.Clientset
}
type Notify struct {
	DeploymentName string `json:"deployment_name"`
}

func NewStatCollector() *StatCollector {
	return &StatCollector{}
}

func (s *StatCollector) initK8sConfig() {
	s.kubeCfg = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	s.restCfg, _ = s.kubeCfg.ClientConfig()
	s.coreClient, _ = kubernetes.NewForConfig(s.restCfg)
}

func (s *StatCollector) executeRemoteCommand(podName string, namespace string, containerName string, command string) (string, int64, string, error) {
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := s.coreClient.CoreV1().RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		Param("container", containerName).
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"/bin/sh", "-c", command},
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(s.restCfg, "POST", request.URL())
	if err != nil {
		fmt.Printf("Error creating Executor: %v\n", err)
		return "", 0, "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, namespace, podName)
	}

	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})

	fmt.Println(err)
	fmt.Println(errBuf.String())
	if err != nil {
		return "", 0, "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, namespace, podName)
	}

	return buf.String(), time.Now().UnixMilli(), errBuf.String(), nil
}

func (s *StatCollector) getAllPods() map[string]entity.Pod {
	// NODE_NAME 환경 변수 읽기
	nodeName := os.Getenv("HOST_NAME")

	if nodeName == "" {
		// fmt.Println("HOST_NAME environment variable is not set.")
		// return nil
	}

	// for test
	nodeName = "minikube-m02"

	// fmt.Printf("Get All Pods in Node Name: %s\n", nodeName)

	// 현재 노드에서 실행 중인 Pod 가져오기
	pods, err := s.coreClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		return nil
	}

	podList := map[string]entity.Pod{}
	// fmt.Printf("Pods on Node %s:\n", nodeName)
	for _, podItem := range pods.Items {
		// fmt.Printf("Name: %s\n", podItem.Name)
		// fmt.Printf("Namespace: %s\n", podItem.Namespace)
		// fmt.Printf("Status: %s\n", podItem.Status.Phase)
		if podItem.Status.Phase != "Running" {
			continue
		}
		// fmt.Println("---")
		containerList := map[string]entity.Container{}
		for _, containerItem := range podItem.Spec.Containers {
			// fmt.Printf("Container Name: %s\n", containerItem.Name)
			container := entity.Container{
				Name:       containerItem.Name,
				CPURequest: containerItem.Resources.Requests.Cpu().MilliValue(),
			}
			containerList[container.Name] = container
		}
		pod := entity.Pod{
			Name:       podItem.Name,
			AppLabel:   podItem.Labels["app"],
			Containers: containerList,
		}
		podList[pod.Name] = pod
	}

	return podList
}

func (s *StatCollector) checkResourceUsage(deployment entity.Deployment) bool {
	// resource 사용량 체크
	return false
}

func (s *StatCollector) updateStats(deployment *entity.Deployment) (entity.Deployment, error) {
	db := database.GetInstance()
	prevDeployment, _ := db.GetStat(deployment.Name)
	if prevDeployment == nil {
		db.UpdateStat(deployment)
		return *deployment, nil
	}

	// 삭제된 파드 체크
	for _, prevPod := range prevDeployment.Pods {
		_, ok := deployment.Pods[prevPod.Name]
		if !ok {
			delete(prevDeployment.Pods, prevPod.Name)
		}
	}

	for _, lastPod := range deployment.Pods {
		pod, ok := prevDeployment.Pods[lastPod.Name]
		if ok { // 기존 파드
			for _, container := range deployment.Pods[lastPod.Name].Containers {
				usages := prevDeployment.Pods[lastPod.Name].Containers[container.Name].Usages
				if len(usages) > 9 {
					usages = usages[:len(usages)-1]
				}
				usages = append(usages, deployment.Pods[lastPod.Name].Containers[container.Name].Usages[0])
				prevDeployment.Pods[lastPod.Name].Containers[container.Name] = entity.Container{
					Name:       container.Name,
					CPURequest: container.CPURequest,
					Usages:     usages,
				}
				fmt.Println(prevDeployment)
			}
		} else { // 새로운 파드
			prevDeployment.Pods[lastPod.Name] = pod
		}
	}
	db.UpdateStat(prevDeployment)
	return *prevDeployment, nil
}

func (s *StatCollector) getStat(deploymentName string, podList map[string]entity.Pod) entity.Deployment {
	// var wg sync.WaitGroup
	podBelonged := map[string]entity.Pod{}
	for _, pod := range podList {
		if pod.AppLabel == deploymentName {
			containers := map[string]entity.Container{}
			for _, container := range pod.Containers {
				// wg.Add(1)
				func() {
					usageStr, timestamp, _, _ := s.executeRemoteCommand(pod.Name, config.NAMESPACE, container.Name, "head -1 /sys/fs/cgroup/cpu.stat")
					usageInt64, err := strconv.ParseInt(strings.Split(strings.Split(usageStr, " ")[1], "\r")[0], 10, 64)
					// fmt.Println(container.Name, usageInt64, timestamp)
					if err != nil {
						fmt.Println("conversion error:", err)
						return
					}
					usage := entity.Usage{
						Usage:     usageInt64,
						Timestamp: timestamp,
					}
					container.Usages = append(container.Usages, usage)
					containers[container.Name] = container
					// wg.Done()
				}()
			}
			// wg.Wait()
			pod.Containers = containers
			podBelonged[pod.Name] = pod
		}
	}
	// fmt.Println(podBelonged)
	return entity.Deployment{
		Name: deploymentName,
		Pods: podBelonged,
	}
}

func (s *StatCollector) Start(controllerServiceName string) {
	s.initK8sConfig()
	db := database.GetInstance()
	for {
		// TODO: Check old deployment in db and remove it
		hpas := db.GetAllHPA()
		podList := s.getAllPods()
		for _, hpa := range hpas {
			// 매 초마다 stat 업데이트
			stat := s.getStat(hpa.Target, podList)
			updatedStat, err := s.updateStats(&stat)
			if err != nil {
				fmt.Println(err)
			}
			if s.checkResourceUsage(updatedStat) {
				deploymentName := Notify{DeploymentName: hpa.Target}
				pbytes, _ := json.Marshal(deploymentName)
				buff := bytes.NewBuffer(pbytes)
				_, err = http.Post("http://"+controllerServiceName+":5001/api/v1/notify", "application/json", buff)
				if err != nil {
					panic(err)
				}
			}
		}
		time.Sleep(1000 * time.Millisecond) // 1s
	}
}
