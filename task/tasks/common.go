package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"
	"encoding/json"
	"strings"

	"github.com/devfeel/dottask"
	"github.com/pkg/errors"
)

// 适配各种类型target推送逻辑
type Pusher interface {
	LogTitle() string
	Push(taskConf *config.TaskInfo, val string) error
}

type PusherBase struct {
	Title string
}

var (
	innerLogger *logger.InnerLogger
)

var (
	NotConfigError = errors.New("no exists such config info")
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

// getQueueData get queue data from redis
func getQueueData(taskConf *config.TaskInfo) (string, error) {
	redisClient := redisutil.GetRedisClient(taskConf.FromServer)
	if redisClient == nil {
		return "", NotConfigError
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

func dealTrigger(logtitle string, taskConf *config.TaskInfo, val string) error {
	var err error
	if taskConf.HasTrigger() {
		if taskConf.HasTriggerFilter() && !matchFilter(val, taskConf.TriggerFilter) {
			logger.Log(logtitle+":Do not send TriggerSignal -> ["+val+"] not match filter", taskConf.TaskID, logdefine.LogLevel_Debug)
			return nil
		}
		err = trigger.SendSignal(taskConf.TriggerServer, taskConf.TriggerQueue, val)
		if err != nil {
			logger.Log(logtitle+":SendTriggerSignal error -> ["+val+"] "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		} else {
			logger.Log(logtitle+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}
	return err
}

func dealCounter(logtitle string, taskConf *config.TaskInfo) error {
	var err error
	if taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(logtitle+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		} else {
			logger.Log(logtitle+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)

		}
	}
	return err
}

// after inserted data into mongodb/kafka/http, deal trigger and counter by config
func afterWork(logtitle string, taskConf *config.TaskInfo, val string) error {
	if err := dealTrigger(logtitle, taskConf, val); err == nil {
		return dealCounter(logtitle, taskConf)
	} else {
		return err
	}
}

func handler(ctx *task.TaskContext, pusher Pusher) error {
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

	if taskConf.HasTargetFilter() && !matchFilter(val, taskConf.TargetFilter) {
		logger.Log(title+":Do not insert JsonData  -> ["+val+"] not match filter", taskConf.TaskID, logdefine.LogLevel_Debug)
	} else {
		pusher.Push(taskConf, val)
	}

	return afterWork(title, taskConf, val)
}
