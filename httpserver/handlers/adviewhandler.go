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
	//"ver",   //php代码中有取但没有用
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
	dataMap := make(map[string]interface{})
	for _, k := range adUrlKeys {
		if k != "data" {
			dataMap[k] = params[strings.ToLower(k)]
		}
	}
	var queryData map[string]interface{}
	err = json.Unmarshal([]byte(params["data"]), &queryData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::" + name + " fail to parse post json: " +
			err.Error() + "\r\n" + params["data"] + "\r\n")
		return nil
	}
	for k, v := range queryData {
		dataMap[k] = v
	}
	ua := getUserAgent(ctx)
	dataMap["useragent"] = ua
	dataMap["writetime"] = getNowFormatTime()
	dataMap["os"] = getAgentOS(ua)
	dataMap["browser"] = getAgentBrowser(ua)
	dataMap["clientip"] = getClientIP(ctx)
	if _, exist := dataMap["appid"]; !exist {
		dataMap["appid"] = ""
	}
	if _, exist := dataMap["logtype"]; !exist {
		dataMap["logtype"] = ""
	}
	// 唯独mid类型为数字
	if _, exist := dataMap["mid"]; !exist {
		dataMap["mid"] = 0
	}
	if _, exist := dataMap["pid"]; !exist {
		dataMap["pid"] = ""
	}
	if _, exist := dataMap["sid"]; !exist {
		dataMap["sid"] = ""
	}
	if _, exist := dataMap["tid"]; !exist {
		dataMap["tid"] = ""
	}
	if _, exist := dataMap["uid"]; !exist {
		dataMap["uid"] = ""
	}
	if _, exist := dataMap["uname"]; !exist {
		dataMap["uname"] = ""
	}
	if _, exist := dataMap["adcode"]; !exist {
		dataMap["adcode"] = ""
	}
	if _, exist := dataMap["targeturl"]; !exist {
		dataMap["targeturl"] = ""
	}
	if _, exist := dataMap["pageurl"]; !exist {
		dataMap["pageurl"] = ""
	}
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
			if err != nil {
				innerLogger.Error(err.Error())
			}
			respstr = respFailed
		}
		return nil
	}
}

func ADView(ctx dotweb.Context) error {
	return adbase(ctx, "ADView")
}
