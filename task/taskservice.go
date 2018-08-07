package task

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/task/tasks"
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
		innerLogger.Info("Task::RegisterTask::RegisterTask => " + v.TaskID)
		switch v.TargetType {
		case config.Target_MongoDB:
			service.CreateLoopTask(v.TaskID, true, taskDueTime, taskInterval, tasks.MongoDBHandler, v)
		case config.Target_Http:
			service.CreateLoopTask(v.TaskID, true, taskDueTime, taskInterval, tasks.HttpHandler, v)
		case config.Target_Kafka:
			service.CreateLoopTask(v.TaskID, true, taskDueTime, taskInterval, tasks.KafkaHandler, v)
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
