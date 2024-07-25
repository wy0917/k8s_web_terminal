package terminal

type Response struct {
	Request *Request `json:"request"`
	Message string   `json:"message"`
	Data    string   `json:"data"`
}

func NewResponse() *Response {
	return &Response{
		Request: NewRequest(),
	}
}
