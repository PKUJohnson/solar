package main

import (
	stdsvc "github.com/PKUJohnson/solar/std/service"
	"github.com/PKUJohnson/solar/user/common"
	"github.com/PKUJohnson/solar/user/models"
	"github.com/PKUJohnson/solar/user/rpcserver"
	"time"
	"github.com/PKUJohnson/solar/gateway/rpcclient"
)

func main() {
	_, _ = time.LoadLocation("Asia/Shanghai")

	stdsvc.StartPprof()

	common.InitModel()
	// table
	models.InitModel(common.GlobalConf.Mysql)
	defer models.CloseDB()
	rpcclient.StartClient(common.GlobalConf.Micro)
	// rpc server
	rpcserver.StartServer()
}
