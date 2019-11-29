package impl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/labstack/echo"
	"github.com/PKUJohnson/solar/gateway/helper"
)

func HttpRedirectUrl(ctx echo.Context) error {
	redirect := ctx.QueryParam("redirect")
	redirectParts, err := url.Parse(redirect)
	if err != nil {
		return helper.ErrorResponse(ctx, err)
	}
	rScheme, rHost, rPath := redirectParts.Scheme, redirectParts.Host, redirectParts.Path
	rQuery := redirectParts.Query()
	fmt.Println(rScheme, rPath, rQuery)
	okRedirect := false

	if okRedirect {
		requestUri := ctx.Request().URL
		requestArgs := requestUri.Query()
		delete(requestArgs, "redirect")
		resRedirectUri := rScheme + "://" + rHost + "/" + rPath
		mergeArgs := []string{}
		for key := range rQuery {
			mergeArgs = append(mergeArgs, key+"="+rQuery.Get(key))
		}
		for key := range requestArgs {
			mergeArgs = append(mergeArgs, key+"="+rQuery.Get(key))
		}
		resRedirectUri = resRedirectUri + "?" + strings.Join(mergeArgs, "&")
		return ctx.Redirect(302, resRedirectUri)
	}
	err = ctx.Redirect(302, redirect)
	return err
}
