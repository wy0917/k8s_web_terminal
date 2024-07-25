package terminal

import (
	"encoding/json"
	"fmt"
)

func init() {
	RegistryCmdHandler("ping", PingHandleFunc)
}

func PingHandleFunc(r *Request, w *Response) {
	w.Data = "pong"
}

type HandleFunc func(*Request, *Response)

var handleFuncs = map[string]HandleFunc{}

func RegistryCmdHandler(cmd string, fn HandleFunc) {
	handleFuncs[cmd] = fn
}

func GetCmdHandleFunc(cmd string) HandleFunc {
	return handleFuncs[cmd]
}

func ParseRequest(payload []byte) (*Request, error) {
	if !json.Valid(payload) {
		return nil, fmt.Errorf("%v must be json", payload)
	}
	req := NewRequest()
	err := json.Unmarshal(payload, req)
	if err != nil {
		return nil, err
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("validate cmd request error, %s", err)
	}

	return req, nil
}
