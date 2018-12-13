package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/devfeel/dotweb"
)

type HPSoftLog struct {
	Uid    string   `json:"uid"`
	Pid    string   `json:"pid"`
	Values []string `json:"values"`
}

type HPSoftLogToQ struct {
	Uid       string `json:"uid"`
	Pid       string `json:"pid"`
	Value     string `json:"value"`
	WriteTime string `json:"writetime"`
}

func HPSoftLogHandle(ctx dotweb.Context) error {
	respstr := respFailed
	defer func() {
		innerLogger.Info("HttpServer::HPSoftLog")
		ctx.WriteString(respstr)
	}()

	jsondata := ctx.Request().PostBody()

	var log HPSoftLog
	if err := json.Unmarshal([]byte(jsondata), &log); err != nil {
		innerLogger.Error("HttpServer::HPSoftLog " + err.Error())
		return nil
	}

	for _, v := range log.Values {
		dataToQ := HPSoftLogToQ{
			Uid:       log.Uid,
			Pid:       log.Pid,
			Value:     v,
			WriteTime: getNowFormatTime(),
		}
		if outJson, err := json.Marshal(dataToQ); err != nil {
			innerLogger.Error("HttpServer::HPSoftLog " + err.Error())
			return nil
		} else if target, err := getImporterTarget("HPSoftLog"); err != nil {
			panic(err)
			return nil
		} else if _, err := target.Push(string(outJson)); err != nil {
			innerLogger.Error("HttpServer::HPSoftLog push queue data failed!")
			innerLogger.Error(err.Error())
			return nil
		}
	}
	respstr = fmt.Sprintf("%d", len(log.Values))
	return nil
}
