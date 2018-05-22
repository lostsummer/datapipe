package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var softJsonKeys = [...]string{
	"App",
	"Version",
	"OPType",
	"Flag",
	"Remark",
	"PageUrl",
}

var softUrlKeys = [...]string{
	"app",
	"version",
	"optype",
	"flag",
	"remark",
	"pageurl",
}

func Soft(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range softUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" || params["optype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::Soft " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::Soft app=[" + params["app"] + "] optype=[" + params["optype"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterInfo("Soft")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Soft " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range softJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["UserAgent"] = ctx.Request().UserAgent()
	dataMap["GlobalID"] = getGlobalID(ctx)
	dataMap["FirstVisitTime"] = getFirstVistTime(ctx)
	dataMap["ClientIP"] = getClientIP(ctx)
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Soft " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::Soft push queue data failed!")
			respstr = respFailed
		}
		return nil
	}
}
