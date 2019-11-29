package api

import (
	std "github.com/PKUJohnson/solar/std"
	stdc "github.com/PKUJohnson/solar/std/common"
	"github.com/PKUJohnson/solar/std/pb/user"
	"github.com/PKUJohnson/solar/std/toolkit/crypto"
	"github.com/PKUJohnson/solar/user/dao"
	"github.com/PKUJohnson/solar/user/models"
	"golang.org/x/net/context"
	"time"
)

func (d *User) UserPasswordSignIn(ctx context.Context, req *user.UserPasswordSignInReq, res *user.UserPasswordSignInRsp) error {

	fuser := dao.GetFUserByMobile(req.Mobile)
	if (fuser == nil) {
		return std.ErrUserNotExist
	}

	check_passwd := crypto.GetMD5Content(req.Password)
	check_passwd2 := crypto.GetMD5Content(fuser.FUnionId + req.Password)

	// 为了兼容老模式的密码，或者内部使用的密码
	if fuser.Password != check_passwd && fuser.Password != check_passwd2 {
		return std.ErrPassword
	}

	now := time.Now().UTC()
	token_exp := now.Add( 7 * 24 * time.Hour).Unix()

	user := dao.GetUserByFUnionId(fuser.FUnionId)

	token := stdc.Token {
		AppType: "Solar",
		UserId: user.UserId,
		OpenId: user.OpenId,
		SessionKey: "",
		Expiration: token_exp,
	}

	token_str := stdc.TokenToString(&token)
	res.Token 	= token_str
	res.Mobile 	= fuser.Mobile
	res.Name   	= user.UserName
	res.Result = true
	return nil
}

func (d *User) CheckUserToken(ctx context.Context, req *user.CheckUserTokenReq, res *user.CheckUserTokenRsp) error {

	token_data, err := stdc.TokenFromString(req.Token)
	if err != nil {
		res.Result = false
		return std.ErrIllegalToken
	}

	user_item := dao.GetUserById(token_data.UserId)

	if user_item == nil {
		res.Result = false
		return std.ErrIllegalToken
	}

	// check if it is valid token (FIXME: if debug, can comment these)
	if len(user_item.SessionKey) == 0 {
		//res.Result = false
		//return std.ErrIllegalToken
	}

	if len(user_item.OpenId) == 0 {
		//res.Result = false
		//return std.ErrIllegalToken
	}

	// compare user_id and session_key extract from the token
	if user_item.SessionKey != token_data.SessionKey {
		//res.Result = false
		//return std.ErrIllegalToken
	}

	if user_item.OpenId != token_data.OpenId {
		//res.Result = false
		//return std.ErrIllegalToken
	}

	if user_item.Expiration != token_data.Expiration {
		//res.Result = false
		//return std.ErrIllegalToken
	}

	res.Result = true
	res.UserId = user_item.UserId
	return nil
}

func (d *User) GetUserInfo(ctx context.Context, req *user.GetUserInfoReq, res *user.GetUserInfoRsp) error {

	std.LogInfoLn("==== receive get get user info ====")
	std.LogInfoLn("user_id: ", req.UserId)
	std.LogInfoLn("session_key: ", req.SessionKey)

	// get user info from db
	var user_item *models.User
	user_item = dao.GetUserById(req.UserId)
	if user_item == nil {
		return std.ErrUserNotExist
	}

	res.UserName = user_item.UserName
	res.UserDisplayName = user_item.UserDisplayName
	res.Gender = user_item.Gender
	res.Signature = user_item.Signature
	res.Mobile = user_item.Mobile
	res.City   = user_item.City
	res.Province = user_item.Province
	res.Country  = user_item.Country
	res.Language = user_item.Language
	res.Email  = user_item.Email
	res.UserId = user_item.UserId
	res.HasUnionId = (user_item.UnionId != "")

	std.LogInfoLn("==== get user from db ====")
	std.LogInfoLn(user_item.UserId, " " + res.UserName + " " +
		res.UserDisplayName + " " + res.Gender + " " + res.City)

	return nil
}


