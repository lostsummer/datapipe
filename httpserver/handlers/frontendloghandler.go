package handlers

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

type FELog struct {
	ID      string `json:"id"`
	AppID   string `json:"appid"`
	Module  string `json:"module"`
	User_ID string `json:"user_id"`
	MsgList string `json:"msglist"`
	Level   string `json:"level"`
	Target  string `json:"target"`
	RowNum  string `json:"rowNum"`
	ColNum  string `json:"colNum"`
	From    string `json:"from"`
	Msg     string `json:"msg"`
	Remark  string `json:"remark"`
	Type    string `json:"type"`
}

type FEActionData struct {
	Logs []FELog `json:"logs"`
}

func FrontEndLog(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::FrontEndLog")
		ctx.WriteString(respstr)
	}()

	importerConf, err := getImporterConf("FrontEndLog")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::FrontEndLog " + err.Error())
		return nil
	}

	datajson := ctx.PostFormValue(postActionDataKey)
	if datajson == "" {
		innerLogger.Error("HttpServer::FrontEndLog " + LessParamError.Error())
		respstr = respFailed
		return nil
	}
	//屏蔽客户端传非法json的问题 ([, {}, {}])
	datajson = strings.Replace(datajson, "[,", "[", 1)
	datajson = strings.Replace(datajson, "\t", "  ", 1)
	datajson = strings.Replace(datajson, "\r", "", 1)
	datajson = strings.Replace(datajson, "\n", "", 1)

	var actionData []FELog
	err = json.Unmarshal([]byte(datajson), &actionData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::FrontEndLog fail to parse post json actionData: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	}

	for _, log := range actionData {
		t := reflect.TypeOf(log)
		v := reflect.ValueOf(log)
		dataMap := make(map[string]string)
		for i := 0; i < t.NumField(); i++ {
			key := strings.ToLower(t.Field(i).Name)
			val := v.Field(i).Interface()
			dataMap[key] = val.(string)

		}
		dataMap["ClientIP"] = getClientIP(ctx)
		dataMap["WriteTime"] = getNowFormatTime()
		if data, err := json.Marshal(dataMap); err != nil {
			respstr = respFailed
			innerLogger.Error("HttpServer::FrontEndLog " + err.Error())
			return nil
		} else {
			qlen, err := pushQueueData(importerConf, string(data))
			if qlen > 0 && err == nil {
				respstr = strconv.FormatInt(qlen, 10)
			} else {
				innerLogger.Error("HttpServer::FrontEndLog push queue data failed!")
				if err != nil {
					innerLogger.Error(err.Error())
				}
				respstr = respFailed
			}
		}
	}
	return nil
}
