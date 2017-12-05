package tasks

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/log"
	"github.com/devfeel/dottask"
	"TechPlat/datapipe/trigger"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/util/kafka"
	"strconv"
)

// KafkaHandler storage message to kafka
func KafkaHandler(ctx *task.TaskContext) error {
	title := "Tasks:KafkaHandler"
	taskConf := ctx.TaskData.(*config.TaskInfo)

	//get kafka config
	kafkaServerUrl:= config.GetKafkaServerUrl()
	if kafkaServerUrl == ""{
		logger.Log(title+":GetKafkaServerUrl no config kafkaServerUrl", taskConf.TaskID, logdefine.LogLevel_Error)
		return NotConfigError
	}

	//get queue data
	val, err := getQueueData(taskConf)
	if err != nil{
		logger.Log(title+":getRedisData error -> "+err.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return err
	}
	logger.Log(title+":getRedisData -> "+val, taskConf.TaskID, logdefine.LogLevel_Debug)

	partition, offset, kafkaErr:=kafka.SendMessage(kafkaServerUrl, taskConf.TargetValue, val)
	if kafkaErr != nil {
		logger.Log(title+":InsertKafkaData["+val+"] error -> "+kafkaErr.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
	} else {
		logger.Log(title+":InsertKafkaData success -> ["+val+"] ["+
			strconv.Itoa(int(partition))+ "," + strconv.FormatInt(offset, 10)+
			"]", taskConf.TaskID, logdefine.LogLevel_Debug)
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


