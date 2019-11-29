package helper

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-os/trace"
	std "github.com/PKUJohnson/solar/std"
	stdc "github.com/PKUJohnson/solar/std/common"
	"github.com/PKUJohnson/solar/std/toolkit/strutil"
)

func SetClientIp(ctx echo.Context, ip string) {
	ctx.Set(ClientIp, ip)
}

func GetClientIp(ctx echo.Context) string {
	if clientIP := ctx.Get(ClientIp); clientIP != nil {
		return clientIP.(string)
	}
	return ""
}

// GetDeviceType extracts device_type from context, login require 时写入
func GetDeviceType(ctx echo.Context) string {
	if deviceType := ctx.Get(DeviceTypeKey); deviceType != nil {
		return deviceType.(string)
	}
	return ""
}

// GetAppVersion extracts app_version from context, login require 时写入
func GetAppVersion(ctx echo.Context) string {
	if appVersion := ctx.Get(AppVersionKey); appVersion != nil {
		return appVersion.(string)
	}
	return ""
}

// GetDeviceID fetches device ID from request header.
func GetDeviceID(req *http.Request) string {
	return req.Header.Get(RequestHeaderDeviceID)
}

// GetTaotieDeviceID fetches taotie device ID from request header.
func GetTaotieDeviceID(req *http.Request) string {
	return req.Header.Get(RequestHeaderTaotieDeviceID)
}

// GetClientType extracts client_type from context, login require 时写入
func GetClientType(ctx echo.Context) string {
	if clientType := ctx.Get(ClientTypeKey); clientType != nil {
		return clientType.(string)
	}
	return ""
}

// GetAppType extracts app_type from context, login require 时写入
func GetAppType(ctx echo.Context) string {
	if appType := ctx.Get(AppTypeKey); appType != nil {
		return appType.(string)
	}
	return ""
}

// GetUserId extracts userId from context, login require 时写入
func GetUserId(ctx echo.Context) int64 {
	if userID := ctx.Get(UserIDKey); userID != nil {
		userIDInt, _ := userID.(int64)
		return userIDInt
	}
	return 0
}

// GetOpenId extracts open_id from context, login require 时写入
func GetOpenId(ctx echo.Context) string {
	if open_id := ctx.Get(OpenIdKey); open_id != nil {
		return open_id.(string)
	}
	return ""
}

// GetSessionKey extracts session_key from context, login require 时写入
func GetSessionKey(ctx echo.Context) string {
	if session_key := ctx.Get(SessionKey); session_key != nil {
		return session_key.(string)
	}
	return ""
}


// GetUserRoles extracts roles from context
func GetUserRoles(ctx echo.Context) []string {
	roleIn := ctx.Get(InternalRolesKey)
	if roleIn != nil {
		return roleIn.([]string)
	}

	return []string{}
}

func GetUserRolesStr(ctx echo.Context) string {
	roles := GetUserRoles(ctx)
	if len(roles) == 0 {
		return ""
	}
	return strings.Join(roles, stdc.CommaSeparator)
}

// CheckParamEmpty checks if parameter is empty.
func CheckParamEmpty(str string) error {
	if str == "" || len(str) == 0 {
		return std.ErrParams
	}
	return nil
}

// GetIdFromPath 从url中获取ID
func GetIdFromPath(ctx echo.Context) (int64, error) {
	entryID := ctx.Param("id")
	entryIDInt, err := strconv.Atoi(entryID)
	if err != nil {
		return 0, std.ErrParams
	}
	return int64(entryIDInt), nil
}

// GetParamFromPath 从url中获取自定义int参数
func GetParamFromPath(ctx echo.Context, param string) (int64, error) {
	entryID := ctx.Param(param)
	entryIDInt, err := strconv.Atoi(entryID)
	if err != nil {
		return 0, std.ErrParams
	}
	return int64(entryIDInt), nil
}

// GetLimit 限制前端API调用传入的limit数量 防止过大
func GetLimit(ctx echo.Context, defaultLimit int64) int64 {
	s := ctx.QueryParam("limit")
	l := strutil.ToInt64(s)
	if l <= 0 {
		l = defaultLimit
	} else if l > 100 {
		l = 100
	}
	return l
}

func GetPlatformFromHeader(ctx echo.Context) string {
	platform := ctx.Request().Header.Get(RequestHeaderPlatform)
	return platform
}

func GetLanguage(ctx echo.Context) string {
	lan := ctx.Request().Header.Get(RequestHeaderLanguage)
	if lan == "zh" || lan == "en" || lan == "zh-tw" {
		return lan
	}
	return "zh"
}

// MicroCtxFromEcho transforms echo.Context to micro context and add trace ID.
func MicroCtxFromEcho(ec echo.Context) context.Context {
	traceID, _ := zipkinFromCtx(ec)

	ctx := metadata.NewContext(context.Background(), metadata.Metadata{
		stdc.Md_USERID:          strconv.FormatInt(GetUserId(ec), 10),
		stdc.Md_APPTYPE:         GetAppType(ec),
		stdc.Md_OPENID:          GetOpenId(ec),
		stdc.Md_SESSIONKEY:      GetSessionKey(ec),
		stdc.Md_CLIENTIP:        GetClientIp(ec),
		stdc.Md_OWNER:           GetHeaderXOwner(ec),
		stdc.Md_DEVICETYPE:      GetHeaderXDevicetype(ec),
		stdc.Md_ADMIN_USER_ROLE: GetUserRolesStr(ec),
		stdc.Md_PLATFORM:        GetPlatformFromHeader(ec),
		stdc.Md_HYSTRIX_STATUS:  hystrixStatus(ec.Path()),
		stdc.Md_APPHEADER:       ec.Request().Header.Get(RequestHeaderApp),
		trace.TraceHeader:       strconv.FormatInt(traceID, 10),
		trace.SpanHeader:        "0",
		trace.ParentHeader:      "0",
	})

	return ctx
}

// hystrixStatus filter the admin API
func hystrixStatus(path string) string {
	path = strings.ToLower(path)
	if strings.Contains(path, "admin") {
		return stdc.Md_HYSTRIX_CLOSED
	}
	return stdc.Md_HYSTRIX_OPEN
}

func zipkinFromCtx(ec echo.Context) (int64, int64) {
	traceIDHex := ec.Get(TraceIDHexKey)
	if traceIDHex != nil {
		traceID, err := strconv.ParseInt(traceIDHex.(string), 16, 64)
		if err != nil {
			std.LogErrorc("zipkin", err, "fail to parse traceID")
		}

		spanID := ec.Get(SpanIDKey).(int64)
		return traceID, spanID
	}

	traceID := RandomInt()
	traceIDHex = RandomHex()
	spanID := RandomInt()
	ec.Set(TraceIDHexKey, traceIDHex)
	ec.Set(SpanIDKey, spanID)
	ec.Response().Header().Add(ResponseHeaderTraceID, traceIDHex.(string))
	return traceID, spanID
}

func GetHeaderXDevicetype(ctx echo.Context) string {
	v := parseHeaderValue(ctx, RequestHeaderAppOffset_DEVICETYPE)
	return v
}

func GetHeaderXOwner(ctx echo.Context) string {
	v := parseHeaderValue(ctx, RequestHeaderAppOffset_OWNER)
	return v
}

func GetHeaderXMarket(ctx echo.Context) string {
	v := parseHeaderValue(ctx, RequestHeaderAppOffset_MARKET)
	return v
}

func parseHeaderValue(ctx echo.Context, index int) string {
	h := ctx.Request().Header.Get(RequestHeaderApp)
	if len(h) == 0 {
		return ""
	}
	for i, v := range strings.Split(h, "|") {
		if index == i {
			return strings.ToLower(strings.TrimSpace(v))
		}
	}
	return ""
}
