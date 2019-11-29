package helper

import (
	"encoding/json"
	"strconv"

	std "github.com/PKUJohnson/solar/std"
	stdc "github.com/PKUJohnson/solar/std/common"
	"strings"
)

// Response API正常返回体
type Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

const (
	RequestHeaderClientIP       = "X-Forwarded-For"
	RequestHeaderTokenName      = "X-Solar-Token"
	RequestHeaderApp            = "X-Solar-App"
	RequestHeaderPlatform       = "X-Solar-Platform"
	RequestHeaderClientType     = "X-Client-Type"
	RequestHeaderDeviceID       = "X-Device-ID"
	RequestHeaderTaotieDeviceID = "X-Taotie-Device-ID"
	RequestHeaderLanguage       = "X-Language"
	RequestHeaderAppOffset_OWNER      = 0
	RequestHeaderAppOffset_DEVICETYPE = 1
	RequestHeaderAppOffset_MARKET     = 4

	ResponseHeaderTraceID = "X-Solar-Trace-ID"

	InternalRolesKey = "admin_user_role"
	UserIDKey        = "user_id"
	AppTypeKey       = "app_type"
	OpenIdKey        = "open_id"
	SessionKey       = "session_key"
	ExpirationKey    = "expiration"
	DeviceTypeKey    = "device_type"
	AppVersionKey    = "app_version"
	ClientTypeKey    = "client_type"
	ClientIp         = "clientip"
	TraceIDHexKey    = "trace_id_hex"
	SpanIDKey        = "span_id"
	ResponseCodeKey  = "response_code"
	InstIdKey        = "inst_id"
)

var (
	RoleListRoles = map[string]interface{}{
	}
)

func CheckAdminUserHasRole(userRoles []string, needRoles map[string]interface{}) bool {
	res := false
	for _, ur := range userRoles {
		if _, ok := needRoles[ur]; ok {
			res = true
			break
		}
	}
	return res
}

func BoolErrorAppType(appType string) bool {
	return !std.AppFlagFromName(appType).Valid()
}

func BoolErrorClientType(clientType string) bool {
	if clientType != stdc.ClientMobile && clientType != stdc.ClientPad &&
		clientType != stdc.ClientPc && clientType != stdc.ClientMweb && clientType != stdc.ClientMinapp {
		return true
	}
	return false
}

func BoolErrorDeviceType(deviceType string) bool {
	if deviceType != stdc.AppAndroidDeviceType && deviceType != stdc.AppIOSDeviceType &&
		deviceType != stdc.AppWebDeviceType {
		return true
	}
	return false
}

func BoolAvailablePassword(pwd string) bool {
	if len(pwd) < 6 {
		return false
	}
	if _, err := strconv.ParseInt(pwd, 10, 64); err == nil {
		return false
	}
	return true
}

func BoolErrorInternalClientType(clientType string) bool {
	if clientType != stdc.InternalClientJuicy && clientType != stdc.InternalClientXuangubao &&
		clientType != stdc.InternalClientSwagger && clientType != stdc.InternalClientTubiaojia && clientType != stdc.InternalClientCong && clientType != stdc.InternalClientAthena {
		return true
	}
	return false
}

func TransferClientType(clientType string) string {
	if clientType == "" {
		return ""
	}
	if strings.ToLower(clientType) == stdc.AppIOSDeviceType || strings.ToLower(clientType) == stdc.AppAndroidDeviceType {
		return stdc.Mobile
	}
	return clientType
}
