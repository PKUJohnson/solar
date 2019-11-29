package common

import (
	std "github.com/PKUJohnson/solar/std"
)

type Config struct {
	Micro std.ConfigService `yaml:"micro"`
	Mysql std.ConfigMysql   `yaml:"mysql"`
	Log   std.ConfigLog     `yaml:"log"`
}
