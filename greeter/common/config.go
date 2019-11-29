package common

import (
	std "github.com/PKUJohnson/solar/std"
)

type Config struct {
	Micro std.ConfigService `yaml:"micro"`
	Log   std.ConfigLog     `yaml:"log"`
}
