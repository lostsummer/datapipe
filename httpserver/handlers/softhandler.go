package handlers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

const (
	defaultSoftQueueKeyForApp = "EMoney.DataPipe:SoftLog_ForApp"
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

var softEncryptUrlKeys = [...]string{
	"app",
	"version",
	"token",
}

var softEncryptJsonKeys = [...]string{
	"App",
	"Version",
	"Token",
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
			importerConf_, err := getImporterConf("SoftLogForApp")
			if err != nil { // no SoftLogApp importer config
				qlen, err = pushQueueDataToSQ(importerConf.ServerUrl,
					defaultSoftQueueKeyForApp,
					string(data))
			} else {
				qlen, err = pushQueueData(importerConf_, string(data))
			}
			if err != nil {
				innerLogger.Error("HttpServer::Soft push queue data failed!")
				if err != nil {
					innerLogger.Error(err.Error())
				}
				respstr = respFailed
				return nil

			}
		}
		qlen, err = pushQueueData(importerConf, string(data))
		if err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::Soft push queue data failed!")
			innerLogger.Error(err.Error())
			respstr = respFailed
		}
		return nil
	}
}

func SoftEncrypt(ctx dotweb.Context) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range softEncryptUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	if params["app"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::SoftEncrypt " + LessParamError.Error())
		return nil
	}
	defer func() {
		innerLogger.Info("HttpServer::SoftEncrypt app=[" + params["app"] + "]")
		ctx.WriteString(respstr)
	}()
	importerConf, err := getImporterConf("SoftEncrypt")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::SoftEncrypt " + err.Error())
		return nil
	}
	dataMap := make(map[string]string)
	for _, k := range softEncryptJsonKeys {
		dataMap[k] = params[strings.ToLower(k)]
	}
	dataMap["Decrypt"] = ""
	dataMap["PageUrl"] = getFullUrl(ctx)
	dataMap["ClientIP"] = getClientIP(ctx)
	dataMap["GlobalID"] = getGlobalID(ctx)

	if data, err := json.Marshal(dataMap); err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::SoftEncrypt " + err.Error())
		return nil
	} else {
		qlen, err := pushQueueData(importerConf, string(data))
		if err == nil {
			respstr = strconv.FormatInt(qlen, 10)
		} else {
			innerLogger.Error("HttpServer::SoftEncrypt push queue data failed!")
			innerLogger.Error(err.Error())
			respstr = respFailed
		}
		return nil
	}
}
