package middleware

import (
	"strings"

	"github.com/avct/uasurfer"
	"github.com/labstack/echo"
	stdc "github.com/PKUJohnson/solar/gateway/common"
	"github.com/PKUJohnson/solar/gateway/helper"
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/pb/user"
)

func ParseUA(myUA string) string {
	ua := uasurfer.Parse(myUA)
	if ua.DeviceType == uasurfer.DeviceComputer {
		return stdc.ClientPc
	} else if ua.DeviceType == uasurfer.DeviceTablet {
		return stdc.ClientPad
	}
	return stdc.ClientMobile
}

//UserInfo gets user info
func UserInfo(handle echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		ctx.Set(helper.UserIDKey, 0)
		ctx.Set(helper.AppTypeKey, "")
		ctx.Set(helper.ExpirationKey, 0)
		ctx.Set(helper.DeviceTypeKey, "")
		ctx.Set(helper.ClientTypeKey, "")
		ctx.Set(helper.AppVersionKey, "")

		headApp := ctx.Request().Header.Get(helper.RequestHeaderApp)
		appInfoList := strings.Split(headApp, "|")
		// wscn|iOS|5.0.2|8.1|0
		if len(appInfoList) >= 2 {
			deviceType := strings.ToLower(appInfoList[1])
			ctx.Set(helper.DeviceTypeKey, deviceType)

			if len(appInfoList) >= 3 {
				appVersion := appInfoList[2]
				ctx.Set(helper.AppVersionKey, appVersion)
			}

			if helper.BoolErrorDeviceType(deviceType) {
				std.LogWarnLn("UserInfo:deviceTypeError", deviceType, headApp, len(deviceType))
			}
		} else {
			ctx.Set(helper.DeviceTypeKey, stdc.AppWebDeviceType)
		}
		ua := ctx.Request().UserAgent()
		clientType := ParseUA(ua)
		headerClientType := ctx.Request().Header.Get(helper.RequestHeaderClientType)
		if headerClientType == "" {
			ctx.Set(helper.ClientTypeKey, clientType)
		} else {
			ctx.Set(helper.ClientTypeKey, helper.TransferClientType(headerClientType))
		}

		xToken := ctx.Request().Header.Get(helper.RequestHeaderTokenName)
		if xToken == "" {
			return handle(ctx)
		}

		rpcReq := new(user.CheckUserTokenReq)
		rpcReq.Token = xToken
		info_res, err := rpcclient.UserClient.CheckUserToken(helper.MicroCtxFromEcho(ctx), rpcReq)
		if err != nil {
			std.LogWarnLn("CheckUserToken:", err)
			return handle(ctx)
		}

		tokenData, err := stdc.NewTokenFromString(xToken)
		if err != nil {
			std.LogInfoc("auth", "invalid user token: "+xToken)
			return handle(ctx)
		}

		ctx.Set(helper.UserIDKey, tokenData.UID)
		ctx.Set(helper.AppTypeKey, tokenData.AppType)
		ctx.Set(helper.InstIdKey, info_res.InstId)
		ctx.Set(helper.ExpirationKey, tokenData.Expiration)
		if tokenData.DeviceType != "" {
			ctx.Set(helper.DeviceTypeKey, tokenData.DeviceType)
		}
		if tokenData.ClientType != "" {
			ctx.Set(helper.ClientTypeKey, tokenData.ClientType)
		}
		return handle(ctx)
	}
}
