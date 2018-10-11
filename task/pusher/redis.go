package pusher

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/queue"
	"TechPlat/datapipe/util/log"
)

type RedisPusher struct{}

func (r RedisPusher) LogTitle() string {
	return logdefine.LogTitle_RedisHandler
}

func (r RedisPusher) Push(taskConf *config.TaskInfo, val string) error {
	title := r.LogTitle()
	redisName := taskConf.TargetName
	var redisConf *config.Redis
	for _, r := range config.CurrentConfig.Redises.RedisList {
		if r.Name == redisName {
			redisConf = &r
		}
	}

	// 今后配置检查提到外部
	if redisConf == nil {
		logger.Log(title+":GetConfig no "+redisName+" define", taskConf.TaskID, logdefine.LogLevel_Error)
		return global.NotConfigError
	}

	q := &queue.Queue{
		redisConf.ServerUrl,
		0,
		taskConf.TargetValue,
	}

	_, err := q.Push(val)

	if err != nil {
		logger.Log(title+":InsertJsonData ["+val+"] error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":InsertJsonData success ->["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	return nil
}
