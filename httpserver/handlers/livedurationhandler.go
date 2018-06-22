package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var liveDuraJsonKeys = [...]string{
	"App",
	"Uid",
	"Pid",
	"ClassID",
	"Remark",
}

var liveDuraUrlKeys = [...]string{
	"app",
	"uid",
	"pid",
	"classid",
	"remark",
}

func LiveDuration(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range liveDuraUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" || params["uid"] == "" || params["classid"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::LiveDuration " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::LiveDuration app=[" + params["app"] +
			"] uid=[" + params["uid"] + "]" +
			"] classid=[" + params["classid"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("LiveDuration")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::LiveDuration " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range liveDuraJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["BeginTime"] = getNowFormatTime()
	dataMap["EndTime"] = getNowFormatTime()
	dataMap["WriteTime"] = getNowFormatTime()
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::LiveDuration " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::LiveDuration push queue data failed!")
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
		return nil
	}
}
