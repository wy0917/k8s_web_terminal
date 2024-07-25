package terminal

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gorilla/websocket"
)

var (
	DefaultWriteBuf = 4 * 1024
)

type WebSocketTerminal struct {
	ws *websocket.Conn
	*TerminalResizer

	timeout  time.Duration
	writeBuf []byte
}

func NewWebSocketTerminal(conn *websocket.Conn) *WebSocketTerminal {
	return &WebSocketTerminal{
		ws:              conn,
		TerminalResizer: NewTerminalResizer(),
		timeout:         3 * time.Second,
		writeBuf:        make([]byte, DefaultWriteBuf),
	}
}

func (t *WebSocketTerminal) Close() error {
	return t.ws.Close()
}

func (t *WebSocketTerminal) Read(p []byte) (n int, err error) {
	mt, m, err := t.ws.ReadMessage()
	if err != nil {
		return 0, err
	}

	switch mt {
	case websocket.TextMessage:
		t.HandleCmd(m)
	case websocket.CloseMessage:
		fmt.Printf("Receive client close: %s\n", m)
	default:
		n = copy(p, m)
	}

	return n, nil
}

func (t *WebSocketTerminal) Write(p []byte) (n int, err error) {
	err = t.ws.WriteMessage(websocket.BinaryMessage, p)
	n = len(p)
	return
}

func (t *WebSocketTerminal) WriteText(msg string) {
	err := t.ws.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		fmt.Printf("Write message error, %s\n", err)
	}
}

func (t *WebSocketTerminal) WriteTextf(fmt_ string, a ...any) {
	t.WriteText(fmt.Sprintf(fmt_, a...))
}

func (t *WebSocketTerminal) WriteTextln(fmt_ string, a ...any) {
	t.WriteTextf(fmt_, a...)
	t.WriteText("\r\n")
}

func (t *WebSocketTerminal) close(code int, msg string) {
	fmt.Printf("close code: %d, msg %s\n", code, msg)
	err := t.ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, msg),
		time.Now().Add(t.timeout),
	)

	if err != nil {
		fmt.Printf("close error %s\n", err)
		t.WriteText("\n" + msg)
	}

}

func (t *WebSocketTerminal) Failed(err error) {
	t.close(websocket.CloseGoingAway, err.Error())
}

func (t *WebSocketTerminal) Success(msg string) {
	t.close(websocket.CloseNormalClosure, msg)
}

func (t *WebSocketTerminal) ResetWriteBuf() {
	t.writeBuf = make([]byte, DefaultWriteBuf)
}

func (t *WebSocketTerminal) HandleCmd(m []byte) {
	resp := NewResponse()
	defer t.Response(resp)

	req, err := ParseRequest(m)
	if err != nil {
		resp.Message = err.Error()
		return
	}
	resp.Request = req

	switch req.Command {
	case "resize":
		payload := NewTerminalSize()
		err := json.Unmarshal(req.Params, payload)
		if err != nil {
			resp.Message = err.Error()
			return
		}
		t.SetSize(*payload)
		fmt.Printf("resize add to queue success %s\n", req)
		return
	}

	fn := GetCmdHandleFunc(req.Command)
	if fn == nil {
		resp.Message = "command not found"
		return
	}

	fn(req, resp)
}

func (t *WebSocketTerminal) Response(resp *Response) {
	if resp.Message != "" {
		fmt.Printf("Response error, %s\n", resp.Message)
	}

	if err := t.ws.WriteJSON(resp); err != nil {
		fmt.Printf("write message error, %s\n", err)
	}
}

func (t *WebSocketTerminal) ReadReq(req any) error {
	mt, data, err := t.ws.ReadMessage()
	if err != nil {
		return err
	}

	if mt != websocket.TextMessage {
		return fmt.Errorf("req must be TextMessage, is %d", mt)
	}
	if !json.Valid(data) {
		return fmt.Errorf("req must be json data, but %s", data)
	}

	return json.Unmarshal(data, req)
}

func (t *WebSocketTerminal) WriteTo(r io.Reader) (err error) {
	_, err = io.CopyBuffer(t, r, t.writeBuf)
	if err != nil {
		return err
	}
	defer t.ResetWriteBuf()

	_, err = t.Write(t.writeBuf)
	return
}
