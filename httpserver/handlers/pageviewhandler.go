package handlers

import (
	"TechPlat/datapipe/global"
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
	//"Ver",   //php代码中从参数获取，但并没有用在下发json上
}

var pageViewUrlKeys = [...]string{
	"code",
	"pageurl",
	"referurl",
	"app",
	"module",
	"remark",
	//"ver",   //php代码中从参数获取，但并没有用在下发json上
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
		innerLogger.Error("HttpServer::PageView " + global.LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::PageView code=[" + params["code"] + "]")
		ctx.WriteString(respstr)
	}()
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
	} else {
		target, err := getImporterTarget("PageView")
		if err != nil {
			panic(err)
			return nil
		}
		qlen, err := target.Push(string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::PageView push queue data failed! ")
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
	}
	return nil
}
