package user

import (
	"github.com/PKUJohnson/solar/gateway/helper"
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	"github.com/PKUJohnson/solar/std/pb/user"
	"github.com/labstack/echo"
)

type SignInArgs struct {
	Mobile string `json:"mobile" form:"mobile" description:"电话" binding:"required" `
	Passwd string `json:"passwd" form:"passwd" description:"密码" binding:"required" `
}

func HttpPasswordSignIn(ctx echo.Context) error {
	args := &SignInArgs{}
	if err :=ctx.Bind(args);err != nil {
		return err
	}

	microCtx := helper.MicroCtxFromEcho(ctx)
	req := &user.UserPasswordSignInReq{
		Mobile: args.Mobile,
		Password: args.Passwd,
	}

	res, err := rpcclient.UserClient.UserPasswordSignIn(microCtx, req)
	if err != nil{
		return helper.ErrorResponse(ctx, err)
	} else {
		return helper.SuccessResponse(ctx, res)
	}
}

func HttpGetUserInfo(ctx echo.Context) error {
	microCtx := helper.MicroCtxFromEcho(ctx)
	req := &user.GetUserInfoReq{
		UserId: helper.GetUserId(ctx),
	}
	res, err := rpcclient.UserClient.GetUserInfo(microCtx, req)
	if err != nil{
		return helper.ErrorResponse(ctx, err)
	} else {
		return helper.SuccessResponse(ctx, res)
	}
}
