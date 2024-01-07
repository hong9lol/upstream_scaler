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
	HPAName        string `json:"hpa_name"`
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
	// TODO: need to change to get this info from kubelet api
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

func calUsage(last entity.Usage, prev entity.Usage) float64 {

	// fmt.Println(container.Usages)
	usage := last.Usage - prev.Usage
	timpestamp := last.Timestamp - prev.Timestamp
	// fmt.Println(container.Name, (usage / timpestamp), "\n", container.Usages)
	return float64(usage) / float64(timpestamp)
}

func (s *StatCollector) calUsagePerSec(container entity.Container) float64 {
	// TODO: need to consider the algorithm
	// use 6 data to get avg usage
	if len(container.Usages) < 6 {
		// return when there is only one value in the list
		return 0
	}
	usagePerSec := 0.0
	for i := 0; i < 5; i++ {
		last := container.Usages[len(container.Usages)-(i+1)]
		prev := container.Usages[len(container.Usages)-(i+2)]
		usagePerSec += calUsage(last, prev)
	}
	// fmt.Println(container.Name, (usage / timpestamp), "\n", container.Usages)
	// usageRatePerSec := usagePerSec / container.CPURequest
	return usagePerSec / 5.0
}

func (s *StatCollector) checkResourceUsage(deployment entity.Deployment, hpa entity.HPA) bool {
	totalCpuUsageRate := 0.0
	podCpuUsageRate := 0.0

	// total cpu usage rate
	for _, pod := range deployment.Pods {
		podCpuUsage := 0.0
		podCpuUsageRate = 0
		for _, container := range pod.Containers {
			containerUsagePerSec := s.calUsagePerSec(container)
			podCpuUsage += containerUsagePerSec
			podCpuUsageRate += (containerUsagePerSec / float64(container.CPURequest)) * 100
		}
		fmt.Println(pod.Name, podCpuUsage)

		podCpuUsageRate = (float64(podCpuUsageRate) / float64(len(pod.Containers)))
		totalCpuUsageRate += podCpuUsageRate
	}
	totalCpuUsageRate = (totalCpuUsageRate / float64(len(deployment.Pods)))
	fmt.Println(totalCpuUsageRate)

	// support only cpu resouce currently
	// check if it is higher than hpa.Metrics
	for _, metric := range hpa.Metrics {
		if metric.Name == "cpu" {
			if totalCpuUsageRate > float64(metric.TargetUtilization) {
				return true
			}
		}
	}

	return false
}

func (s *StatCollector) updateStat(deployment entity.Deployment) (entity.Deployment, error) {
	db := database.GetInstance()
	prevDeployment, err := db.GetStat(deployment.Name)

	if err != nil { // not found
		db.UpdateStat(deployment)
		return deployment, nil
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
					usages = usages[1:]
				}
				usages = append(usages, deployment.Pods[lastPod.Name].Containers[container.Name].Usages[0])
				prevDeployment.Pods[lastPod.Name].Containers[container.Name] = entity.Container{
					Name:       container.Name,
					CPURequest: container.CPURequest,
					Usages:     usages,
				}
				// fmt.Println(prevDeployment)
			}
		} else { // 새로운 파드 추가
			prevDeployment.Pods[lastPod.Name] = pod
		}
	}
	db.UpdateStat(prevDeployment)
	return prevDeployment, nil
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
					// TODO: get this info from runc project
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
			updatedStat, err := s.updateStat(stat)
			if err != nil {
				fmt.Println(err)
			}
			if s.checkResourceUsage(updatedStat, hpa) {
				notify := Notify{DeploymentName: hpa.Target, HPAName: hpa.Name}
				pbytes, _ := json.Marshal(notify)
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

// func getKubeletAPI() {
// 	// kubelet API 엔드포인트 URL 설정 (일반적으로 10250 포트 사용)
// 	kubeletURL := "http://localhost:10250/pods"

// 	// HTTP GET 요청 보내기
// 	resp, err := http.Get(kubeletURL)
// 	if err != nil {
// 		fmt.Printf("HTTP GET 요청 실패: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	// 응답 본문 읽기
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Printf("응답 본문 읽기 실패: %v\n", err)
// 		return
// 	}

// 	// 응답 출력
// 	fmt.Printf("kubelet API 응답:\n%s\n", string(body))
// }
