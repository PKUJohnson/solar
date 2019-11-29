package api

import (
	"fmt"
	"github.com/PKUJohnson/solar/std/pb/msg"
	"golang.org/x/net/context"
)

func (d *Greeter) SayHello(ctx context.Context, req *msg.HelloRequest, res *msg.HelloReply) error {
	fmt.Println(req.Name)
	res.Message = "I have received " + req.Name
	return nil
}

