package global

import "github.com/devfeel/dottask"

var TaskService *task.TaskService

func init(){
	TaskService = task.StartNewService()
}

