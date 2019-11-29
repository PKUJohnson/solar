package rpcserver

import (
	"time"
	"github.com/PKUJohnson/solar/std/pb/user"
	"github.com/PKUJohnson/solar/std/service"
)

var cacheSettings = []service.CacheSetting{
	{user.UserHandler.CheckUserToken, time.Second * 300, false},
}

