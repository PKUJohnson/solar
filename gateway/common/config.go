package common

import (
	"fmt"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/ec"
	"github.com/PKUJohnson/solar/std/toolkit"
	"gopkg.in/redis.v5"
	"time"
)

var (
	GlobalConf *Config
	RedisClient *redis.Client
)

func InitConfig() {
	_, _ = time.LoadLocation("Asia/Shanghai")
	GlobalConf = new(Config)
	std.LogInfoLn("================init config begin")
	std.LoadConf(GlobalConf, "./conf/gateway.yaml")
	RedisClient = toolkit.InitRedis(GlobalConf.Redis)
	std.InitLog(GlobalConf.Log)
	std.LogInfoLn("InitModel ", fmt.Sprintf("%+v",GlobalConf))
	std.LogInfoLn("================init config end")
}

type Config struct {
	Environment       ConfigEnv           `json:"environment" yaml:"environment"`
	Micro             std.ConfigService   `json:"micro" yaml:"micro"`
	Log               std.ConfigLog       `json:"logger" yaml:"logger"`
	Redis             std.ConfigRedis     `json:"redis" yaml:"redis"`
	EventCollector    ec.Config           `json:"event_collector" yaml:"event_collector"`
}

type ConfigEnv struct {
	Env     string `json:"env" yaml:"env"`
	Bind    string `json:"bind" yaml:"bind"`
	CertPem string `json:"cert_pem" yaml:"cert_pem"`
	KeyPem  string `json:"key_pem" yaml:"key_pem"`
}
