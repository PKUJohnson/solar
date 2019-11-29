package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/labstack/echo"
	"github.com/wednesdaysunny/gpinyin"
	pb "github.com/PKUJohnson/solar/std/pb"
	std "github.com/PKUJohnson/solar/std"
	"github.com/PKUJohnson/solar/std/toolkit/crypto"
	//"github.com/PKUJohnson/solar/std/toolkit/jsonpb"
	"github.com/golang/protobuf/jsonpb"
)

var pbMarshaler jsonpb.Marshaler

func init() {
	pbMarshaler = jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
		EnumsAsInts:  true,
	}
}

var translations = map[string]string{
	"none available": "服务正在维护中，请稍后再试",
	"not found":      "服务正在维护中，请稍后再试",
	"EOF":            "EOF异常，请检查服务log",
}

/*
 * 正确返回
 */
func SuccessResponse(ctx echo.Context, data interface{}, fields ...string) error {
	return successResponse(ctx, data, true, fields...)
}

func SuccessResponseNoETag(ctx echo.Context, data interface{}, fields ...string) error {
	return successResponse(ctx, data, false, fields...)
}

func successResponse(ctx echo.Context, data interface{}, writeETag bool, fields ...string) error {
	if data == nil || fmt.Sprintf("%v", reflect.ValueOf(data)) == "<nil>" { // TODO FBI warning reflect.ValueOf(data).IsNil()
		data = new(pb.CommonResultResp)
	}
	pbMsg, ok := data.(proto.Message)
	var buffer bytes.Buffer
	if ok {
		if err := pbMarshaler.Marshal(&buffer, pbMsg); err != nil {
			std.LogWarnc("response", err, "fail to marshal response")
			return ErrorResponse(ctx, std.ErrIllegalJson)
		}
	} else {
		m := json.NewEncoder(&buffer)
		if err := m.Encode(data); err != nil {
			std.LogWarnc("response", err, "fail to encode response")
			return ErrorResponse(ctx, std.ErrIllegalJson)
		}
	}

	ctx.Set(ResponseCodeKey, std.ErrOK.Code)

	if writeETag && ctx.Request().Method == "GET" {
		// set ETag
		etag := crypto.Md5Str(buffer.Bytes())
		ctx.Response().Header().Set("ETag", etag)
		// compare ETag
		oldTag := ctx.Request().Header.Get("If-None-Match")
		if oldTag == etag {
			return ctx.JSON(http.StatusNotModified, "")
		}
	}
	//return ctx.JSON(http.StatusOK, Payload(buffer.Bytes(), fields...))
	return ctx.JSON(http.StatusOK, Payload(responseDataTransfer(ctx, buffer.Bytes()), fields...))
}

/*
 * 错误返回
 */
func ErrorResponse(ctx echo.Context, err error) error {
	singleErr := std.ErrFromGoErr(err)
	code := int(singleErr.Code)
	std.LogErrorc("gateway", singleErr, "Error response")
	if singleErr.Code == std.ErrInternalFromString.Code { // this should be a micro error
		// Hystrix errors as well
		microErr := &struct {
			ID     string `json:"id"`
			Code   int    `json:"code"`
			Detail string `json:"detail"`
			Status string `json:"status"`
		}{}
		if err := json.Unmarshal([]byte(singleErr.Message), microErr); err != nil {
			std.LogErrorc("json", err, "fail to unmarshal micro-error: "+singleErr.Message+", api: "+ctx.Request().URL.Path)
		} else {
			code = microErr.Code + std.ErrMicro.Code
			if translated, ok := translations[microErr.Detail]; ok {
				singleErr.Message = translated
			} else {
				singleErr.Message = microErr.Detail
				std.LogWarnc("gateway", nil, "fail to translate error message: "+microErr.Detail)
			}
		}
	}
	ctx.Set(ResponseCodeKey, code)
	return ctx.JSON(http.StatusOK, Payload([]byte("{}"), strconv.Itoa(code), singleErr.Message))
}

/*
 * Payload 构建API 返回体
 */
func Payload(data json.RawMessage, fields ...string) *Response {
	res := &Response{
		Code:    std.ErrOK.Code,
		Message: std.ErrOK.Message,
		Data:    data,
	}

	length := len(fields)
	if length >= 1 {
		singleCode, _ := strconv.Atoi(fields[0])
		res.Code = int(singleCode)
	}
	if length >= 2 {
		res.Message = fields[1]
	}

	return res
}

func responseDataTransfer(ctx echo.Context, data []byte) []byte {
	lan := GetLanguage(ctx)
	if lan == "zh-tw" {
		s := gpinyin.ConvertToTraditionalChinese(string(data))
		return []byte(s)
	}
	return data
}
