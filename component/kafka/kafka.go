package kafka

import(
	"github.com/Shopify/sarama"
	"strings"
	"github.com/pkg/errors"
	"fmt"
	"sync"
)

var(
	producerConfig *sarama.Config
	producerMap sync.Map

)

func init(){
	producerConfig := sarama.NewConfig()
	producerConfig.ClientID = "datapipe"
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	producerConfig.Producer.Return.Successes = true
}

//GetKafkaProducer get producer by serverUrls
func GetKafkaProducer(kafkaServerUrls string)  (producer sarama.SyncProducer, err error){
	producerObj, isOk := producerMap.Load(kafkaServerUrls)
	if !isOk{
		producer, err = sarama.NewSyncProducer(strings.Split(kafkaServerUrls, ","), producerConfig)
		if err != nil {
			return nil, err
		}else{
			producerMap.Store(kafkaServerUrls, producer)
		}
	}else{
		producer = producerObj.(sarama.SyncProducer)
	}
	return producer, nil
}


// SendMessage send topic and value to kafka server
func SendMessage(kafkaServerUrls string, topic, value string) (partition int32, offset int64, err error){
	producer, err := GetKafkaProducer(kafkaServerUrls)
	if err != nil {
		return 0,0,errors.New(fmt.Sprintf("Failed to produce message: %s", err))
	}
	//defer producer.Close()
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Key = sarama.StringEncoder("default")
	msg.Value = sarama.ByteEncoder(value)
	partition, offset, err = producer.SendMessage(msg)
	if err != nil {
		return 0,0,errors.New(fmt.Sprintf("Failed to produce message: %s", err))
	}
	return partition, offset, nil
}
