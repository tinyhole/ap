package transport

import (
	"fmt"
	"github.com/dubbogo/getty"
	"github.com/sirupsen/logrus"
	"github.com/tinyhole/ap/protocol/pack"
	"github.com/tinyhole/ap/protocol/tcprpc"
	"net"
	"time"
)

type tcpTransport struct {
	opts Options
}

func NewTcpTransport() Transport {
	return &tcpTransport{opts: Options{
		Codec: tcprpc.NewTcpCodec(),
	}}
}

func (t *tcpTransport) Read(session getty.Session, data []byte) (interface{}, int, error) {
	apPack := &pack.ApPackage{}
	err := t.opts.Codec.Unmarshal(data, apPack)
	return apPack, len(data), err
}

func (t *tcpTransport) Write(session getty.Session, pack interface{}) ([]byte, error) {
	data, err := t.opts.Codec.Marshal(pack)
	return data, err
}

type tcpTransportListener struct {
	l          getty.Server
	tTransport *tcpTransport
}

func (t *tcpTransportListener) Accept(setUp, destroy func(socket Socket)) error {
	go t.l.RunEventLoop(func(session getty.Session) error {

		var (
			ok      bool
			tcpConn *net.TCPConn
		)
		logrus.Debug("======")
		if tcpConn, ok = session.Conn().(*net.TCPConn); !ok {
			panic(fmt.Sprintf("%s, session.conn{%#v} is not tcp connection\n", session.Stat(), session.Conn()))
		}

		tcpConn.SetNoDelay(true)
		tcpConn.SetKeepAlive(true)
		tcpConn.SetKeepAlivePeriod(time.Duration(time.Second * 6))
		tcpConn.SetReadBuffer(262144)
		tcpConn.SetWriteBuffer(65536)

		session.SetName(fmt.Sprintf("tcp-%s", session.RemoteAddr()))
		session.SetMaxMsgLen(65536)
		session.SetRQLen(1024)
		session.SetWQLen(1024)
		session.SetReadTimeout(time.Second)
		session.SetWriteTimeout(time.Second * 5)
		session.SetCronPeriod(int(time.Second * 6))
		session.SetWaitTime(time.Second * 7)
		//session.SetTaskPool(t.taskPool)

		session.SetPkgHandler(t.tTransport)
		eventListener := NewTcpTransportSocket(session,destroy)
		session.SetEventListener(eventListener)
		setUp(eventListener.(*tcpTransportSocket))
		return nil
	})
	logrus.Warn("________--")
	return nil
}

func (t *tcpTransportListener) Addr() string {
	return t.l.Listener().Addr().String()
}

func (t *tcpTransportListener) Close() error {
	return t.l.Listener().Close()
}

func (t *tcpTransport) Listen(addr string, opts ...ListenOption) (Listener, error) {
	var (
		options ListenOptions
	)

	for _, o := range opts {
		o(&options)
	}

	listener := &tcpTransportListener{
		l:          getty.NewTCPServer(getty.WithLocalAddress(addr)),
		tTransport: t,
	}
	return listener, nil
}
