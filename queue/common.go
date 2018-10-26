package queue

// 适配各种类型target推送逻辑
type Target interface {
	Push(val string) (int64, error)
}

type Source interface {
	Pop() (string, error)
}
