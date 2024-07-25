package terminal

import (
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalResizer struct {
	sizeChan chan remotecommand.TerminalSize
	doneChan chan struct{}
}

func NewTerminalResizer() *TerminalResizer {
	size := &TerminalResizer{
		sizeChan: make(chan remotecommand.TerminalSize, 10),
		doneChan: make(chan struct{}),
	}

	return size
}

func (i *TerminalResizer) SetSize(ts TerminalSize) {
	i.sizeChan <- remotecommand.TerminalSize{
		Width:  ts.Width,
		Height: ts.Height,
	}
}

func (i *TerminalResizer) Next() *remotecommand.TerminalSize {
	select {
	case size := <-i.sizeChan:
		return &size
	case <-i.doneChan:
		return nil
	}
}
