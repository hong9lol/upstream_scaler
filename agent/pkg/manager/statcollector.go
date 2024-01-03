package manager

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	database "github.com/hong9lol/upstream_scaler/tree/master/agent/pkg/db"
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

func NewStatCollector() *StatCollector {
	return &StatCollector{}
}

func (s *StatCollector) executeRemoteCommand(name string, namespace string, command string) (string, string, error) {
	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := s.coreClient.CoreV1().RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(name).
		SubResource("exec").
		VersionedParams(&v1.PodExecOptions{
			Command: []string{"/bin/sh", "-c", command},
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(s.restCfg, "POST", request.URL())
	err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		return "", "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, namespace, name)
	}

	return buf.String(), errBuf.String(), nil
}

func (s *StatCollector) initK8sConfig() {
	s.kubeCfg = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	s.restCfg, _ = s.kubeCfg.ClientConfig()
	s.coreClient, _ = kubernetes.NewForConfig(s.restCfg)
}

func (s *StatCollector) getAllPods() []string {
	// NODE_NAME 환경 변수 읽기
	nodeName := os.Getenv("HOST_NAME")

	if nodeName == "" {
		fmt.Println("HOST_NAME environment variable is not set.")
		return nil
	}

	fmt.Printf("Node Name: %s\n", nodeName)

	// for test
	nodeName = "minikube-m02"

	// 현재 노드에서 실행 중인 Pod 가져오기
	pods, err := s.coreClient.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", nodeName),
	})
	if err != nil {
		fmt.Printf("Error getting pods: %v\n", err)
		return nil
	}

	// Pod 목록 출력
	var podNameList []string
	fmt.Printf("Pods on Node %s:\n", nodeName)
	for _, pod := range pods.Items {
		fmt.Printf("Name: %s\n", pod.Name)
		fmt.Printf("Namespace: %s\n", pod.Namespace)
		fmt.Printf("Status: %s\n", pod.Status.Phase)
		fmt.Println("---")
		podNameList = append(podNameList, pod.Name)
	}

	return podNameList
}

func (s *StatCollector) Start() {
	s.initK8sConfig()
	db := database.GetInstance()
	for {
		hpas := db.GetAllHPA()
		pods := s.getAllPods()
		for _, hpa := range hpas {
			// hpa에서 target의 이름과 겹치는 pod 찾기
			// 각 파드의 resource 사용량 가져오기
			// 컨디션 확인
			// request
			fmt.Println(hpa)
		}
		// ret, _, _ := ExecuteRemoteCommand("media-mongodb-6575f756d7-t5ngn", "default", "head -1 /sys/fs/cgroup/cpu.stat")
		// fmt.Println(ret)

		time.Sleep(1000 * time.Millisecond) // 1s
	}
}
