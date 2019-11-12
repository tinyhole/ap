package main

import (
	"fmt"
	"github.com/dubbogo/getty"
	"github.com/pkg/errors"
	"github.com/tinyhole/ap/protocol/pack"
	"github.com/tinyhole/ap/protocol/tcprpc"
	"github.com/tinyhole/ap/transport"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	cli getty.Client
	tt  = transport.NewTcpTransport()
)

const (
	CronPeriod = 6e9
)

func main(){
	cli = getty.NewTCPClient(getty.WithServerAddress("127.0.0.1:8080"),getty.WithConnectionNumber(1))
	socket := &tcpTransportSocket{}
	cli.RunEventLoop(func(session getty.Session) error {
		var (
			ok      bool
			tcpConn *net.TCPConn
		)
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
		session.SetCronPeriod(int(CronPeriod)/1e6) //6 second
		session.SetWaitTime(time.Second * 7)
		//session.SetTaskPool(t.taskPool)

		session.SetPkgHandler(&tcpTransport{})
		session.SetEventListener(socket)
		return nil
	})

	wait()
}

type tcpTransport struct{}


func (t *tcpTransport) Read(session getty.Session, data []byte) (interface{}, int, error) {
	apPack := &pack.ApPackage{}
	unCodec := tcprpc.GetCodec()
	defer tcprpc.PutCodec(unCodec)
	err := unCodec.Unmarshal(data, apPack)
	if err != nil {
		if err == tcprpc.ErrTotalLengthNotEnough || err == tcprpc.ErrFlagLengthNotEnough {
			return nil, 0, nil
		}
		return nil, 0, errors.WithStack(err)
	}
	totalLen := int(tcprpc.HeadLenBytesSize + tcprpc.TotalLenBytesSize + unCodec.(*tcprpc.TcpCodec).TotalLen)
	return apPack, totalLen, err
}

func (t *tcpTransport) Write(session getty.Session, pkg interface{}) ([]byte, error) {
	unCodec := tcprpc.GetCodec()
	defer tcprpc.PutCodec(unCodec)
	data,err := unCodec.Marshal(pkg)
	return data, err
}


type tcpTransportSocket struct{
	session getty.Session
}


func (t *tcpTransportSocket) OnOpen(session getty.Session) error {
	t.session = session
	return nil
}

func (t *tcpTransportSocket) OnError(session getty.Session, err error) {
	t.session.Close()
}

//OnClose 回调上层错误处理函数，关闭msgchan
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

	if pbPkg.Header.Request != nil{
		if pbPkg.Header.Request.ServiceName == "ap" &&
			pbPkg.Header.Request.MethodName == "ping"{
			fmt.Println("ping")
			return
		}

		if pbPkg.Header.Request.ServiceName == "ap" &&
			pbPkg.Header.Request.MethodName == "pong"{
			fmt.Println("pong")
			return
		}

	}
}

func (t *tcpTransportSocket) OnCron(session getty.Session) {
	var (
		err error
	)
	req :=&pack.ApPackage{
		Header: &pack.Header{
			Request:              &pack.RequestMeta{
				ServiceName:          "ap",
				MethodName:           "ping",
				CallType:             pack.CallType_Push,
			},
			Seq:                  0,
		},
	}
	err = session.WritePkg(req,time.Duration(5*time.Second))

	if err != nil {
		session.Close()
		return
	}

	active := session.GetActive()
	if CronPeriod < time.Since(active).Nanoseconds()  {
		session.Close()
	}
}



func wait() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}