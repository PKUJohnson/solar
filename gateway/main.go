package main

import (
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	"github.com/PKUJohnson/solar/gateway/service"
	"github.com/PKUJohnson/solar/gateway/common"
	stdsvc "github.com/PKUJohnson/solar/std/service"
)

func main() {
	common.InitConfig()
	rpcclient.StartClient(common.GlobalConf.Micro)
	stdsvc.StartPprof()
	service.RunServer(common.GlobalConf)
}
