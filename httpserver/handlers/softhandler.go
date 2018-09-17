package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

const (
	softQueueKeyForApp = "EMoney.DataPipe:SoftLog_ForApp"
)

var softJsonKeys = [...]string{
	"App",
	"Version",
	"OPType",
	"Flag",
	"Remark",
}

var softUrlKeys = [...]string{
	"app",
	"version",
	"optype",
	"flag",
	"remark",
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
	importerConf, err := getImporterConf("Soft")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Soft " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range softJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["PageUrl"] = getFullUrl(ctx)
	//dataMap["UserAgent"] = getUserAgent(ctx)  //php版本中有，线上.Net版本无

	dataMap["GlobalID"] = getGlobalID(ctx)
	//dataMap["FirstVisitTime"] = getFirstVistTime(ctx) //php版本中有，线上.Net版本无

	dataMap["ClientIP"] = getClientIP(ctx)
	dataMap["LogID"] = "0" //php版本无， 线上.Net版本有

	//php版本无， 线上.Net版本有
	dataMap["WriteTime"] = getNowFormatTime()
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Soft " + err.Error())
		return nil
	} else {
		var qlen int64
		var err error
		if params["remark"] == "ShowUnInstallBK" {
			qlen, err = pushQueueDataToSQ(importerConf.ServerUrl,
				softQueueKeyForApp,
				string(data))
		}
		qlen, err = pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::Soft push queue data failed!")
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
		return nil
	}
}
