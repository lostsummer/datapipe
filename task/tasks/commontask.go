package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/http"
	"TechPlat/datapipe/util/log"
	"TechPlat/datapipe/util/redis"
	"TechPlat/datapipe/repository/impl"
	"github.com/devfeel/dottask"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/counter"
)

var (
	innerLogger *logger.InnerLogger
)

func init() {
	innerLogger = logger.GetInnerLogger()
}



//将redis获取的数据存入mongodb
func MongoDBHandler(ctx *task.TaskContext) error{
	title:= "Tasks:MongoDBHandler"
	taskConf := ctx.TaskData.(*config.TaskInfo)

	redisClient := redisutil.GetRedisClient(taskConf.FromServer)
	val, err := redisClient.BRPop(taskConf.FromQueue)
	if err != nil {
		logger.Log(title+":BRPop error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":GetRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)

	mongoHandler := new(impl.MongoHandler)
	mongoHandler.SetConn(config.CurrentConfig.MongoDB.ServerIP, config.CurrentConfig.MongoDB.DBName)
	//insert data to mongo
	err = mongoHandler.InsertJsonData(taskConf.TargetValue, val)
	if err != nil {
		logger.Log(title+":InsertJsonData ["+val+"] error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertJsonData success ->["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}

	//deal trigger
	if err == nil && taskConf.HasTrigger() {
		err = trigger.SendSignal(taskConf.TriggerServer, taskConf.TriggerQueue, val)
		if err != nil {
			logger.Log(title+":SendTriggerSignal error -> ["+val+"] "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	//deal counter
	if err == nil && taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(title+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)

		}
	}
	return nil
}

//将redis获取的数据发送到指定http接口，默认post
func HttpHandler(ctx *task.TaskContext) error{
	title:= "Tasks:HttpHandler"
	taskConf := ctx.TaskData.(*config.TaskInfo)

	redisClient := redisutil.GetRedisClient(taskConf.FromServer)
	val, err := redisClient.BRPop(taskConf.FromQueue)
	if err != nil {
		logger.Log(title+":BRPop error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":GetRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)

	//insert data to HttpPost
	retBody, _, _, httpErr := httputil.HttpPost(taskConf.TargetValue, val, "")
	if httpErr != nil {
		logger.Log(title+"InsertJsonData["+val+"] error -> "+httpErr.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertJsonData success -> ["+val+"] ["+retBody+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}

	//deal trigger
	if err == nil && taskConf.HasTrigger() {
		err = trigger.SendSignal(taskConf.TriggerServer, taskConf.TriggerQueue, val)
		if err != nil {
			logger.Log(title+":SendTriggerSignal error -> ["+val+"] "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":SendTriggerSignal success -> ["+val+"]", taskConf.TaskID, logdefine.LogLevel_Debug)
		}
	}

	//deal counter
	if err == nil && taskConf.HasCounter() {
		err = counter.Count(taskConf.CounterServer, taskConf.CounterKey)
		if err != nil {
			logger.Log(title+":Counter error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		}else{
			logger.Log(title+":Counter success", taskConf.TaskID, logdefine.LogLevel_Debug)

		}
	}
	return nil
}
