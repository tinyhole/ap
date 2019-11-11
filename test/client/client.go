package main

import (
	"flag"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/tinyhole/ap/protocol"
	"github.com/tinyhole/ap/protocol/pack"
	"github.com/tinyhole/ap/protocol/tcprpc"
	"github.com/tinyhole/ap/test/client/hello"
	"net"
	"os"
)

var (
	host  = flag.String("host", "localhost", "host")
	port  = flag.String("port", "8080", "port")
	codec = tcprpc.NewTcpCodec().(protocol.Codec)
)

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	tcpc := conn.(*net.TCPConn)
	tcpc.SetKeepAlive(true)
	tcpc.SetNoDelay(true)

	defer conn.Close()
	fmt.Println("Connecting to " + *host + ":" + *port)
	done := make(chan string)
	go handleWrite(conn, done)
	go handleRead(conn, done)
	fmt.Println(<-done)
	fmt.Println(<-done)
}

func handleWrite(conn net.Conn, done chan string) {
	req := &hello.SayHelloReq{
		Name: "Jerry",
	}

	reqByte, _ := proto.Marshal(req)

	reqPack := &pack.ApPackage{
		Header: &pack.Header{
			Request: &pack.RequestMeta{
				ServiceName: "hello",
				MethodName:  "Hello.SayHello",
				CallType:    pack.CallType_Sync,
			},
			Seq: 0,
			Auth: &pack.AuthInfo{
				Uid:   11,
				Token: "11",
			},
		},
		Body: reqByte,
	}

	reqPackByte, _ := codec.Marshal(reqPack)
	for i := 10; i > 0; i-- {
		_, e := conn.Write(reqPackByte)
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	done <- "Sent"
}
func handleRead(conn net.Conn, done chan string) {
	for {
		buf := make([]byte, 10240)
		reqLen, err := conn.Read(buf)
		if err != nil {
			//fmt.Println("Error to read message because of ", err)
			continue
		}
		fmt.Println(string(buf[:reqLen-1]))
	}
	done <- "Read"
}
