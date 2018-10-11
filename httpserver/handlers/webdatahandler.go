package handlers

import (
	"TechPlat/datapipe/global"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var webDataJsonKeys = [...]string{
	"App",
	"Module",
	"DataKey",
	"Data",
	"PageUrl",
	"Remark",
	//"Ver",
	"ClientIP",
}

var webDataUrlKeys = [...]string{
	"app",
	"module",
	"datakey",
	"data",
	"pageurl",
	"remark",
	//"ver",
	"clientip",
}

func WebData(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range webDataUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" || params["datakey"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::WebData " + global.LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::WebData app=[" + params["app"] + "] datakey=[" + params["datakey"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("WebData")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::WebData " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range webDataJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	if dataMap["ClientIP"] == "" {
		dataMap["ClientIP"] = getClientIP(ctx)
	}
	dataMap["UserAgent"] = getUserAgent(ctx)
	dataMap["GlobalID"] = getGlobalID(ctx)
	dataMap["VisitTime"] = getNowFormatTime()
	dataMap["LogID"] = "0" //php版本没有，线上.Net版本有此字段
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::WebData " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::WebData push queue data failed!")
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
		return nil
	}
}
