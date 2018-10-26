package handlers

import (
	"TechPlat/datapipe/global"
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
		innerLogger.Error("HttpServer::PayLog " + global.LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PayLog appid=[" + params["appid"] + "] logtype=[" + params["logtype"] + "]")
		ctx.WriteString(respstr)
	}()
	data := params["jsondata"]

	target, err := getImporterTarget("PayLog")
	if err != nil {
		panic(err)
		return nil
	}
	qlen, err := target.Push(string(data))
	if qlen > 0 && err == nil {
		respstr = strconv.FormatInt(qlen, 10)
	} else {
		innerLogger.Error("HttpServer::PayLog push queue data failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
	}
	return nil
}
