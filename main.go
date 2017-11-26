/*emoney.tongjiservice
* Author: Panxinming
* LastUpdateTime: 2016-10-17 10:00
 */
package main

import (
	"emoney/tongjiservice/config"
	"emoney/tongjiservice/util/common"
	"emoney/tongjiservice/util/log"
	"emoney/tongjiservice/task"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"emoney/tongjiservice/core/exception"
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
			ex := exception.CatchError("TongjiService main->recover error", err)
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

	//阻塞主线程，同时等待操作系统信号
	waitSignal()
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
			innerLogger.Info("main::waitSignal reload config end")
		default:
			return
		}
	}
}
