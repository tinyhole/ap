package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/tinyhole/ap/test/server/hello"
)

func main() {
	service := micro.NewService(micro.Name("hello"),micro.Registry(etcd.NewRegistry()))

	service.Init()

	hello.RegisterHelloHandler(service.Server(), &Handler{})
	service.Run()
}

type Handler struct{}

func (h *Handler) SayHello(ctx context.Context, req *hello.SayHelloReq, rsp *hello.SayHelloRsp) error {
	rsp.Reply = fmt.Sprintf("hello %s", req.Name)
	fmt.Println(rsp.Reply)
	return nil
}
