package kafka

import(
	"github.com/Shopify/sarama"
	"strings"
	"github.com/pkg/errors"
	"fmt"
)


// SendMessage send topic and value to kafka server
func SendMessage(kafkaServerUrls string, topic, value string) (partition int32, offset int64, err error){
	config := sarama.NewConfig()
	config.ClientID = "datapipe"
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(kafkaServerUrls, ","), config)
	if err != nil {
		return 0,0,errors.New(fmt.Sprintf("Failed to produce message: %s", err))
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	//msg.Partition = 1 //仅当设定为ManualPartitioner时，该设置才生效
	msg.Key = sarama.StringEncoder("default")
	msg.Value = sarama.ByteEncoder(value)
	partition, offset, err = producer.SendMessage(msg)
	if err != nil {
		return 0,0,errors.New(fmt.Sprintf("Failed to produce message: %s", err))
	}
	return partition, offset, nil
}
