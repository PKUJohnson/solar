package msg

import (
	"github.com/labstack/echo"
)

func MountDemoAPI(g *echo.Group) {
	api := g.Group("/demo")
	api.GET("/hello", HttpHello)
}