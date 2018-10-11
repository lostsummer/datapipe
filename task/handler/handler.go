package handler

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/task/pusher"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"
	"encoding/json"
	"strings"

	"github.com/devfeel/dottask"
)

var (
	innerLogger *logger.InnerLogger
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

// getQueueData get queue data from redis
func getQueueData(taskConf *config.TaskInfo) (string, error) {
	redisClient := redisutil.GetRedisClient(taskConf.FromServer)
	if redisClient == nil {
		return "", global.NotConfigError
	}
	val, err := redisClient.BRPop(taskConf.FromQueue)
	return val, err
}

// 暂时支持单个string 类型字段的完全匹配, "|" 分隔多个匹配值
// 例如 triggerfilter="App=12345|12354"
// 非string/key不存在, 都认为无匹配返回false
func matchFilter(val string, fltStr string) bool {
	var jsMap map[string]interface{}
	if err := json.Unmarshal([]byte(val), &jsMap); err != nil {
		return false
	}
	fKeyValues := strings.Split(fltStr, "=")
	key := fKeyValues[0]
	fValues := strings.Split(fKeyValues[1], "|")
	jsValueIf, exist := jsMap[key]
	if !exist {
		return false
	}
	jsValue := ""
	switch jsValueIf.(type) {
	case string:
		jsValue = jsValueIf.(string)
	default:
		return false
	}
	for _, v := range fValues {
		if v == jsValue {
			return true
		}
	}
	return false
}

func handlerWork(ctx *task.TaskContext, pusher pusher.Pusher) error {
	title := pusher.LogTitle()
	taskConf := ctx.TaskData.(*config.TaskInfo)
	// 今后开关加载提到外部
	if taskConf.Enable == false {
		return nil
	}

	val, err := getQueueData(taskConf)
	if err != nil {
		logger.Log(title+":getRedisData error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}

	logger.Log(title+":getRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)

	// 处理 target, 当前逻辑是返回错误不影响后续工作
	if taskConf.HasTargetFilter() && !matchFilter(val, taskConf.TargetFilter) {
		logger.Log(title+":Do not insert JsonData  -> ["+val+"] not match filter", taskConf.TaskID, logdefine.LogLevel_Debug)
	} else {
		err = pusher.Push(taskConf, val)
	}

	// 处理 trigger, 当前逻辑是返回错误不影响后续工作
	if taskConf.HasTrigger() {
		if taskConf.HasTriggerFilter() && !matchFilter(val, taskConf.TriggerFilter) {
			logger.Log(title+":Do not send TriggerSignal -> ["+val+"] not match filter", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
		err = trigger.SendSignal(taskConf.TriggerServer, taskConf.TriggerQueue, val)
		if err != nil {
			logger.Log(title+":SendTriggerSignal error -> ["+val+"] "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		} else {
			logger.Log(title+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	// 处理 counter, 当前逻辑是返回错误不影响后续工作
	if taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(title+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		} else {
			logger.Log(title+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	return err
}

func CreateHandler(taskType string) task.TaskHandle {
	var p pusher.Pusher
	switch taskType {
	case config.Target_MongoDB:
		p = pusher.MongoDBPusher{}
	case config.Target_Kafka:
		p = pusher.KafkaPusher{}
	case config.Target_Http:
		p = pusher.HttpPusher{}
	default:
		return nil
	}
	return func(ctx *task.TaskContext) error {
		handlerWork(ctx, p)

		//适应looptask
		return nil
	}
}
