package model

import "github.com/wutong-paas/wutong-gateway/cmd/option"

// Stream -
type Stream struct {
	StreamPort int
}

// NewStream creates a new stream.
func NewStream(conf *option.Config) *Stream {
	return &Stream{
		StreamPort: conf.ListenPorts.Stream,
	}
}
