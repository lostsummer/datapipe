package handlers

import (
	"TechPlat/datapipe/global"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

var pageClickJsonKeys = [...]string{
	"App",
	"Module",
	"ClickKey",
	"ClickData",
	"ClickRemark",
	"PageUrl",
	"HtmlType",
	"Remark",
	//"Ver",	//php代码中从参数获取，但并没有用在下发json上
}

var pageClickUrlKeys = [...]string{
	"app",
	"module",
	"clickkey",
	"clickdata",
	"clickremark",
	"pageurl",
	"htmltype",
	"remark",
	//"ver",  //php代码中从参数获取，但并没有用在下发json上
}

func PageClick(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range pageClickUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" || params["clickkey"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::PageClick " + global.LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PageClick app=[" + params["app"] + "] clickkey=[" + params["clickkey"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("PageClick")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::PageClick " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range pageClickJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["UserAgent"] = getUserAgent(ctx)
	dataMap["GlobalID"] = getGlobalID(ctx)
	dataMap["FirstVisitTime"] = getFirstVistTime(ctx)
	dataMap["ClientIP"] = getClientIP(ctx)
	dataMap["WriteTime"] = getNowFormatTime()
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
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
		return nil
	}
}
