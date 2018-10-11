package task

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/task/handler"
	"TechPlat/datapipe/util/log"
	"fmt"

	"github.com/devfeel/dottask"
)

var (
	innerLogger *logger.InnerLogger
	taskService *task.TaskService
)

func init() {
	innerLogger = logger.GetInnerLogger()
	taskService = global.TaskService
}

const (
	taskDueTime  = 0
	taskInterval = 1 //1毫秒
)

func LoadTasks(service *task.TaskService) {
	innerLogger.Info("Task::RegisterTask begin...")
	for _, v := range config.CurrentConfig.TaskMap {
		taskHandler := handler.CreateHandler(v.TargetType)
		if taskHandler != nil {
			service.CreateLoopTask(v.TaskID, true, taskDueTime,
				taskInterval, taskHandler, v)
		} else {
			innerLogger.Error("Task::CreateTask fail to create task type " +
				v.TargetType)
		}
	}

	//load queue task
	global.TaskService.CreateQueueTask(counter.QueueTaskName, true, 1, counter.DealMessage, nil, counter.QueueSize)
	innerLogger.Info("Task::RegisterTask end")
}

func StartTaskService() {
	//step 1: LoadTasks
	LoadTasks(taskService)

	//step 2: start all task
	taskService.StartAllTask()

	fmt.Println("StartTaskService", taskService.PrintAllCronTask())
}

func ReStartTaskService() {
	//step 1: stop and remove all task
	taskService.RemoveAllTask()

	//step 2: LoadTasks
	LoadTasks(taskService)

	//step 3: start all task
	taskService.StartAllTask()

	fmt.Println("ReStartTaskService", taskService.PrintAllCronTask())
}
