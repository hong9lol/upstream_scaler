package manager

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	database "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db"
	entity "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db/entity"
	"github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/util"

	v1 "k8s.io/api/core/v1"
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

	client *http.Client
	req    *http.Request
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

func (s *StatCollector) readStat(metricPath string) (string, int64) {
	// fmt.Println(metricPath + "/cpu.stat")
	ret, err := util.ExecCommand("cat", metricPath+"/cpu.stat")
	if err != nil {
		fmt.Printf("Error to read stat: %v\n", err)
		return "0 0\n", time.Now().UnixMilli()
	}
	return ret, time.Now().UnixMilli()
}

func (s *StatCollector) getAllPods() map[string]entity.Pod {
	// NODE_NAME 환경 변수 읽기
	nodeName := os.Getenv("NODE_NAME")

	if nodeName == "" {
		fmt.Println("NODE_NAME environment variable is not set.")
		return nil
	}

	fmt.Println("Node name: ", nodeName)

	// 현재 노드에서 실행 중인 Pod 가져오기
	// TODO: need to change to get this info from kubelet api
	// pods, err := s.coreClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
	// 	FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	// })

	pods, err := s.getPods(s.client, s.req)
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		return nil
	}

	// TODO
	// add remove old pod
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
		// fmt.Printf("Container Name: %s\n", podItem)
		for _, containerSpec := range podItem.Spec.Containers {
			// fmt.Printf("Container Name: %s\n", containerSpec.Name)
			container := entity.Container{
				Name:       containerSpec.Name,
				CPURequest: containerSpec.Resources.Requests.Cpu().MilliValue(),
			}
			containerList[container.Name] = container
		}
		// fmt.Printf("Container #1: %s\n", containerList)
		for _, containerStatus := range podItem.Status.ContainerStatuses {
			_container := containerList[containerStatus.Name]
			containerID := containerStatus.ContainerID
			// start := time.Now()
			containerID = strings.Split(containerID, "containerd://")[1]
			_arg := "*" + containerID + "*"
			path, err := util.ExecCommand("find", "/host", "-name", _arg)
			if err != nil {
				panic(err)
			}
			// end := time.Since(start)
			// fmt.Println("execution time:", end)
			_container.MetricPath = strings.Split(path, "\n")[0]
			containerList[containerStatus.Name] = _container
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

// TODO: weight moving average
func (s *StatCollector) calUsagePerSec(container entity.Container) float64 {
	// TODO: need to consider the algorithm
	// use 6 data to get avg usage
	const totalDataNum = 6
	if len(container.Usages) < totalDataNum {
		// return when there is less than totalDataNum value in the list
		return 0
	}
	usagePerSec := 0.0
	cnt := 5.0
	wmaWeight := 1.0
	for i := 0; i < 5; i++ {
		if container.Usages[len(container.Usages)-(i+1)].Usage == int64(0) || container.Usages[len(container.Usages)-(i+2)].Usage == int64(0) {
			cnt -= 1.0
			continue
		}
		last := container.Usages[len(container.Usages)-(i+1)]
		prev := container.Usages[len(container.Usages)-(i+2)]
		usagePerSec += calUsage(last, prev) * wmaWeight
	}
	// fmt.Println(container.Name, (usage / timpestamp), "\n", container.Usages)
	// usageRatePerSec := usagePerSec / container.CPURequest
	return usagePerSec / cnt
}

func (s *StatCollector) checkResourceUsage(deployment entity.Deployment, hpa entity.HPA) bool {
	totalCpuUsageRate := 0.0
	podCpuUsageRate := 0.0

	// total cpu usage rate
	fmt.Println("\n[ Resource Usage Info of dployment:", deployment.Name, "]")
	fmt.Println("Pods [", len(deployment.Pods), "]")
	for _, pod := range deployment.Pods {
		podCpuUsage := 0.0
		podCpuUsageRate = 0
		for _, container := range pod.Containers {
			containerUsagePerSec := s.calUsagePerSec(container)
			podCpuUsage += containerUsagePerSec
			podCpuUsageRate += (containerUsagePerSec / float64(container.CPURequest)) * 100
		}
		fmt.Println(" -", pod.Name, ", CPU Usage:", math.Ceil(podCpuUsage))

		podCpuUsageRate = (float64(podCpuUsageRate) / float64(len(pod.Containers)))
		totalCpuUsageRate += podCpuUsageRate
	}
	totalCpuUsageRate = (totalCpuUsageRate / float64(len(deployment.Pods)))
	fmt.Println("Total CPU Usage Rate(%):", totalCpuUsageRate)
	fmt.Fprintln(os.Stdout, []any{"\n"}...)

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
					// usageStr, timestamp, _, _ := s.executeRemoteCommand(pod.Name, config.NAMESPACE, container.Name, "head -1 /sys/fs/cgroup/cpu.stat")
					// usageInt64, err := strconv.ParseInt(strings.Split(strings.Split(usageStr, " ")[1], "\r")[0], 10, 64)
					usageStr, timestamp := s.readStat(container.MetricPath)
					usageInt64, err := strconv.ParseInt(strings.Split(strings.Split(usageStr, " ")[1], "\n")[0], 10, 64)
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
	s.initKubeletClient()
	db := database.GetInstance()
	var hpas []entity.HPA = db.GetAllHPA()
	var podList map[string]entity.Pod = s.getAllPods()
	statusInterval := 0
	for {
		if statusInterval > 14 { // update every 15s
			hpas = db.GetAllHPA()
			podList = s.getAllPods()
			statusInterval = 0
		}
		var notifyList []Notify
		for _, hpa := range hpas {
			// 매 초마다 stat 업데이트
			stat := s.getStat(hpa.Target, podList)
			if len(stat.Pods) < 1 { // no pod in this node
				continue
			}

			updatedStat, err := s.updateStat(stat)
			if err != nil {
				fmt.Println(err)
			}
			if s.checkResourceUsage(updatedStat, hpa) {
				notify := Notify{DeploymentName: hpa.Target, HPAName: hpa.Name}
				notifyList = append(notifyList, notify)
				pbytes, _ := json.Marshal(notifyList)
				buff := bytes.NewBuffer(pbytes)
				_, err = http.Post("http://"+controllerServiceName+".upstream-system.svc.cluster.local:5001/api/v1/notify", "application/json", buff)
				if err != nil {
					panic(err)
				}
			}
		}
		statusInterval++
		time.Sleep(1000 * time.Millisecond) // 1s
	}
}

func (s *StatCollector) getPods(client *http.Client, req *http.Request) (v1.PodList, error) {
	// Load the CA certificate.
	// Set up HTTPS client with the loaded CA certificate and token.
	//nodeNamehttps://$NODE_NAME:10250/pods
	// resp, err := client.Get("https://" + nodeName + ":10250/pods")
	// Prepare the request.
	// client, req := newFunction()

	// Perform the request.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("HTTP GET 요청 실패: %v\n", err)
		// return
	}
	defer resp.Body.Close()

	// 응답 본문 읽기
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("응답 본문 읽기 실패: %v\n", err)
		// return
	}

	// 응답 출력
	// fmt.Printf("kubelet API 응답:\n%s\n", string(body))

	podList := v1.PodList{}
	err = json.Unmarshal([]byte(body), &podList)
	if err != nil {
		log.Fatal(err)
	}

	return podList, err
}

func (s *StatCollector) initKubeletClient() (*http.Client, *http.Request) {
	token, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		panic(err)
	}

	caCert, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            caCertPool,
				InsecureSkipVerify: true,
			},
		},
	}

	nodeName := os.Getenv("NODE_NAME")

	req, err := http.NewRequest("GET", "https://"+nodeName+":10250/pods", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+string(token))
	s.req = req
	s.client = client
	return client, req
}
