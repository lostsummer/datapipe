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
	configFile  string
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			ex := exception.CatchError("DataPipe main->recover error", err)
			innerLogger.Error(ex.GetDefaultLogString())
			os.Stdout.Write([]byte(ex.GetDefaultLogString()))
		}
	}()

	currentBaseDir := common.GetCurrentDirectory()
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	if configFile == "" {
		configFile = currentBaseDir + "/app.conf"
	}
	//启动内部日志服务
	logger.StartInnerLogHandler(currentBaseDir)

	//加载xml配置文件
	config := config.InitConfig(configFile)

	//启动日志服务
	logger.StartLogHandler(config.Log.FilePath)

	//启动Task Service
	task.StartTaskService()

	//启动计数器服务
	counter.StartCounter()

	//异步处理操作系统信号
	go waitSignal()

	srvStatus := make(chan string)
	runServer := func() {
		err := httpserver.StartServer()
		if err != nil {
			innerLogger.Error("HttpServer.StartServer error: " + err.Error())
			srvStatus <- "end"
		} else {
			srvStatus <- "restart"
		}
	}

	//开启httpserver
	go runServer()

	//httpserver简单控制
	for s := range srvStatus {
		switch s {
		case "end":
			return
		case "restart":
			go runServer()
		default:
			return
		}
	}
}

func waitSignal() {
	c := make(chan os.Signal, 1)
	//syscall.SIGSTOP
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		innerLogger.Info("main::waitSignal [" + s.String() + "]")
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			return
		case syscall.SIGHUP: //配置重载
			innerLogger.Info("main::waitSignal reload config begin...")
			//重新加载xml配置文件
			config.InitConfig(configFile)
			//重启启动Task集合
			task.ReStartTaskService()
			//使httpserver.StartServer返回nil
			httpserver.RestartServer()
			innerLogger.Info("main::waitSignal reload config end")

		default:
			return
		}
	}
}
