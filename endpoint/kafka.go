package endpoint

import (
	"TechPlat/datapipe/component/kafka"
)

type KafkaTarget struct {
	URL   string
	Topic string
}

func (k *KafkaTarget) Push(val string) (int64, error) {
	//partition, offset, kafkaErr := kafka.SendMessage(k.URL, k.Topic, val)
	_, _, kafkaErr := kafka.SendMessage(k.URL, k.Topic, val)
	if kafkaErr != nil {
		return -1, kafkaErr
	}
	return 1, nil
}
