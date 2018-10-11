package pusher

import (
	"TechPlat/datapipe/component/kafka"
	"fmt"
	"strconv"
	"testing"
)

func TestKafkaSendMessage(t *testing.T) {

	kafkaServerUrl := "192.168.240.28:9092,192.168.240.29:9092,192.168.240.26:9092"
	val := "11111"

	for i := 0; i < 3; i++ {
		partition, offset, kafkaErr := kafka.SendMessage(kafkaServerUrl, "test", "11111")
		if kafkaErr != nil {
			fmt.Println("InsertKafkaData[" + val + "] error -> " + kafkaErr.Error())
		} else {
			fmt.Println("InsertKafkaData success -> [" + val + "] [" +
				strconv.Itoa(int(partition)) + "," + strconv.FormatInt(offset, 10) +
				"]")
		}
	}
}
