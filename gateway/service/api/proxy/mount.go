package proxy

import (
	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/service/api/proxy/impl"
)

func MountProxy(e *echo.Echo) {
	e.GET("/auth/redirect", impl.HttpRedirectUrl)
}
