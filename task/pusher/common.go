package pusher

import (
	"TechPlat/datapipe/config"
)

// 适配各种类型target推送逻辑
type Pusher interface {
	LogTitle() string
	Push(taskConf *config.TaskInfo, val string) error
}
