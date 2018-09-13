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

func OutputFileAdapter (conf config.OutputAdapter, appid string, logstr string) {
	logFile := getLogFile(appid, conf.Url)
	logstr = logstr + "\r\n"
	writeFile(logFile, logstr)
}

func getLogFile(appid string, path string) string {
	if path == "" {
		path = defaultFileRoot
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if appid != "" {
		path = path + appid + "/"
	}
	os.MkdirAll(path, os.ModePerm)

	logFile := path + "log_" + time.Now().Format(defaultTimeLayout) + ".log"
	return logFile
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