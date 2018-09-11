package global

import "github.com/devfeel/dottask"

var TaskService *task.TaskService

var (
	Version   = ""
	Branch    = ""
	CommitID  = ""
	BuildTime = ""
)

func init() {
	TaskService = task.StartNewService()
}
