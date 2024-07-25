package terminal

type TerminalSize struct {
	// @gotags: json:"width"
	Width uint16 `json:"width"`

	// @gotags: json:"height"
	Height uint16 `json:"height"`
}

func NewTerminalSize() *TerminalSize {
	return &TerminalSize{}
}
