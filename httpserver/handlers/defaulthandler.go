package handlers

import (
	"fmt"

	"github.com/devfeel/dotweb"
)

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("welcome to datapipe | version=1")
	return nil
}

func Test(ctx dotweb.Context) error {
	retstr := fmt.Sprintf("user_agent: %s\nglobal_id: %s\nfirst_visit_time: %s\nclient_ip: %s\nwrite_time: %s",
		getUserAgent(ctx), getGlobalID(ctx), getFirstVistTime(ctx), getClientIP(ctx), getNowFormatTime())

	ctx.WriteString(retstr)
	return nil
}
