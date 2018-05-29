package handlers

import (
	"strconv"

	"github.com/devfeel/dotweb"
)

var userLogUrlKeys = [...]string{
	"appid",
	"logtype",
	"jsondata",
	"ver",
}

func UserLog(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range userLogUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["appid"] == "" || params["logtype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::UserLog " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::UserLog appid=[" + params["appid"] + "] logtype=[" + params["logtype"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("UserLog")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::UserLog " + err.Error())
		return nil
	}
	data := params["jsondata"]
	qlen, err := pushQueueData(importerConf, string(data))
	if qlen > 0 && err == nil {
		respstr = strconv.FormatInt(qlen, 10)
	} else {
		innerLogger.Error("HttpServer::UserLog push queue data failed!")
		respstr = respFailed
	}
	return nil
}
