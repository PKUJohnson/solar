package service

import (
	std "github.com/PKUJohnson/solar/std"
)

var configIns microConfig

const (
	UserSvcName      = "solar.user"
	GreeterSvcName   = "solar.greeter"
)

var (
	configDev = microConfig{
		Transport:        "utp",
		ClientPoolSize:   64,
		RegisterInterval: 60,
		RegisterTTL:      120,
	}
	configTest = microConfig{
		Transport:        "utp",
		ClientPoolSize:   64,
		RegisterInterval: 300,
		RegisterTTL:      600,
	}
	configProd = microConfig{
		Transport:        "utp",
		ClientPoolSize:   64,
		RegisterInterval: 60,
		RegisterTTL:      120,
	}
)

type microConfig struct {
	Transport        string
	ClientPoolSize   int
	RegisterInterval int
	RegisterTTL      int
}

func init() {
	configIns = configDev
	env := std.ConfEnv()
	if env == std.ConfEnvProduction {
		configIns = configProd
	} else if env == std.ConfEnvTest {
		configIns = configTest
	}
}
