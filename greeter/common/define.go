package common

import (
	"fmt"
	std "github.com/PKUJohnson/solar/std"
	svc "github.com/PKUJohnson/solar/std/service"
)

var (
	GlobalConf *Config
)

func InitModel() {
	GlobalConf = new(Config)
	std.LogInfoLn("================init config begin")
	std.LoadConf(GlobalConf, "./conf/greeter.yaml")
	std.InitLog(GlobalConf.Log)
	GlobalConf.Micro.SvcName = svc.GreeterSvcName
	std.LogInfoLn("InitModel ", fmt.Sprintf("%+v",GlobalConf))
	std.LogInfoLn("================init config end")
}

type KvStore interface {
	Set(k, v string, timeout int) error
	Get(k string) (string, error)
}
