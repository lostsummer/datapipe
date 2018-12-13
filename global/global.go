package global

import (
	"github.com/devfeel/dottask"
	"github.com/pkg/errors"
)

var TaskService *task.TaskService

var (
	Version   = ""
	Branch    = ""
	CommitID  = ""
	BuildTime = ""
)

var (
	NotConfigError = errors.New("not exists such config info")
	LessParamError = errors.New("less param")
	GetRedisError  = errors.New("get rediscli failed")
	EscapeError    = errors.New("invalid URL escape")
)

func init() {
	TaskService = task.StartNewService()
}
