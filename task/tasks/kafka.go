package tasks

import (
	"TechPlat/datapipe/component/kafka"
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/const/log"
	"TechPlat/datapipe/util/log"
	"strconv"

	"github.com/devfeel/dottask"
)

type KafkaPusher PusherBase

const (
	kafkaLogTitle = "Tasks:KafkaHandler"
)

func (k KafkaPusher) LogTitle() string {
	return k.Title
}

func (k KafkaPusher) Push(taskConf *config.TaskInfo, val string) error {
	title := k.Title
	// 今后将检查配置的逻辑移至外部，启动加载时检查一次
	kafkaServerUrl := config.GetKafkaServerUrl()
	if kafkaServerUrl == "" {
		logger.Log(title+":GetKafkaServerUrl no config kafkaServerUrl", taskConf.TaskID, logdefine.LogLevel_Error)
		return NotConfigError
	}
	partition, offset, kafkaErr := kafka.SendMessage(kafkaServerUrl, taskConf.TargetValue, val)
	if kafkaErr != nil {
		logger.Log(title+":InsertKafkaData["+val+"] error -> "+kafkaErr.Error(), taskConf.TaskID, logdefine.LogLevel_Error)
		return kafkaErr
	} else {
		logger.Log(title+":InsertKafkaData success -> ["+val+"] ["+
			strconv.Itoa(int(partition))+","+strconv.FormatInt(offset, 10)+
			"]", taskConf.TaskID, logdefine.LogLevel_Debug)
	}
	return nil
}

// KafkaHandler storage message to kafka
func KafkaHandler(ctx *task.TaskContext) error {
	handler(ctx, KafkaPusher{kafkaLogTitle})
	return nil

}
