/*emoney.tongjiservice
* Author: Panxinming
* LastUpdateTime: 2016-10-17 10:00
 */
package main

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/core/exception"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/httpserver"
	"TechPlat/datapipe/task"
	"TechPlat/datapipe/util/common"
	"TechPlat/datapipe/util/log"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var (
	innerLogger *logger.InnerLogger
	configPath  string
	configFile  string
	basePath    string
)

func init() {
	innerLogger = logger.GetInnerLogger()
	basePath = common.GetCurrentDirectory()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError("DataPipe main->recover error", err)
			innerLogger.Error(ex.GetDefaultLogString())
			os.Stdout.Write([]byte(ex.GetDefaultLogString()))
		}
	}()

	parseFlag()
	//启动内部日志服务
	logger.StartInnerLogHandler(basePath)

	//加载xml配置文件
	config := config.InitConfig(configFile)

	//启动日志服务
	logger.StartLogHandler(config.Log.FilePath)

	//启动Task Service
	task.StartTaskService()

	//启动计数器服务
	counter.StartCounter()

	//异步处理操作系统信号
	go listenSignal()

	//根据配置启动 http server, 阻塞
	httpserver.StartServer()

}

func parseFlag() {
	var runEnv string
	if runEnv = os.Getenv(server.RunEnv_Flag); runEnv == "" {
		runEnv = httpserver.RunEnv_Develop
	}

	configPath = basePath + "/conf/" + runEnv
	httpserver.RunEnv = runEnv
	httpserver.ConfigPath = configPath

	//从命令行参数读取配置路径
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	if configFile == "" {
		configFile = configPath + "/app.conf"
	}

}

func listenSignal() {
	c := make(chan os.Signal, 1)
	//syscall.SIGSTOP
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		innerLogger.Info("main::listenSignal [" + s.String() + "]")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP: //配置重载
			innerLogger.Info("main::listenSignal reload config begin...")
			//重新加载xml配置文件
			config.InitConfig(configFile)
			//重启启动Task集合
			task.ReStartTaskService()
			//使httpserver.StartServer返回nil
			httpserver.RestartServer()
			innerLogger.Info("main::listenSignal reload config end")

		default:
			return
		}
	}
}
