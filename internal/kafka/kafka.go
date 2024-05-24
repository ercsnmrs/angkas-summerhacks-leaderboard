package kafka

import (
	"context"
	"time"

	"log/slog"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaProducer is an interface of Confluent Kafka producer.
type KafkaProducer interface {
	Produce(msg *ckafka.Message, deliveryChan chan ckafka.Event) error
	Flush(timeoutMs int) int
	Close()
}

// Writer is interface of Kafka producer.
type Writer interface {
	Produce(ctx context.Context, key []byte, value []byte, optionalTopic ...string) error
	Close() error
}

type (
	Client struct {
		logger    *slog.Logger
		configmap ckafka.ConfigMap
		producer  KafkaProducer
		now       func() time.Time
		topic     string
		quit      chan struct{}
		timeout   int
	}
	WriterConfig struct {
		Servers          string
		SecurityProtocol string
		SSLMechanism     string
		Username         string
		Password         string
		Topic            string
		FlushTimeout     int
		BatchSize        int
		LingerMs         int
		CompressionType  string
		Acks             string
		BufferMemory     int
		ProducerTimeout  int
	}
)
