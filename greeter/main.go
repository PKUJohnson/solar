package main

import (
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	"github.com/PKUJohnson/solar/greeter/common"
	"github.com/PKUJohnson/solar/greeter/rpcserver"
	stdsvc "github.com/PKUJohnson/solar/std/service"
	"time"
)

func main() {
	_, _ = time.LoadLocation("Asia/Shanghai")

	stdsvc.StartPprof()

	common.InitModel()
	// table
	rpcclient.StartClient(common.GlobalConf.Micro)
	// rpc server
	rpcserver.StartServer()
}
