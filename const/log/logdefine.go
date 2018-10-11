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
	LogTitle_HttpHandler    = "Tasks:HttpHandler"
	LogTitle_KafkaHandler   = "Tasks:KafkaHandler"
	LogTitle_MongoDBHandler = "Tasks:MongDBHandler"
	LogTitle_RedisHandler   = "Tasks:RedisHandler"
)
