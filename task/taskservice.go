package task

import (
	"github.com/devfeel/dottask"
	"emoney/tongjiservice/config"
	"fmt"
	"emoney/tongjiservice/util/log"
	"emoney/tongjiservice/task/tasks"
)

var (
	innerLogger *logger.InnerLogger
	taskService *task.TaskService
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

const(
	taskDueTime = 0
	taskInterval = 1 //1毫秒
)

func LoadTasks(service *task.TaskService) {
	innerLogger.Info("Task::RegisterTask begin...")
	for _, v := range config.CurrentConfig.TaskMap {
		innerLogger.Info("Task::RegisterTask::RegisterTask => " + v.TaskID)
		if v.TargetType == config.Target_MongoDB{
			service.CreateLoopTask(v.TaskID, true, taskDueTime, taskInterval, tasks.MongoDBHandler, v)
		}else if v.TargetType == config.Target_Http{
			service.CreateLoopTask(v.TaskID, true, taskDueTime, taskInterval, tasks.HttpHandler, v)
		}
	}
	innerLogger.Info("Task::RegisterTask end")
}

func StartTaskService() {
	taskService = task.StartNewService()

	//step 1: LoadTasks
	LoadTasks(taskService)

	//step 2: start all task
	taskService.StartAllTask()

	fmt.Println("StartTaskService", taskService.PrintAllCronTask())
}

func ReStartTaskService(){
	//step 1: stop and remove all task
	taskService.RemoveAllTask()

	//step 2: LoadTasks
	LoadTasks(taskService)

	//step 3: start all task
	taskService.StartAllTask()

	fmt.Println("ReStartTaskService", taskService.PrintAllCronTask())
}