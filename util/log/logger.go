// logger
package logger

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

type ChanLog struct {
	Content   string
	LogTarget string
}

var (
	logChan_Debug  chan ChanLog
	logChan_Info   chan ChanLog
	logChan_Warn   chan ChanLog
	logChan_Error  chan ChanLog
	logChan_Custom chan ChanLog
)

var (
	logRootPath string
)

func init() {
	logChan_Debug = make(chan ChanLog, 100000)
	logChan_Info = make(chan ChanLog, 100000)
	logChan_Warn = make(chan ChanLog, 100000)
	logChan_Error = make(chan ChanLog, 100000)
	logChan_Custom = make(chan ChanLog, 100000)
}

func Debug(log string, logTarget string) {
	chanLog := ChanLog{
		LogTarget: logTarget,
		Content:   log,
	}
	logChan_Debug <- chanLog
}

func Info(log string, logTarget string) {
	chanLog := ChanLog{
		LogTarget: logTarget,
		Content:   log,
	}
	logChan_Info <- chanLog
}

func Warn(log string, logTarget string) {
	chanLog := ChanLog{
		LogTarget: logTarget,
		Content:   log,
	}
	logChan_Warn <- chanLog
}

func Error(log string, logTarget string) {
	chanLog := ChanLog{
		LogTarget: logTarget,
		Content:   log,
	}
	logChan_Error <- chanLog
}

func Log(log string, logTarget string, logLevel string) {
	chanLog := ChanLog{
		LogTarget: logTarget + "_" + logLevel,
		Content:   log,
	}
	logChan_Custom <- chanLog
}

//开启日志处理器
func StartLogHandler(rootPath string) {
	//设置日志根目录
	logRootPath = rootPath
	if !strings.HasSuffix(logRootPath, "/") {
		logRootPath = logRootPath + "/"
	}

	go handleDebug()
	go handleInfo()
	go handleWarn()
	go handleError()
	go handleCustom()
}

//处理日志内部函数
func handleDebug() {
	for {
		log := <-logChan_Debug
		writeLog(log, "debug")
	}
}
func handleInfo() {
	for {
		log := <-logChan_Info
		writeLog(log, "info")
	}
}
func handleWarn() {
	for {
		log := <-logChan_Warn
		writeLog(log, "warn")
	}
}
func handleError() {
	for {
		log := <-logChan_Error
		writeLog(log, "error")
	}
}

func handleCustom() {
	for {
		log := <-logChan_Custom
		writeLog(log, "custom")
	}
}

func writeLog(chanLog ChanLog, level string) {
	filePath := logRootPath + chanLog.LogTarget
	switch level {
	case "debug":
		filePath = filePath + "_debug_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "info":
		filePath = filePath + "_info_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "warn":
		filePath = filePath + "_warn_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "error":
		filePath = filePath + "_error_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
	case "custom":
		filePath = filePath + "_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
		break
	}
	log := time.Now().Format(defaultFullTimeLayout) + " " + chanLog.Content
	writeFile(filePath, log)
}

func writeFile(logFile string, log string) {
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	logstr := log + "\r\n"
	file, err := os.OpenFile(logFile, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(logFile, err)
		return
	}
	file.WriteString(logstr)
}
