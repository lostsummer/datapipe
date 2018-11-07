package endpoint

import (
	"TechPlat/datapipe/component/kafka"
)

type Kafka struct {
	URL   string
	Topic string
}

func (k *Kafka) Push(val string) (int64, error) {
	//partition, offset, kafkaErr := kafka.SendMessage(k.URL, k.Topic, val)
	_, _, kafkaErr := kafka.SendMessage(k.URL, k.Topic, val)
	if kafkaErr != nil {
		return -1, kafkaErr
	}
	return 1, nil
}
