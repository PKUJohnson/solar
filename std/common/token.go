package common

import (
	"encoding/json"
	std "github.com/PKUJohnson/solar/std"
	crypto "github.com/PKUJohnson/solar/std/toolkit/crypto"
)

const (
	TokenKey    = "tushuangshuang2501hahaha"
	SysTokenKey = "tushuangshuang20190815st"
)

func TokenFromString(content string) (*Token, error) {
	return tokenFromStringCore(content, TokenKey)
}

func tokenFromStringCore(content, tokenKey string) (*Token, error) {
	data, err := crypto.Base64Decode(content)
	if err != nil {
		std.LogInfoLn("Token ParseToken Base64Decode Error:", err)
		return nil, std.ErrIllegalToken
	}
	data, err = crypto.Decrypt(data, []byte(tokenKey))
	if err != nil {
		std.LogInfoLn("Token ParseToken Decrypt Error:", err)
		return nil, std.ErrIllegalToken
	}
	token := new(Token)
	if err = json.Unmarshal(data, &token); err != nil {
		std.LogInfoLn("Token ParseToken Unmarshal Error:", err)
		return nil, std.ErrIllegalToken
	}
	return token, nil

}

func SysTokenFromString(content string) (*Token, error) {
	return tokenFromStringCore(content, SysTokenKey)
}

//func NewToken(appType string, uid int64, outTime int64, deviceType string) *Token {
//	token := new(Token)
//	token.AppType = appType
//	token.UID = uid
//	token.Expiration = time.Now().Unix() + outTime
//	token.DeviceType = deviceType
//	token.ClientType = ""
//	return token
//}
//
//func NewTokenWithClient(appType string, uid int64, outTime int64, deviceType string, clientType string) *Token {
//	token := new(Token)
//	token.AppType = appType
//	token.UID = uid
//	token.Expiration = time.Now().Unix() + outTime
//	token.DeviceType = deviceType
//	token.ClientType = clientType
//	return token
//}

//type Token struct {
//	AppType    string
//	UID        int64
//	Expiration int64
//	DeviceType string
//	ClientType string
//}

type Token struct {
	AppType    string
	UserId     int64
	OpenId     string
	SessionKey string
	Expiration int64  // FIXME: not used for now
}

func TokenToString(o *Token) string {
	return tokenToString(o, TokenKey)
}

func tokenToString(o *Token, tokenKey string) string {
	data, err := json.Marshal(o)
	if err != nil {
		std.LogInfoLn("GenTokenString Error:", err)
		return ""
	}
	data, err = crypto.Encrypt(data, []byte(tokenKey))
	if err != nil {
		std.LogInfoLn("GenTokenString Encrypt Error:", err)
		return ""
	}
	return crypto.Base64Encode(data)
}

func SysTokenToString(o *Token) string {
	return tokenToString(o, SysTokenKey)
}
