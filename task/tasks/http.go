package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/http"
	"TechPlat/datapipe/util/log"

	"github.com/devfeel/dottask"
)

type HttpPusher PusherBase

const (
	httpLogTitle = "Tasks:HttpHandler"
)

func (h HttpPusher) LogTitle() string {
	return h.Title
}

func (h HttpPusher) Push(taskConf *config.TaskInfo, val string) error {
	title := h.Title
	retBody, _, _, httpErr := httputil.HttpPost(taskConf.TargetValue, val, "")
	if httpErr != nil {
		logger.Log(title+"InsertJsonData["+val+"] error -> "+httpErr.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertJsonData success -> ["+val+"] ["+retBody+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}

	return httpErr
}

//将redis获取的数据发送到指定http接口，默认post
func HttpHandler(ctx *task.TaskContext) error {
	handler(ctx, HttpPusher{httpLogTitle})
	return nil
}
