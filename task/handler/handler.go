package handler

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/util/log"
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

func Handler(ctx *task.TaskContext) error {
	taskConf := ctx.TaskData.(*config.TaskInfo)
	id := taskConf.ID
	source, exist := config.CurrentConfig.TaskSourceMap[id]
	if !exist {
		logger.Log("no source exist!", taskConf.ID, logdefine.LogLevel_Error)
		return global.NotConfigError
	}
	target, exist := config.CurrentConfig.TaskTargetMap[id]
	if !exist {
		logger.Log("no target exist!", taskConf.ID, logdefine.LogLevel_Error)
		return global.NotConfigError
	}

	val, err := source.Pop()
	if err != nil {
		logger.Log("get source data error -> "+err.Error(), taskConf.ID, logdefine.LogLevel_Error)
		return err
	}

	// 处理 target, push不成功终止后续工作
	if taskConf.HasTargetFilter() && !matchFilter(val, taskConf.Target.Filter) {
		logger.Log("do not insert data  -> ["+val+"] not match filter", taskConf.ID, logdefine.LogLevel_Debug)
	} else {
		_, err := target.Push(val)
		if err != nil {
			logger.Log("insert data ["+val+"] error -> "+err.Error(), taskConf.ID, logdefine.LogLevel_Error)
		} else {
			logger.Log("insert data success ->["+val+"]", taskConf.ID, logdefine.LogLevel_Debug)
		}
	}

	// 处理 trigger, 当前逻辑是返回错误不影响后续工作
	if taskConf.HasTrigger() {
		trigger := config.CurrentConfig.TaskTriggerMap[id]
		if taskConf.HasTriggerFilter() && !matchFilter(val, taskConf.Trigger.Filter) {
			logger.Log("do not trigger data  -> ["+val+"] not match filter", taskConf.ID, logdefine.LogLevel_Debug)
		} else {
			_, err := trigger.Push(val)
			if err != nil {
				logger.Log("trigger error -> "+err.Error(), taskConf.ID, logdefine.LogLevel_Error)
			} else {
				logger.Log("trigger data success ->["+val+"]", taskConf.ID, logdefine.LogLevel_Debug)
			}
		}

	}

	// 处理 counter, 当前逻辑是返回错误不影响后续工作
	if taskConf.HasCounter() {
		err = counter.Count(taskConf.Counter.GetServer(), taskConf.Counter.Key)
		if err != nil {
			logger.Log("Counter error -> "+err.Error(), taskConf.ID, logdefine.LogLevel_Error)
		} else {
			logger.Log("Counter success", taskConf.ID, logdefine.LogLevel_Debug)
		}
	}

	//适应looptask
	return nil
}
