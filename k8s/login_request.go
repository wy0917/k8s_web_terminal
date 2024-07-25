package k8s

import (
	"encoding/json"
	"io"

	"github.com/tidwall/pretty"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	shellCmd = []string{
		"sh",
		"-c",
		`TERM=xterm-256color; export TERM; [ -x /bin/bash ] && ([ -x /usr/bin/script ] && /usr/bin/script -q -c "/bin/bash" /dev/null || exec /bin/bash) || exec /bin/sh`,
	}
)

type LoginContainerRequest struct {
	Namespace     string            `json:"namespace" validate:"required"`
	PodName       string            `json:"pod_name" validate:"required"`
	ContainerName string            `json:"container_name"`
	Command       []string          `json:"command"`
	Executor      ContainerTerminal `json:"-"`
}

func NewLoginContainerRequest(ce ContainerTerminal) *LoginContainerRequest {
	return &LoginContainerRequest{
		Command:  shellCmd,
		Executor: ce,
	}
}

func (req *LoginContainerRequest) String() string {
	jsonReq, _ := json.Marshal(req)
	return string(pretty.Pretty(jsonReq))
}

type ContainerTerminal interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}
