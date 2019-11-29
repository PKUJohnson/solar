package service

import (
	"github.com/PKUJohnson/solar/gateway/common"
	"net/http"

	"github.com/PKUJohnson/solar/gateway/service/api/data"
	"github.com/PKUJohnson/solar/gateway/service/api/msg"
	"github.com/PKUJohnson/solar/gateway/service/api/proxy"
	"github.com/PKUJohnson/solar/gateway/service/api/user"
	waymid "github.com/PKUJohnson/solar/gateway/service/middleware"
	std "github.com/PKUJohnson/solar/std"
	"github.com/labstack/echo"
	echomw "github.com/labstack/echo/middleware"
)

var app = echo.New()

func init() {
	// Recover
	app.Use(echomw.Recover())

	// Cors
	app.Use(waymid.RequestCORS())
	app.Use(waymid.ProcessGetBindData)
	app.Use(waymid.RequestRealIp)
	app.HTTPErrorHandler = waymid.ServerErrorHandler

	MountAPIModule(app)
	MountProxyModule(app)
	MountRootAccess(app)
}

/*
 * 挂在gateway api
 */
func MountAPIModule(e *echo.Echo) {
	apiV1 := e.Group("/apiv1")
	user.MountTuUserAPI(apiV1)
	msg.MountDemoAPI(apiV1)
}

/*
* 挂在代理请求
 */
func MountProxyModule(e *echo.Echo) {
	proxy.MountProxy(e)
}

// 挂载根路径
func MountRootAccess(e *echo.Echo) {
	e.GET("/", func(ctx echo.Context) error {
		return ctx.String(http.StatusOK, "ok")
	})
}

// RunServer starts
func RunServer(config * common.Config) {
	MountAPIModule(app)

	data.Init(&config.EventCollector)
	defer data.Close()

	// add API stats middleware after data collector is initialized
	//app.Use(waymid.APIStats(data.Collector())) // TODO:

	// awesome HTTPs, are you scared?
	if config.Environment.CertPem != "" && config.Environment.KeyPem != "" {
		std.LogInfoLn("⇛ https server started on ", config.Environment.Bind)
		err := app.StartTLS(
			config.Environment.Bind,
			config.Environment.CertPem,
			config.Environment.KeyPem,
		)
		if err!= nil{
			std.LogFatalLn(err, "Gateway server stops unexpectedly.")
		}
	} else {
		std.LogInfoLn("⇛ http server started on ", config.Environment.Bind)
		err := app.Start(config.Environment.Bind)
		if err!= nil{
			std.LogFatalLn(err, "Gateway server stops unexpectedly.")
		}
	}
}
