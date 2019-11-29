package std

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes
// 这样定义错误码:
// 首先看可以归类到哪个HTTP Code, 把这个Code作为错误码的前三位
// 然后后面两位递增
// 如果无法/不想归类到HTTP Code, 请用600开头
// 500** 代表系统错误
// 依次类推f
var (
	AllErrors map[int]string = make(map[int]string)

	ErrOK                 = newErr(20000, "OK")
	ErrNotFound           = newErr(40400, "您访问的资源不存在")
	ErrTooManyRequests    = newErr(42900, "您的操作过于频繁，请稍后再试")
	ErrInternal           = newErr(50000, "Internal error")
	ErrInternalFromString = newErr(50001, "[Should never be returned]")
	ErrDatabase           = newErr(50002, "数据库提交数据失败")
	ErrIllegalJson        = newErr(50003, "json数据格式不正确")
	ErrParams             = newErr(50004, "参数不正确")
	ErrRpcCache           = newErr(50005, "RPC cache 错误")
	ErrRpcCacheMarshal    = newErr(50006, "RPC cache 序列化错误")
	ErrRpcCacheUnmarshal  = newErr(50007, "RPC cache 反序列化错误")
	ErrIllegalToken       = newErr(50008, "非法的token")
	ErrNotImplemented     = newErr(50009, "Not Implemented")
	ErrNoData             = newErr(50010, "No Data")
	ErrDumplicate         = newErr(50011, "记录已经存在")
	ErrOtherClientSignIn  = newErr(50012, "其他客户端登录了")
	ErrUserBan            = newErr(50013, "用户已被禁用")
	ErrTokenExpired       = newErr(50014, "登录过期")
	ErrRpcCacheTimeout    = newErr(50015, "RPC cache timeout")
	ErrServerTooBusy      = newErr(50016, "服务器正忙，请稍后再试")
	ErrBadFileToken       = newErr(50017, "Bad file token")
	ErrFileTokenExpired   = newErr(50018, "File token expired")
	ErrSearchNoDisplay    = newErr(50019, "相关搜索结果不予显示")
	ErrPubPriKeyWrong     = newErr(50020, "秘钥不正确")
	ErrIllegalCooperation = newErr(50021, "非白名单合作方")

	ErrUserNotExist       = newErr(60001, "用户不存在")
	ErrPassword           = newErr(60002, "用户密码错误")

	ErrMicro = newErr(70000, "未定义的 Micro 错误")
	ErrEcho  = newErr(71000, "未定义的 Echo 错误")

	// 系统管理专用代码段
	ErrNotSysAdmin                   = newErr(90001, "非系统管理员")
	ErrSysAdminHasBeenInit           = newErr(90002, "系统管理员已初始化")

	errorMapping             map[string]*Err
)

func init() {
	errorMapping = make(map[string]*Err)
	errorMapping[gorm.ErrRecordNotFound.Error()] = ErrNotFound
	errorMapping[gorm.ErrInvalidSQL.Error()] = ErrDatabase
	errorMapping[gorm.ErrInvalidTransaction.Error()] = ErrDatabase
	errorMapping[gorm.ErrUnaddressable.Error()] = ErrDatabase
}

func newErr(code int, msg string) *Err {
	if code != 99999 {
		if _, ok := AllErrors[code]; ok {
			panic("Duplicated error code!!!")
		}
		AllErrors[code] = msg
	}
	return &Err{code, msg}
}

// NewErr registers a new error.
func NewErr(code int, msg string) *Err {
	return newErr(code, msg)
}

// Err represents the error.
type Err struct {
	Code    int
	Message string
}

func (e *Err) Error() string {
	r, _ := json.Marshal(e)
	return string(r)
}

// GetMessage fetches message from error.
func (e Err) GetMessage() string {
	return e.Message
}

// ErrFromString assembles an error from string.
func ErrFromString(str string) *Err {
	var e Err
	if err := json.Unmarshal([]byte(str), &e); err != nil || e.Code < 10000 {
		return &Err{ErrInternalFromString.Code, str}
	}
	return &e
}

// ErrFromGoErr transforms the golang error object to error.
func ErrFromGoErr(err error) *Err {
	if sysErr, ok := errorMapping[err.Error()]; ok {
		return sysErr
	}

	if e, ok := err.(*Err); ok {
		return e
	}
	return ErrFromString(err.Error())
}

// IsSysErr indicates if the error is error.
func IsSysErr(err error, ivkerr *Err) bool {
	e, ok := err.(*Err)
	if !ok {
		return false
	}
	return e.Code == ivkerr.Code
}
