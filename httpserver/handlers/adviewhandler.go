package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

// adview 中url参数和jsonkey统一到小写，而pageview就不是，这是原先php代码比较随意而不一致之处
var adUrlKeys = [...]string{
	"appid",
	"logtype",
	"data",
	"referurl",
	"ver",
}

var adExtraJsonKeys = [...]string{
	"mid",
	"pid",
	"sid",
	"tid",
	"uid",
	"uname",
	"adcode",
	"targeturl",
	"pageurl",
}

func adbase(ctx dotweb.Context, name string) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range adUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["appid"] == "" || params["logtype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::" + name + " " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::" + name + " appid=[" + params["appid"] + "] logtype=[" + params["logtype"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf(name)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::" + name + " " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range adUrlKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	for _, k := range adExtraJsonKeys {
		dataMap[k] = "" //php代码中都留空, 这是一个疑惑点，先逻辑照搬
	}
	ua := getUserAgent(ctx)
	dataMap["useragent"] = ua
	dataMap["writetime"] = getNowFormatTime()
	dataMap["os"] = getAgentOS(ua)
	dataMap["browser"] = getAgentBrowser(ua)
	dataMap["clientip"] = getClientIP(ctx)
	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::" + name + " " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if qlen > 0 && err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::" + name + " push queue data failed!")
			respstr = respFailed
		}
		return nil
	}
}

func ADView(ctx dotweb.Context) error {
	return adbase(ctx, "ADView")
}
