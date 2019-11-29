package middleware

import (
	"bytes"
	"net"
	"strings"

	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/helper"
)

func RequestRealIp(h echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		ip := ctx.Request().RemoteAddr

		if x := ctx.Request().Header.Get("x-client-ip"); len(x) > 0 {
			ip = x
		} else if xs := ctx.Request().Header.Get("x-forwarded-for"); len(xs) > 0 {
			for _, x := range strings.Split(xs, ",") {
				if !isPrivateIpv4(x) {
					ip = x
					break
				}
			}
		}

		helper.SetClientIp(ctx, ipHost(ip))
		//ctx.Set(stdc.ClientIp, ipHost(ip))
		return h(ctx)
	}
}

var privateRanges = []struct {
	Lower net.IP
	Upper net.IP
}{
	{Lower: net.ParseIP("10.0.0.0"), Upper: net.ParseIP("10.255.255.255")},
	{Lower: net.ParseIP("172.16.0.0"), Upper: net.ParseIP("172.31.255.255")},
	{Lower: net.ParseIP("192.168.0.0"), Upper: net.ParseIP("192.168.255.255")},
	{Lower: net.ParseIP("169.254.0.0"), Upper: net.ParseIP("169.254.255.255")},
	{Lower: net.ParseIP("127.0.0.0"), Upper: net.ParseIP("127.255.255.255")},
}

func ipHost(ip string) string {
	h := strings.TrimSpace(ip)
	if strings.Index(h, ":") != -1 {
		h, _, _ = net.SplitHostPort(h)
	}
	return h
}

func isPrivateIpv4(ip string) bool {
	h := ipHost(ip)
	if ipv4 := net.ParseIP(h); ipv4 != nil && ipv4.To4() != nil {
		for _, r := range privateRanges {
			if bytes.Compare(ipv4, r.Lower) >= 0 && bytes.Compare(ipv4, r.Upper) <= 0 {
				return true
			}
		}
	}
	return false
}
