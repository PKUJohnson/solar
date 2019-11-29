package common

// app_type
const (
	Internal = "internal"

	Mobile  = "mobile"
	Email   = "email"
	Wechat  = "wechat"
	Weibo   = "weibo"
	Tencent = "tencent"
	Huawei  = "huawei"
	System  = "system"
)

const (
	InternalTokenOutTime = 7 * 24 * 3600
	UserTokenOutTime     = 10 * 365 * 3600 * 24
)

const (
)

var (
)

func GetRolesMap() map[string]string {
	return map[string]string{
	}
}

const (
	UserLikeSearch = "like"
)

const (
	StatusFrozen = "frozen"
)
