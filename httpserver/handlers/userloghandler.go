package handlers

import (
	"TechPlat/datapipe/global"
	"strconv"

	"github.com/devfeel/dotweb"
)

var userLogUrlKeys = [...]string{
	"appid",
	"logtype",
	"ver",
}

//对php代码照搬然而有疑惑，appid, logtype, ver并没有实际用途
func userlogbase(ctx dotweb.Context, name string) error {
	params := make(map[string]string)
	respstr := respFailed
	for _, k := range userLogUrlKeys {
		params[k] = ctx.Request().QueryString(k)
	}
	jsondata := ctx.Request().PostBody()
	if params["appid"] == "" || params["logtype"] == "" {
		ctx.WriteString(respFailed)
		innerLogger.Error("HttpServer::" + name + " " + global.LessParamError.Error())
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
	qlen, err := pushQueueData(importerConf, string(jsondata))
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

func UserLog(ctx dotweb.Context) error {
	return userlogbase(ctx, "UserLog")
}
