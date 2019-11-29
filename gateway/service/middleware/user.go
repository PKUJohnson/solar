package middleware

import (
	"bytes"
	"github.com/PKUJohnson/solar/std/pb/user"
	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/helper"
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	std "github.com/PKUJohnson/solar/std"
	stdb "github.com/PKUJohnson/solar/std/common"
	"io/ioutil"
	"time"
)

func getXToken(ctx echo.Context) string {
	xToken := ctx.Request().Header.Get(helper.RequestHeaderTokenName)
	if xToken == "" {
		xToken = ctx.FormValue(helper.RequestHeaderTokenName)
	}
	return xToken
}


/*
 * 验证是否登录
 */
func UserLoginRequired(handle echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		xToken := getXToken(ctx)
		if xToken == "" && len(xToken) == 0 {
			return helper.ErrorResponse(ctx, std.ErrIllegalToken)
		}
		tokenData, err := stdb.TokenFromString(xToken)
		if err != nil {
			return helper.ErrorResponse(ctx, err)
		}
		if tokenData.Expiration < time.Now().UTC().Unix() {
			return helper.ErrorResponse(ctx, std.ErrTokenExpired) //ErrTokenExpired
		}
		rpcReq := new(user.CheckUserTokenReq)
		rpcReq.Token = xToken
		userinfo_res, err := rpcclient.UserClient.CheckUserToken(helper.MicroCtxFromEcho(ctx), rpcReq)
		if err != nil {
			std.LogWarnLn("UserLoginRequired:", err)
			return helper.ErrorResponse(ctx, err)
		}

		if userinfo_res.Result == false {
			return helper.ErrorResponse(ctx, std.NewErr(99999, "token检查失败"))
		}

		ctx.Set(helper.AppTypeKey, tokenData.AppType)
		ctx.Set(helper.InstIdKey, userinfo_res.InstId)
		ctx.Set(helper.UserIDKey, tokenData.UserId)
		ctx.Set(helper.OpenIdKey, tokenData.OpenId)
		ctx.Set(helper.SessionKey, tokenData.SessionKey)
		ctx.Set(helper.ExpirationKey, tokenData.Expiration)

		return handle(ctx)
	}
}

func GetUserInfoOrDefaultInfo(handle echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		xToken := getXToken(ctx)
		if xToken == "" && len(xToken) == 0 {
			ctx.Set(helper.AppTypeKey, "solar")
			ctx.Set(helper.UserIDKey, 0)
			ctx.Set(helper.OpenIdKey, "")
			ctx.Set(helper.SessionKey, "")
			ctx.Set(helper.ExpirationKey,"")
			return handle(ctx)
			//return helper.ErrorResponse(ctx, std.ErrIllegalToken)
		}
		tokenData, err := stdb.TokenFromString(xToken)
		if err != nil {
			return helper.ErrorResponse(ctx, err)
		}
		//if tokenData.Expiration < time.Now().Unix() {
		//	return helper.ErrorResponse(ctx, std.ErrOtherClientSignIn) //ErrTokenExpired
		//}
		rpcReq := new(user.CheckUserTokenReq)
		rpcReq.Token = xToken
		info_res, err := rpcclient.UserClient.CheckUserToken(helper.MicroCtxFromEcho(ctx), rpcReq)
		if err != nil {
			std.LogWarnLn("UserLoginRequired:", err)
			return helper.ErrorResponse(ctx, err)
		}

		ctx.Set(helper.AppTypeKey, tokenData.AppType)
		ctx.Set(helper.UserIDKey, tokenData.UserId)
		ctx.Set(helper.OpenIdKey, tokenData.OpenId)
		ctx.Set(helper.SessionKey, tokenData.SessionKey)
		ctx.Set(helper.ExpirationKey, tokenData.Expiration)
		ctx.Set(helper.InstIdKey, info_res.InstId)

		return handle(ctx)
	}
}

func CloneBody(ctx echo.Context) []byte{
	//todo: need to consider concurrency issue
	req := ctx.Request()
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body.Close() //  must close
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func Either(fallbackErr error, orHandlers ...ArgsCheckHandlerFunc) echo.MiddlewareFunc {
	return func(handle echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var responseChannel = make(chan bool, 4)
			defer close(responseChannel)
			bodyBytes := CloneBody(ctx)

			for _, h := range orHandlers {
				go func(theHandler ArgsCheckHandlerFunc) {
					err := theHandler(ctx, bodyBytes)
					if err != nil {
						responseChannel <- false
					} else {
						responseChannel <- true
					}
				}(h)
			}

			passed := false
			responsedReq := 0
			totalReq := len(orHandlers)
			for result := range responseChannel {
				//don't use shortcut break as send to a closed channel will panic
				passed = passed || result

				responsedReq += 1
				if responsedReq == totalReq {
					break
				}
			}

			if passed {
				return handle(ctx)
			} else {
				return helper.ErrorResponse(ctx, fallbackErr)
			}
		}
	}
}

type ArgsCheckHandlerFunc func(echo.Context, []byte) error

func All(fallbackErr error, andHandlers ...ArgsCheckHandlerFunc) echo.MiddlewareFunc {
	return func(handle echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var responseChannel = make(chan bool, 4)
			defer close(responseChannel)
			bodyBytes := CloneBody(ctx)

			for _, h := range andHandlers {
				go func(theHandler ArgsCheckHandlerFunc) {
					err := theHandler(ctx, bodyBytes)
					if err != nil {
						responseChannel <- false
					} else {
						responseChannel <- true
					}
				}(h)
			}

			passed := true
			responsedReq := 0
			totalReq := len(andHandlers)
			for result := range responseChannel {
				//don't use shortcut break as send to a closed channel will panic
				passed = passed && result

				responsedReq += 1
				if responsedReq == totalReq {
					break
				}
			}

			if passed {
				return handle(ctx)
			} else {
				return helper.ErrorResponse(ctx, fallbackErr)
			}
		}
	}
}
