package transport

import (
	"github.com/dubbogo/getty"
	"github.com/tinyhole/ap/protocol/pack"
	"time"
)

type tcpTransportSocket struct {
	session   getty.Session
	msgChan   chan *pack.ApPackage
	authState bool
	destroyFn func(socket Socket)
}

func NewTcpTransportSocket(session getty.Session, destroyFn func(socket Socket)) getty.EventListener {
	return &tcpTransportSocket{
		session:   session,
		msgChan:   make(chan *pack.ApPackage, 1024),
		authState: false,
		destroyFn: destroyFn,
	}
}

func (t *tcpTransportSocket) GetAuthState() bool {
	return t.authState
}

func (t *tcpTransportSocket) UpdateAuthState(state bool) {
	t.authState = state
}

func (t *tcpTransportSocket) Recv() *pack.ApPackage {
	pkg := <-t.msgChan
	return pkg
}

func (t *tcpTransportSocket) Send(intrepidPackage *pack.ApPackage) {
	t.session.WritePkg(intrepidPackage, time.Second*5)
}

func (t *tcpTransportSocket) Close() error {
	t.destroyFn(t)
	t.session.Close()
	close(t.msgChan)
	return nil
}

func (t *tcpTransportSocket) Local() string {

	return t.session.LocalAddr()
}

func (t *tcpTransportSocket) Remote() string {
	return t.session.RemoteAddr()
}

func (t *tcpTransportSocket) ID() uint32 {
	return t.session.ID()
}

func (t *tcpTransportSocket) OnOpen(session getty.Session) error {
	//t.session = session
	return nil
}

func (t *tcpTransportSocket) OnError(session getty.Session, err error) {
	t.session.Close()
}

func (t *tcpTransportSocket) OnClose(session getty.Session) {

}

func (t *tcpTransportSocket) OnMessage(session getty.Session, pkg interface{}) {
	var (
		pbPkg *pack.ApPackage
		ok    bool
	)
	pbPkg, ok = pkg.(*pack.ApPackage)
	if !ok {
		return
	}
	t.msgChan <- pbPkg
}

func (t *tcpTransportSocket) OnCron(seession getty.Session) {
}
