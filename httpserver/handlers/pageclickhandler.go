package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var pageClickJsonKey = [...]string{
	"App",
	"Module",
	"ClickKey",
	"PageUrl",
	"HtmlType",
	"Remark",
	"Ver",
}

var pageClickUrlKeys = [...]string{
	"app",
	"module",
	"clickkey",
	"clickremark",
	"pageurl",
	"htmltype",
	"remark",
	"ver",
}

func PageClick(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range pageClickUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" || params["clickkey"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::PageClick " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PageClick app=[" + params["app"] + "] clickkey=[" + params["clickkey"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterInfo("PageClick")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PageClick " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range pageClickJsonKey {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["UserAgent"] = ctx.Request().UserAgent()
	dataMap["GlobalID"] = getGlobalID(ctx)
	dataMap["FirstVisitTime"] = getFirstVistTime(ctx)
	dataMap["ClientIP"] = getClientIP(ctx)
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PageClick " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::PageClick push queue data failed!")
			respstr = respFailed
		}
		return nil
	}
}
