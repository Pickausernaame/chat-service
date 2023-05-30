package afcverdictsprocessor

import (
	"context"
	"io"
	"strconv"

	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/dlq_writer_mock.gen.go -package=afcverdictsprocessormocks

type KafkaDLQWriter interface {
	io.Closer
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

func NewKafkaDLQWriter(brokers []string, topic string) KafkaDLQWriter {
	return &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}
}

func prepareDLQMessage(msg kafka.Message, err error) kafka.Message {
	return kafka.Message{
		Headers: []kafka.Header{
			{
				Key:   "LAST_ERROR",
				Value: []byte(err.Error()),
			},
			{
				Key:   "ORIGINAL_PARTITION",
				Value: []byte(strconv.Itoa(msg.Partition)),
			},
		},
		Value: msg.Value,
	}
}
