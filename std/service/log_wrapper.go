package service

import (
	"github.com/micro/go-micro/server"
	std "github.com/PKUJohnson/solar/std"
	"golang.org/x/net/context"
)

// LogWrapper wraps the RPC handler for capturing the panic log.
func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		defer func() {
			if e := recover(); e != nil {
				std.LogRecover(e)
			}
		}()

		std.LogDebugLn(std.LogFields{"Method": req.Method()}, "Pre RPC")
		err := fn(ctx, req, rsp)
		std.LogDebugLn(std.LogFields{"Method": req.Method(), "ContentType": req.ContentType(), "Service": req.Service()}, "Post RPC")
		return err
	}
}
