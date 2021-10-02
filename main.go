package main

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	internalapi "k8s.io/cri-api/pkg/apis"
	"k8s.io/kubernetes/pkg/kubelet/remote"
	utilruntime "k8s.io/kubernetes/cmd/kubeadm/app/util/runtime"
	utilsexec "k8s.io/utils/exec"
)

// remoteRuntimeEndpoint: /var/run/dockershim.sock, 与 docker.sock 同目录.
// 每个 docker 容器在启动时都会创建一个新的 containerd-shim 进程,
// 并指定 dockershim.sock 路径
func getRuntimeAndImageServices(
	remoteRuntimeEndpoint string, 
	remoteImageEndpoint string, 
	runtimeRequestTimeout metav1.Duration,
) (internalapi.RuntimeService, internalapi.ImageManagerService, error) {
	rs, err := remote.NewRemoteRuntimeService(
		remoteRuntimeEndpoint, 
		runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	is, err := remote.NewRemoteImageService(
		remoteImageEndpoint, 
		runtimeRequestTimeout.Duration,
	)
	if err != nil {
		return nil, nil, err
	}
	return rs, is, err
}

func main(){
	dockerEp := "/var/run/dockershim.sock"
	// 这个 runtime 是 kubeadm 为 `kubeadm config image pull` 构建的一个简单的 runtime.
	// ta只有 list 和 pull 两个命令.
	// 但是这两个其实是通过 exec 执行的命令(不过 dockershim.sock 仍然需要)
	containerRuntime, err := utilruntime.NewContainerRuntime(utilsexec.New(), dockerEp)
	if err != nil {
		log.Printf("failed to init docker api: %s", err)
		return 
	}
	err = containerRuntime.PullImage("redis")
	if err != nil {
		log.Printf("failed to pull image: %s", err)
		return 
	}
}
