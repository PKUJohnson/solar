package common

import (
	"context"

	"regexp"
	"strings"

	"github.com/micro/go-micro/metadata"
	std "github.com/PKUJohnson/solar/std"
)

const (
	Success        = "success"
	Failed         = "failed"
	CommaSeparator = ","
	DefaultLimit   = 10
)

const (
	LanguageEn = "en"
	LanguageZh = "zh"
)

const (
)

const (
	AppIOSDeviceType        = "ios"
	AppAndroidDeviceType    = "android"
	AppWebDeviceType        = "web"
	AppWebIosDeviceType     = "web_ios"
	AppWebAndroidDeviceType = "web_android"

	ClientPc     = "pc"
	ClientMobile = "mobile"
	ClientPad    = "pad"
	ClientMweb   = "mweb"
	ClientMinapp = "minapp"

	InternalClientJuicy     = "juicy"
	InternalClientXuangubao = "XGB"
	InternalClientSwagger   = "swagger"
	InternalClientTubiaojia = "TBJ"
	InternalClientCong      = "cong"
	InternalClientAthena    = "athena"
)

var (
	CanUseBalanceDeviceType = []string{AppIOSDeviceType, AppAndroidDeviceType, AppWebIosDeviceType, AppWebAndroidDeviceType}
)

const (
	EmailRegex = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	//MobileRegex = `^\+1([38][0-9]|14[57]|5[^4])\d{8}$`
	MobileRegex = `^\+[1-9]{1}[0-9]{3,14}$`
)

const (
	Md_USERID          = "x-userid"
	Md_CLIENTIP        = "x-clientip"
	Md_DEVICETYPE      = "x-devicetype"
	Md_OWNER           = "x-owner"
	Md_ADMIN_USER_ROLE = "x-admin-user-role"
	Md_HYSTRIX_STATUS  = "x-hystrix-status"
	Md_PLATFORM        = "x-platform"
	Md_HYSTRIX_CLOSED  = "closed"
	Md_HYSTRIX_OPEN    = "open"
	Md_APPHEADER       = "x-app-header"
	Md_Token           = "x-token"
	Md_DEVICEID        = "x-deviceid"

	RequestHeaderAppOffset_DEVICETYPE = 1
)

var (
	ContentOrderMap = map[string]string{
		"display_time_desc": "article_display_time desc",
		"display_time_asc":  "article_display_time asc",
		"updated_at_asc":    "updated_at asc",
		"updated_at_desc":   "updated_at desc",
	}
)

func NsqTopicName(id string) string {
	env := std.ConfEnv()
	return id + ".env_" + env
}

func NsqDefaultChannelName(id string) string {
	return NsqChannelName(id, "")
}

func NsqChannelName(id string, nameParam string) string {
	return nsqChannelNameImpl(id, nameParam, false)
}

func NsqEphChannelName(id string, nameParam string) string {
	return nsqChannelNameImpl(id, nameParam, true)
}

func nsqChannelNameImpl(id string, nameParam string, ephemeral bool) string {
	if nameParam == "" {
		nameParam = "default_channel"
	}
	env := std.ConfEnv()
	builder := []string{
		id,
		nameParam,
		"env_" + env,
	}
	ret := strings.Join(builder, ".")
	if ephemeral {
		ret += "#ephemeral"
	}
	return ret
}

func PbMetaGet(k string, ctx context.Context) string {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return ""
	}
	var v string

	for sk, sv := range md {
		if regexp.MustCompile("(?i)^" + k + "$").Match([]byte(sk)) {
			v = sv
			break
		}
	}
	return v
}
