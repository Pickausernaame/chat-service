package logger

import (
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _ kafka.Logger = (*KafkaAdapted)(nil)

type KafkaAdapted struct {
	lg  *zap.Logger
	lvl zapcore.Level
}

func (k KafkaAdapted) Printf(s string, i ...interface{}) {
	switch k.lvl {
	case zapcore.ErrorLevel:
		k.lg.Error(fmt.Sprintf(s, i...))
	case zapcore.InfoLevel:
		k.lg.Info(fmt.Sprintf(s, i...))
	}
}

func NewKafkaAdapted() *KafkaAdapted {
	return &KafkaAdapted{
		lg:  zap.L(),
		lvl: zapcore.InfoLevel,
	}
}

func (k *KafkaAdapted) ForErrors() *KafkaAdapted {
	k.lvl = zapcore.ErrorLevel
	return k
}

func (k *KafkaAdapted) WithServiceName(name string) *KafkaAdapted {
	k.lg = zap.L().Named(name)
	return k
}
