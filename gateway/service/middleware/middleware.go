package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	echomw "github.com/labstack/echo/middleware"
	"github.com/PKUJohnson/solar/gateway/helper"
	std "github.com/PKUJohnson/solar/std"
	stdb "github.com/PKUJohnson/solar/std/common"
)

/*
 * 跨域配置
 */
type MessageLog struct {
	UserId       int64
	OpenId       string
	Path         string
	InputStream  string
	OutPutStream error
	KeyWord      string
	ClientIP     string
	AccessToken  string
}

func RequestCORS() echo.MiddlewareFunc {

	return echomw.CORSWithConfig(echomw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin", "No-Cache", "X-Requested-With", "If-Modified-Since",
			"Pragma", "Last-Modified", "Cache-Control", "Expires", "Content-Type", "X-E4M-With",
			"Authorization", helper.RequestHeaderTokenName, helper.RequestHeaderPlatform, helper.RequestHeaderApp,
			helper.RequestHeaderClientType, helper.RequestHeaderDeviceID, helper.RequestHeaderLanguage},
		AllowCredentials: true,
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	})
}

func ProcessGetBindData(handle echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var message MessageLog
		if ctx.Path() != "/*" {
			message.AccessToken = ctx.Request().Header.Get(helper.RequestHeaderTokenName)
			xToken := getXToken(ctx)
			if xToken != "" || len(xToken) == 0 {
				tokenData, err := stdb.TokenFromString(xToken)
				if err == nil {
					message.OpenId = tokenData.OpenId
					message.UserId = tokenData.UserId
				} else {
					message.UserId = 0
					message.AccessToken = "$"
					message.OpenId = "$"
				}
			}
			message.Path = ctx.Path()
			message.ClientIP = ctx.RealIP()
			message.InputStream = ctx.QueryString()
		}
		if strings.ToUpper(ctx.Request().Method) == echo.GET {
			queryParams := ctx.QueryParams()
			message.InputStream = queryParams.Encode()
			for key, val := range queryParams {
				keySplits := strings.Split(key, "_")
				newKey := make([]string, len(keySplits))
				for _, v := range keySplits {
					newKey = append(newKey, strings.Title(v))
				}
				queryParams[strings.Join(newKey, "")] = val
				message.KeyWord += val[0]
			}
		} else {
			if len(ctx.Request().Form) > 0 {
				message.InputStream = ctx.Request().Form.Encode()
			} else {
				result, err := ioutil.ReadAll(ctx.Request().Body)
				if err == nil {
					ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(result))
					message.InputStream = bytes.NewBuffer(result).String()
				}
			}
		}
		if message.InputStream == "" {
			message.InputStream = "$"
		}

		if message.KeyWord == "" {
			message.KeyWord = "$"
		}
		message.OutPutStream = handle(ctx)

		// 太长的日志不记录
		input_stream := ""
		if len(message.InputStream) > 1024 {
			input_stream = "日志超过长度"
		} else {
			input_stream = message.InputStream
		}

		if message.Path == "" {

		}
		std.LogInfo(std.LogFields{
			"keyword":     message.KeyWord,
			"InputStream": input_stream, //message.InputStream,
			"ClientIP":    message.ClientIP,
			"URL":         ctx.Request().URL,
			"ReturnCode":  ctx.Get(helper.ResponseCodeKey),
			"OpenID":      message.OpenId,
			"UserId":      message.UserId,
			"AccessToken": message.AccessToken,
			"Path":        message.Path}, "APIGateWay")
		return message.OutPutStream
	}
}

func ProcessGetBindDataO(handle echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if strings.ToUpper(ctx.Request().Method) == echo.GET {
			queryParams := ctx.QueryParams()
			for key, val := range queryParams {
				keySplits := strings.Split(key, "_")
				newKey := make([]string, len(keySplits))
				for _, v := range keySplits {
					newKey = append(newKey, strings.Title(v))
				}
				queryParams[strings.Join(newKey, "")] = val
			}
		}
		return handle(ctx)
	}
}

func ServerErrorHandler(err error, c echo.Context) {
	code := std.ErrInternal.Code
	c.Set(helper.ResponseCodeKey, code)

	msg := std.ErrInternal.Message
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code + std.ErrEcho.Code
		msg = fmt.Sprintf("%v", he.Message)
	} else {
		msg = fmt.Sprintf("%s (original message: %s)", msg, err.Error())
	}

	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(http.StatusOK)
		} else {
			c.JSON(http.StatusOK, helper.Payload([]byte("{}"), strconv.Itoa(code), msg))
		}
	}
}
