package msgproducer

import (
	"github.com/segmentio/kafka-go"

	"github.com/Pickausernaame/chat-service/internal/logger"
)

const serviceName = "msg-producer"

func NewKafkaWriter(brokers []string, topic string, batchSize int) KafkaWriter {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.Murmur2Balancer{Consistent: true},
		BatchSize:    batchSize,
		RequiredAcks: kafka.RequireOne,
		Async:        false,
		Logger:       logger.NewKafkaAdapted().WithServiceName(serviceName),
		ErrorLogger:  logger.NewKafkaAdapted().WithServiceName(serviceName).ForErrors(),
	}
}
