package k8s_test

import (
	"context"
	"io"
	"os"
	"testing"
	"web_terminal/k8s"

	"k8s.io/client-go/tools/remotecommand"
)

var (
	client *k8s.Client
)

func TestServerVersion(t *testing.T) {
	v, err := client.ServerVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
}

func TestWatchContainerLog(t *testing.T) {
	req := k8s.NewWatchContainerLogRequest()
	req.Namespace = "default"
	req.PodName = "hello-pod"
	stream, err := client.WatchContainerLog(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	_, err = io.Copy(os.Stdout, stream)
	if err != nil {
		t.Fatal(err)
	}
}

func init() {
	kubeConf := k8s.MustReadContentFile("c:/Users/John/code/web_terminal/kube_config.yml")
	c, err := k8s.NewClient(kubeConf)
	if err != nil {
		panic(err)
	}
	client = c
}

type MockContainerTerminal struct {
	In io.Reader
}

func (t *MockContainerTerminal) Read(p []byte) (n int, err error) {
	return t.In.Read(p)
}

func (t *MockContainerTerminal) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (t *MockContainerTerminal) Next() *remotecommand.TerminalSize {
	return &remotecommand.TerminalSize{
		Width:  100,
		Height: 100,
	}
}

func TestLoginContainer(t *testing.T) {
	reader, writer := io.Pipe()

	term := &MockContainerTerminal{
		In: reader,
	}

	go func() {
		writer.Write([]byte("ls -al / \n"))
	}()

	req := k8s.NewLoginContainerRequest(term)
	req.Namespace = "default"
	req.PodName = "hello-pod"
	err := client.LoginContainer(context.Background(), req)
	if err != nil {
		panic(err)
	}
}
