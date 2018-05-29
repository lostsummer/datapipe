package handlers

import (
	"strconv"

	"github.com/devfeel/dotweb"
)

var payLogUrlKeys = [...]string{
	"appid",
	"logtype",
	"jsondata",
	"ver",
}

func PayLog(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range payLogUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["appid"] == "" || params["logtype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::PayLog " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PayLog appid=[" + params["appid"] + "] logtype=[" + params["logtype"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("PayLog")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PayLog " + err.Error())
		return nil
	}
	data := params["jsondata"]
	qlen, err := pushQueueData(importerConf, string(data))
	if qlen > 0 && err == nil {
		respstr = strconv.FormatInt(qlen, 10)
	} else {
		innerLogger.Error("HttpServer::PayLog push queue data failed!")
		respstr = respFailed
	}
	return nil
}
