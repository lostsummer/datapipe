package handlers

import (
	"strconv"

	"github.com/devfeel/dotweb"
)

var appDataUrlKeys = [...]string{
	"appid",
	"logtype",
	"jsondata",
	"ver",
}

func AppData(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range appDataUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["appid"] == "" || params["logtype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::AppData " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::AppData appid=[" + params["appid"] + "] logtype=[" + params["logtype"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterInfo("AppData")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::AppData " + err.Error())
		return nil
	}
	data := params["jsondata"]
	qlen, err := pushQueueData(importerConf, string(data))
	if qlen > 0 && err == nil {
		respstr = strconv.FormatInt(qlen, 10)
	} else {
		innerLogger.Error("HttpServer::AppData push queue data failed!")
		respstr = respFailed
	}
	return nil
}
