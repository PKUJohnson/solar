package msg

import (
	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/helper"
	"github.com/PKUJohnson/solar/gateway/rpcclient"
	"github.com/PKUJohnson/solar/std/pb/msg"
)

type HelloArgs struct {
	Name string `json:"name" form:"name" description:"姓名" binding:"required" `
}

func HttpHello(ctx echo.Context) error {
	args := &HelloArgs{}
	if err :=ctx.Bind(args);err != nil {
		return err
	}

	microCtx := helper.MicroCtxFromEcho(ctx)
	req := &msg.HelloRequest{
		Name: args.Name,
	}

	res,err := rpcclient.MsgClient.SayHello(microCtx, req)
	if err != nil{
		return helper.ErrorResponse(ctx, err)
	} else {
		return helper.SuccessResponse(ctx, res.Message)
	}
}
