package rpcserver

import (
	"github.com/micro/go-micro"
	std "github.com/PKUJohnson/solar/std"
	stdservice "github.com/PKUJohnson/solar/std/service"
	"github.com/PKUJohnson/solar/greeter/common"
	"github.com/PKUJohnson/solar/greeter/rpcserver/api"
	"time"
)

func registerHandler(service micro.Service, process interface{}) {
	handlerProcess := service.Server().NewHandler(process)
	service.Server().Handle(handlerProcess)
}

func StartServer() {
	// Register Handlers
	configObj := common.GlobalConf
	svc := stdservice.NewService(
		configObj.Micro,
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)
	svc.Init()
	stdservice.ConfigRpcCache(cacheSettings)
	// Register Handlers
	registerHandler(svc, new(api.Greeter))

	std.LogInfoc("greeter",  "Start Greeter server.")
	// Run server
	if err := svc.Run(); err != nil {
		std.LogFatalLn(err, "Greeter server stops unexpectedly.")
	}
}
