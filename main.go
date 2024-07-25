package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"time"
	"web_terminal/k8s"
	"web_terminal/terminal"

	"github.com/gorilla/websocket"
	"github.com/infraboard/mcube/v2/http/response"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout: 60 * time.Second,
		ReadBufferSize:   8192,
		WriteBufferSize:  8192,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

//go:embed ui
var ui embed.FS

func main() {
	http.HandleFunc("/ws/pod/terminal/log", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			response.Failed(w, err)
			return
		}
		// ws.WriteMessage(websocket.TextMessage, []byte("hello websocket server"))

		term := terminal.NewWebSocketTerminal(ws)

		kubeConf := k8s.MustReadContentFile("kube_config.yml")
		k8sClient, err := k8s.NewClient(kubeConf)
		if err != nil {
			term.Failed(err)
			return
		}

		req := k8s.NewWatchContainerLogRequest()
		if err := term.ReadReq(req); err != nil {
			term.Failed(err)
			return
		}

		podReader, err := k8sClient.WatchContainerLog(r.Context(), req)
		if err != nil {
			term.Failed(err)
			return
		}

		_, err = io.Copy(term, podReader)
		if err != nil {
			term.Failed(err)
			return
		}
	})

	http.HandleFunc("/ws/pod/terminal/login", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			response.Failed(w, err)
			return
		}

		// ws.ReadMessage()
		// ws.WriteMessage(websocket.TextMessage, []byte("hello websocket server"))

		term := terminal.NewWebSocketTerminal(ws)

		kubeConf := k8s.MustReadContentFile("kube_config.yml")
		k8sClient, err := k8s.NewClient(kubeConf)
		if err != nil {
			term.Failed(err)
			return
		}

		req := k8s.NewLoginContainerRequest(term)
		if err = term.ReadReq(req); err != nil {
			term.Failed(err)
			return
		}

		err = k8sClient.LoginContainer(r.Context(), req)
		if err != nil {
			term.Failed(err)
			return
		}

	})

	web, _ := fs.Sub(ui, "ui")
	http.Handle("/", http.FileServer(http.FS(web)))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
