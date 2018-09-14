package protocol

import (
	"encoding/json"
)

type logMessage struct {
	Appid   string
	Level   string
	Message string
}

func OnHandleJsonProtocol(buf [4096]byte, n int) LogInfo {
	logStr := string(buf[0:n])

	var msg logMessage
	err := json.Unmarshal(buf[0:n], &msg)
	if err != nil {
		return LogInfo{logStr, "UNKNOW", "", err.Error()}
	} else {
		return LogInfo{msg.Message, msg.Level, msg.Appid, ""}
	}
}