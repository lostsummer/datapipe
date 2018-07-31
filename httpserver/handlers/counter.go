package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/devfeel/dotweb"
)

type IncrData struct {
	Category string `json:"category"`
	AppID    string `json:"appid"`
	Key      string `json:"key"`
	Increase int    `json:"increase"`
	Time     int64  `json:"time"`
}

func Counter(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::Counter")
		ctx.WriteString(respstr)
	}()

	accConf, err := getAccumulatorConf("Counter")
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Counter " + err.Error())
		return nil
	}

	datajson := ctx.PostFormValue(postActionDataKey)
	if datajson == "" {
		innerLogger.Error("HttpServer::Counter " + LessParamError.Error())
		respstr = respFailed
		return nil
	}

	var incrData IncrData
	err = json.Unmarshal([]byte(datajson), &incrData)
	if err != nil {
		respstr = respFailed
		innerLogger.Error("HttpServer::Counter fail to parse post json incrData: " +
			err.Error() + "\r\n" + datajson + "\r\n")
		return nil
	}

	category, appid, key, val := incrData.Category, incrData.AppID, incrData.Key, incrData.Increase
	field := fmt.Sprintf("%s:%s", appid, key)
	ret, err := counterIncrBy(accConf, category, field, val)

	if ret > 0 && err == nil {
		respstr = strconv.FormatInt(ret, 10)
	} else {
		innerLogger.Error("HttpServer::Counter incr failed!")
		if err != nil {
			innerLogger.Error(err.Error())
		}
		respstr = respFailed
	}
	return nil
}
