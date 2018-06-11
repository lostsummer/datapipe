package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var pageViewJsonKeys = [...]string{
	"Code",
	"PageUrl",
	"ReferUrl",
	"App",
	"Module",
	"Remark",
	"Ver",
}

var pageViewUrlKeys = [...]string{
	"code",
	"pageurl",
	"referurl",
	"app",
	"module",
	"remark",
	"ver",
}

func PageView(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	//request := ctx.Request()
	for _, k := range pageViewUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["code"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::PageView " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PageView code=[" + params["code"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("PageView")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PageView " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range pageViewJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["UserAgent"] = getUserAgent(ctx)
	dataMap["GlobalID"] = getGlobalID(ctx)
	dataMap["FirstVisitTime"] = getFirstVistTime(ctx)
	dataMap["ClientIP"] = getClientIP(ctx)
	dataMap["WriteTime"] = getNowFormatTime()
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PageView " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::PageView push queue data failed!")
			respstr = respFailed
		}
		return nil
	}
}
