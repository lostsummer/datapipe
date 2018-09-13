package outputadapter

import (
	"time"
	"os"
	"syscall"
	"fmt"
	"strings"
	"TechPlat/datapipe/config"
)

const (
	defaultFileRoot   = "/emoney/logserver/logs/"
	defaultTimeLayout = "2006_01_02"
)

var currentFileRoot = defaultFileRoot

func OutputFileAdapter (conf config.OutputAdapter, appid string, logstr string) {
	if conf.Url != "" {
		if strings.HasSuffix(conf.Url, "/") {
			currentFileRoot = conf.Url
		} else {
			currentFileRoot = conf.Url + "/"
		}
	}

	if appid != "" {
		writeLogWithApp(appid, logstr)
	} else {
		writeLog(logstr)
	}
}

func writeLog(log string) {
	os.MkdirAll(currentFileRoot, os.ModePerm)
	logFile := currentFileRoot + "log_" + time.Now().Format(defaultTimeLayout) + ".log"
	logstr := log + "\r\n"

	writeFile(logFile, logstr)
}

func writeLogWithApp(appid string, log string) {
	logpath := currentFileRoot + appid
	os.MkdirAll(logpath, os.ModePerm)
	logFile := logpath + "/log_" + time.Now().Format(defaultTimeLayout) + ".log"
	logstr := log + "\r\n"

	writeFile(logFile, logstr)
}

func writeFile(logFile string, logstr string) {
	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666

	file, err := os.OpenFile(logFile, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(logFile, err, err.Error())
		return
	}
	file.WriteString(logstr)
}