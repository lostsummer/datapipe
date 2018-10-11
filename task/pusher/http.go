package pusher

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/http"
	"TechPlat/datapipe/util/log"
)

type HttpPusher struct{}

func (h HttpPusher) LogTitle() string {
	return logdefine.LogTitle_HttpHandler
}

func (h HttpPusher) Push(taskConf *config.TaskInfo, val string) error {
	title := h.LogTitle()
	retBody, _, _, httpErr := httputil.HttpPost(taskConf.TargetValue, val, "")
	if httpErr != nil {
		logger.Log(title+"InsertJsonData["+val+"] error -> "+httpErr.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertJsonData success -> ["+val+"] ["+retBody+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}

	return httpErr
}
