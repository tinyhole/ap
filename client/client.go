package client

import (
	"github.com/pkg/errors"
	"github.com/tinyhole/ap/bucket"
	"github.com/tinyhole/ap/protocol/pack"
)

type Client interface {
	Push(fid int64, service, method string, data []byte)
}

type client struct{}

var (
	APClient = &client{}
)

func (c *client) Push(fid int64, service, method string, data []byte) error {
	socket, err := bucket.DefaultSocketBucket.Get(fid)
	if err != nil {
		return errors.Wrapf(err, "DefaultSocketBucket.Get(%d)", fid)
	}

	req := &pack.ApPackage{
		Header: &pack.Header{
			Request: &pack.RequestMeta{
				ServiceName: service,
				MethodName:  method,
				CallType:    pack.CallType_Push,
			},
			Seq: 0,
		},
		Body: data,
	}
	socket.Send(req)
	return nil
}
