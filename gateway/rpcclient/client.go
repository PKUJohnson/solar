package rpcclient

import (
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/service"
	"github.com/PKUJohnson/solar/std/pb/msg"
	"github.com/PKUJohnson/solar/std/pb/user"
	"github.com/PKUJohnson/solar/std/service/hystrix"
)

var (
	UserClient user.UserService
	MsgClient msg.GreeterService
)

func StartClient(cnf std.ConfigService) {
	svc := service.NewService(cnf)
	svc.Init()
	hystrix.InitHystrix()

	UserClient = user.UserServiceClient(service.UserSvcName, svc.Client())
	MsgClient = msg.NewGreeterService(service.GreeterSvcName, svc.Client())
}
