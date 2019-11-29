package middleware

import (
	"time"

	"fmt"

	"github.com/labstack/echo"
	"github.com/anacrolix/tollbooth"
	tollbooth_config "github.com/anacrolix/tollbooth/config"
	"github.com/PKUJohnson/solar/gateway/helper"
	std "github.com/PKUJohnson/solar/std"
)

// RateLimit protects the APIs.
func RateLimit(limiter *std.ConfigRateLimiter) echo.MiddlewareFunc {
	ratelimiter := tollbooth_config.NewLimiter(limiter.Max, limiter.TTL, &limiter.Redis)
	ratelimiter.Message = limiter.Message
	ratelimiter.MessageContentType = limiter.MessageContentType
	ratelimiter.StatusCode = limiter.StatusCode
	ratelimiter.IPLookups = limiter.IPLookups
	ratelimiter.Methods = limiter.Methods
	ratelimiter.Headers = limiter.Headers

	for _, settings := range limiter.APISettings {
		tollbooth.RegisterAPI(settings.Path, settings.Method, settings.Max, time.Duration(settings.Duration))
	}
	fmt.Println("进入了限速中间件")

	return func(h echo.HandlerFunc) echo.HandlerFunc {
		if !limiter.Enabled {
			return func(c echo.Context) (err error) {
				err = h(c)
				return err
			}
		}

		return func(c echo.Context) (err error) {
			httpError := tollbooth.LimitByRequest(ratelimiter, c.Request())
			if httpError != nil {
				std.LogInfoc("ratelimiter", c.Request().URL.Path)
				return helper.ErrorResponse(c, std.ErrTooManyRequests)
			}

			err = h(c)
			return err
		}
	}
}
