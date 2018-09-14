package protocol

import (
	"strings"
)

func OnHandleTextProtocol(buf [4096]byte, n int) LogInfo {
	logStr := string(buf[0:n])
	appid := ""

	appidInfo := Substring(logStr, 0, 7)
	if len(appidInfo) == 7 {
		appid = strings.Trim(strings.Trim(appidInfo, "["), "]")
	}

	return LogInfo{logStr, "UNKNOWN", appid, ""}
}

func Substring(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}
