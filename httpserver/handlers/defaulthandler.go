package handlers

import (
	"TechPlat/datapipe/global"
	"fmt"

	"github.com/devfeel/dotweb"
)

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	retstr := fmt.Sprintf("welcome to datapipe|version=%s", global.Version)
	ctx.WriteString(retstr)
	return nil
}

func Test(ctx dotweb.Context) error {
	versionInfo := fmt.Sprintf("Version: %s, Branch: %s, Build: %s, Build time: %s\n",
		global.Version, global.Branch, global.CommitID, global.BuildTime)
	clientInfo := fmt.Sprintf("user_agent: %s\nglobal_id: %s\nfirst_visit_time: %s\nclient_ip: %s\nnow_time: %s",
		getUserAgent(ctx), getGlobalID(ctx), getFirstVistTime(ctx), getClientIP(ctx), getNowFormatTime())

	ctx.WriteString(versionInfo + clientInfo)
	return nil
}
