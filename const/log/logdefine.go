package logdefine

const (
	LogTarget_Default = "Default"
	LogTarget_Task    = "Task"
	LogTarget_MongoDB = "MongoDB"

	LogLevel_Debug = "debug"
	LogLevel_Info  = "info"
	LogLevel_Warn  = "warn"
	LogLevel_Error = "error"
)

const (
	LogTitle_HttpSource     = "Source:Http"
	LogTitle_KafkaSource    = "Source:Kafka"
	LogTitle_MongoDBSource  = "Source:MongDB"
	LogTitle_RedisSource    = "Source:Redis"
	LogTitle_HttpTarget     = "Target:Http"
	LogTitle_KafkaTarget    = "Target:Kafka"
	LogTitle_MongoDBTarget  = "Target:MongDB"
	LogTitle_RedisTarget    = "Target:Redis"
	LogTitle_HttpTrigger    = "Trigger:Http"
	LogTitle_KafkaTrigger   = "Trigger:Kafka"
	LogTitle_MongoDBTrigger = "Trigger:MongDB"
	LogTitle_RedisTrigger   = "Trigger:Redis"
)
