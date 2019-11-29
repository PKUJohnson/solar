package common

import (
	"encoding/json"
	"time"

	std "github.com/PKUJohnson/solar/std"
	crypto "github.com/PKUJohnson/solar/std/toolkit/crypto"
)

const (
	TokenKey = "token_1234567890"
)

func NewTokenFromString(content string) (*Token, error) {
	data, err := crypto.Base64Decode(content)
	if err != nil {
		std.LogInfoLn("Token ParseToken Base64Decode Error:", err)
		return nil, std.ErrIllegalToken
	}
	data, err = crypto.Decrypt(data, []byte(TokenKey))
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

func NewToken(appType string, uid int64, outTime int64, deviceType string) *Token {
	token := new(Token)
	token.AppType = appType
	token.UID = uid
	token.Expiration = time.Now().Unix() + outTime
	token.DeviceType = deviceType
	token.ClientType = ""
	return token
}

func NewTokenWithClient(appType string, uid int64, outTime int64, deviceType string, clientType string) *Token {
	token := new(Token)
	token.AppType = appType
	token.UID = uid
	token.Expiration = time.Now().Unix() + outTime
	token.DeviceType = deviceType
	token.ClientType = clientType
	return token
}

type Token struct {
	AppType    string
	UID        int64
	Expiration int64
	DeviceType string
	ClientType string
}

func (o *Token) GenTokenString() string {
	data, err := json.Marshal(o)
	if err != nil {
		std.LogInfoLn("GenTokenString Error:", err)
		return ""
	}
	data, err = crypto.Encrypt(data, []byte(TokenKey))
	if err != nil {
		std.LogInfoLn("GenTokenString Encrypt Error:", err)
		return ""
	}
	return crypto.Base64Encode(data)
}
