package k8s

import (
	"context"
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/remotecommand"
)

type Client struct {
	kubeconf *clientcmdapi.Config
	restconf *rest.Config
	client   *kubernetes.Clientset
}

func NewClient(kubeConfigYaml string) (*Client, error) {
	kubeConf, err := clientcmd.Load([]byte(kubeConfigYaml))
	if err != nil {
		return nil, err
	}

	restConf, err := clientcmd.BuildConfigFromKubeconfigGetter("",
		func() (*clientcmdapi.Config, error) {
			return kubeConf, nil
		},
	)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		return nil, err
	}

	return &Client{
		kubeconf: kubeConf,
		restconf: restConf,
		client:   client,
	}, nil
}

func (c *Client) ServerVersion() (string, error) {
	si, err := c.client.ServerVersion()
	if err != nil {
		return "", err
	}
	return si.String(), nil
}

func (c *Client) WatchContainerLog(
	ctx context.Context,
	req *WatchContainerLogRequest) (io.ReadCloser, error) {
	restReq := c.client.CoreV1().Pods(req.Namespace).GetLogs(req.PodName, req.PodLogOptions)
	return restReq.Stream(ctx)
}

func (c *Client) LoginContainer(ctx context.Context, req *LoginContainerRequest) error {
	restReq := c.client.CoreV1().RESTClient().Post().Resource("Pods").Name(req.PodName).Namespace(req.Namespace).SubResource("exec")

	restReq.VersionedParams(&v1.PodExecOptions{
		Container: req.ContainerName,
		Command:   req.Command,
		Stdin:     true,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(c.restconf, "POST", restReq.URL())
	if err != nil {
		return err
	}

	return executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:             req.Executor,
		Stdout:            req.Executor,
		Stderr:            req.Executor,
		Tty:               true,
		TerminalSizeQueue: req.Executor,
	})
}
