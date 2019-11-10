package transport

import (
	"context"
	"github.com/tinyhole/ap/protocol"
)

type Options struct {
	Addr  string
	Codec protocol.Codec
}

type ListenOptions struct {
	Ctx context.Context
}

type ListenOption func(o *ListenOptions)
