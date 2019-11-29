package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	std "github.com/PKUJohnson/solar/std"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}

func DecryptWeChatData(encryptedData, iv, appId, sessionKey string) (map[string]interface{}, error) {
	if len(sessionKey) != 24 {
		return nil, errors.New("sessionKey length is error")
	}
	aesKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	if len(iv) != 24 {
		return nil, errors.New("iv length is error")
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = PKCS7UnPadding(aesPlantText)

	var decrypted map[string]interface{}
	aesPlantText = []byte(strings.Replace(string(aesPlantText), "\a", "", -1))
	err = json.Unmarshal([]byte(aesPlantText), &decrypted)
	if err != nil {
		return nil, err
	}

	if decrypted["watermark"].(map[string]interface{})["appid"] != appId {
		return nil, errors.New("appId is not match")
	}

	return decrypted, nil
}


func Struct2Map(obj interface{}) map[string]interface{} {
	var data = make(map[string]interface{})
	j, _ := json.Marshal(obj)
	json.Unmarshal(j, &data)
	return data
}

type QRCode struct {
	Scene string  `json:"scene"`
	Page  string  `json:"page"`
	Width int32   `json:"width"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   time.Duration `json:"expires_in"`
}

func GetAccessToken() *AccessToken {
	wechat_url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" +
		WeiChat_AppId + "&secret=" + WeiChat_AppSecret

	_, body, errs := gorequest.New().Get(wechat_url).EndBytes()

	var result AccessToken

	if errs != nil {
		std.LogInfoLn("get wechat access token error: ", string(body))
		return nil
	}
	if err := json.Unmarshal(body, &result); err != nil {
		std.LogInfoLn(std.LogFields{
			"error": err,
			"body":  string(body),
		}, "Failed to unmarshal access token result")
		return nil
	}
	std.LogInfoLn("Token: ", result.AccessToken)
	std.LogInfoLn("Token expire in ", result.ExpiresIn)
	return &result
}

type CheckRsp struct {
	Errcode int32   `json:"errcode"`
	Errmsg  string  `json:"errmsg"`
}

func CheckTextRisk(text string, token string) bool {
	v := url.Values{}
	v.Add("access_token",token)
	if text == "" {
		return true
	}

	check_url := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?%s",v.Encode())
	data := map[string]string{}
	data["content"] = text
	d, _ := json.Marshal(data)
	str := string(d)
	resp, body, _ := gorequest.New().Post(check_url).
		Type("json").Send(str).End()
	rsp := CheckRsp{}

	if resp.StatusCode == 200 {
		json.Unmarshal([]byte(body),&rsp)
	} else {
		std.LogInfoc("WeChat", fmt.Sprintf("Check text risk error %s",resp.Status))
	}
	if rsp.Errcode == 0 {
		return true
	} else {
		std.LogInfoc("WeChat", fmt.Sprintf("Check text risk failed: %s",rsp.Errmsg))
		return false
	}
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

