package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/tinyhole/ap/test/server/hello"
)

func main() {
	service := micro.NewService(micro.Name("go.micro.srv.hello"))

	service.Init()

	hello.RegisterHelloHandler(service.Server(), &Handler{})
	service.Run()
}

type Handler struct{}

func (h *Handler) SayHello(ctx context.Context, req *hello.SayHelloReq, rsp *hello.SayHelloRsp) error {
	rsp.Reply = fmt.Sprintf("hello %s", req.Name)
	return nil
}
