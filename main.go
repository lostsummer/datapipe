/*emoney.tongjiservice
* Author: Panxinming
* LastUpdateTime: 2016-10-17 10:00
 */
package main

import (
	"TechPlat/datapipe/config"
	"TechPlat/datapipe/core/exception"
	"TechPlat/datapipe/counter"
	"TechPlat/datapipe/global"
	"TechPlat/datapipe/httpserver"
	"TechPlat/datapipe/task"
	"TechPlat/datapipe/util/common"
	"TechPlat/datapipe/util/log"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
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
	if runEnv = os.Getenv(httpserver.RunEnv_Flag); runEnv == "" {
		runEnv = httpserver.RunEnv_Develop
	}

	httpserver.RunEnv = runEnv
	configPath = basePath + "/conf/" + runEnv

	//从命令行参数读取配置路径
	var version bool
	flag.BoolVar(&version, "v", false, "-v")
	flag.StringVar(&configFile, "config", "", "配置文件路径")
	flag.Parse()
	if version {
		fmt.Printf("Version: %s, Branch: %s, Build: %s, Build time: %s\n",
			global.Version, global.Branch, global.CommitID, global.BuildTime)
		return
	}
	if configFile == "" {
		configFile = configPath + "/app.conf"
	} else {
		if dir, err := filepath.Abs(filepath.Dir(configFile)); err == nil {
			configPath = dir
		}
	}
	httpserver.ConfigPath = configPath

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
