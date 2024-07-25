package terminal

import (
	"encoding/json"

	"github.com/tidwall/pretty"
)

type Request struct {
	Id      string          `json:"id"`
	Command string          `json:"command"`
	Params  json.RawMessage `json:"params"`
}

func NewRequest() *Request {
	return &Request{}
}

func (r *Request) Validate() error {
	return nil
}

func (r *Request) String() string {
	jsonReq, _ := json.Marshal(r)
	return string(pretty.Pretty(jsonReq))
}
