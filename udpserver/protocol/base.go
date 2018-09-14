package protocol

type LogInfo struct {
	Message  string
	Level    string
	Appid    string
	Clientip string
}

type ProtocolHandler func (buf [4096]byte, n int) LogInfo