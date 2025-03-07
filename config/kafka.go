package config

import (
	"github.com/segmentio/kafka-go"
)

// InitKafka initializes the Kafka writer.
func InitKafka(conf Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(conf.Kafka.Brokers...),
		Topic:    conf.Kafka.Topic,
		Balancer: &kafka.LeastBytes{},
	}
}
