package server

import (
	"context"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tinyhole/ap/bucket"
	"github.com/tinyhole/ap/protocol/pack"
	"github.com/tinyhole/ap/transport"
)

type apServer struct {
	opts      *Options
	transport transport.Transport
	rpcClient client.Client
	ctx       context.Context
	cancelFn  context.CancelFunc
	listener  transport.Listener
}

func NewAPServer() Server {
	return &apServer{}
}

func (a *apServer) Init(opts ...Option) error {
	a.opts = newOptions()

	for _, o := range opts {
		o(a.opts)
	}
	a.transport = transport.NewTcpTransport()
	a.ctx, a.cancelFn = context.WithCancel(context.Background())
	a.rpcClient = client.NewClient(client.Registry(etcd.NewRegistry()))
	return nil
}

func (a *apServer) Start() error {
	listener, err := a.transport.Listen(a.opts.Addr)

	if err != nil {
		return errors.Wrap(err, "transport.Listen failed")
	}

	listener.Accept(a.SetUpCon, a.DestroyCon)
	a.listener = listener
	return nil
}

func (a *apServer) Stop() error {
	a.listener.Close()
	a.cancelFn()
	return nil
}

func (a *apServer) SetUpCon(socket transport.Socket) {
	go func(ctx context.Context, socket transport.Socket) {
		var (
			err error
		)
		for {
			select {
			case <-ctx.Done():
				socket.Close()
			default:
				msg := socket.Recv()
				err = a.ProcessMsg(socket, msg)
				if err != nil {
					if err == ErrAuthFailed{
						logrus.Warnf("socket not auth will close")
						socket.Close()
						return
					}
				logrus.Warnf("process msg error [%v]", err)
				}
			}
		}
	}(a.ctx, socket)
}

func (a *apServer) DestroyCon(socket transport.Socket) {
	fID := a.GenerateFid(socket.ID())
	bucket.DefaultSocketBucket.Remove(fID)
}

func (a *apServer) ProcessMsg(socket transport.Socket, reqPack *pack.ApPackage) error {
	var (
		err    error
		reqTmp *Message
		rspTmp *Message
	)
	if socket.GetAuthState() == false {
		//1.身份信息
		uid := reqPack.Header.Auth.Uid
		token := reqPack.Header.Auth.Token
		//2.认证
		//fmt.Sprintf("[%d][%s]", uid, token)
		logrus.Debugf("auth uid [%v] token[%v]", uid, token)
		socket.UpdateAuthState(true)
		//3.认证后处理

		//socket.Close()
		//return ErrAuthFailed
	}
	if reqPack.Header.Request != nil {
		logrus.Debugf("proxy send rpc call")
		reqTmp = NewMessage(reqPack.Body)
		req := a.rpcClient.NewRequest(reqPack.Header.Request.ServiceName, reqPack.Header.Request.MethodName, reqTmp)
		rspPack := &pack.ApPackage{
			Header: &pack.Header{
				Response: &pack.ResponseMeta{
					ErrCode: 0,
					ErrText: "",
				},
				Seq: reqPack.Header.Seq,
			},
		}
		rspTmp = NewMessage([]byte{})
		if err = a.rpcClient.Call(a.ctx, req, rspTmp); err != nil {
			logrus.Errorf("rpc call failed call type [%v] service [%v], method [%v] error [%v]",
				reqPack.Header.Request.CallType,
				reqPack.Header.Request.ServiceName,
				reqPack.Header.Request.MethodName,
				err)

			rspPack.Header.Response.ErrCode = Failed
			rspPack.Header.Response.ErrText = err.Error()
		} else {
			rspPack.Header.Response.ErrCode = OK
			rspPack.Body = rspTmp.data
		}
		socket.Send(rspPack)
	}

	return nil
}

func (a apServer) GenerateFid(id uint32) int64 {
	base := uint64(a.opts.SrvID) << 32
	return int64(base + uint64(id))
}
