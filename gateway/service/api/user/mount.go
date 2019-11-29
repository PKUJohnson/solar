package user

import (
	"github.com/PKUJohnson/solar/gateway/service/middleware"
	"github.com/labstack/echo"
)

// MountUserAPI 挂载delegate API模块
func MountTuUserAPI(g *echo.Group) {
	api := g.Group("/user")
	api.POST("/signin/passwd", HttpPasswordSignIn)
	api.POST("/info/me", HttpGetUserInfo, middleware.UserLoginRequired)
}
