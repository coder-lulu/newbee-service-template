package main

import (
	"flag"
	"fmt"

	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/internal/config"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/internal/server"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/internal/svc"
	"github.com/coder-lulu/newbee-{{SERVICE_NAME}}-rpc/types/{{SERVICE_NAME}}"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/{{SERVICE_NAME}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{SERVICE_NAME}}.Register{{SERVICE_NAME_UPPER}}Server(grpcServer, server.New{{SERVICE_NAME_UPPER}}Server(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
