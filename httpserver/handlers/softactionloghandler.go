package handlers

import (
	"TechPlat/datapipe/endpoint"
	"TechPlat/datapipe/global"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb"
)

type Log struct {
	Col_1              string `json:"col_1"`
	Col_2              string `json:"col_2"`
	Pid                string `json:"pid"`
	User_ID            string `json:"user_id"`
	Seq_Order          string `json:"seq_order"`
	Oper_Type          string `json:"oper_type"`
	Name               string `json:"name"`
	Parameter          string `json:"parameter"`
	Interface_Name     string `json:"interface_name"`
	User_Name          string `json:"user_name"`
	Client_Version     string `json:"client_version"`
	Log_Date           string `json:"log_date"`
	Session_Start_Time string `json:"session_start_time"`
	Mac                string `json:"mac"`
}

//type ActionData struct {
//	Logs []Log `json:"logs"`
//}

const (
	freeuserQueuePostfix = "_FreeUser"
)

var freeUserPids = [...]string{
	"200000301",
	"200000302",
	"200000000",
	"100002000",
	"200000303",
	"202002000",
	"202002300",
}

func isFreeUserPid(pid string) bool {
	for _, p := range freeUserPids {
		if p == pid {
			return true
		}
	}
	return false
}

func SoftActionLog(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::SoftActionLog")
		ctx.WriteString(respstr)
	}()

	datajson := ctx.PostFormValue(postActionDataKey)
	if datajson == "" {
		innerLogger.Error("HttpServer::SoftActionLog " + global.LessParamError.Error())
		respstr = respFailed
		return nil
	}
	//屏蔽客户端传非法json的问题
	datajson = strings.Replace(datajson, "[,", "[", 1)
	datajson = strings.Replace(datajson, "\t", "  ", -1)
	datajson = strings.Replace(datajson, "\r", "", -1)
	datajson = strings.Replace(datajson, "\n", "", -1)

	var actionData []Log
	err := json.Unmarshal([]byte(datajson), &actionData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::SoftActionLog fail to parse post json actionData: " +
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
		if len(dataMap["name"]) > 1000 {
			dataMap["name"] = dataMap["name"][:1000]
		}
		dataMap["writetime"] = getNowFormatTime()
		dataMap["client_ip"] = getClientIP(ctx)
		if data, err := json.Marshal(dataMap); err != nil {
			respstr = respFailed
			innerLogger.Error("HttpServer::SoftActionLog " + err.Error())
			return nil
		} else {
			var qlen int64
			var err error

			target, err := getImporterTarget("SoftActionLog")
			if err != nil {
				panic(err)
				return nil
			}

			// 此处免费用户写入不同redis队列的逻辑继承自php版本
			// 因为不能体现在配置上我十分想移至task中借助trigger filter过滤
			// 但filter目前不支持反向过滤，这使得付费队列中没法滤除免费用户数据
			// 今后考虑filter支持正则匹配，类似mongodb的find语法
			if isFreeUserPid(dataMap["pid"]) {
				r, ok := target.(*endpoint.Redis)
				if ok {
					rFree := new(endpoint.Redis)
					*rFree = *r
					rFree.Key += freeuserQueuePostfix
					target = rFree
				}
			}

			qlen, err = target.Push(string(data))
			if qlen > 0 && err == nil {
				respstr = strconv.FormatInt(qlen, 10)
			} else {
				innerLogger.Error("HttpServer::SoftActionLog push queue data failed!")
				if err != nil {
					innerLogger.Error(err.Error())
				}
				respstr = respFailed
			}

		}
	}
	return nil
}
