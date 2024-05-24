package kafka

import (
	"context"
	"errors"
	"strings"
	"time"

	"log/slog"

	ckafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"gitlab.angkas.com/avengers/microservice/incentive-service/worker"
)

func NewKafkaProducer(cfg *WriterConfig) (*ckafka.Producer, error) {
	if cfg == nil {
		return nil, errors.New("kafka config is nil")
	}
	cm := ckafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
		"security.protocol": cfg.SecurityProtocol,
		"batch.size":        cfg.BatchSize,
		"linger.ms":         cfg.LingerMs,
		"compression.type":  cfg.CompressionType,
		"acks":              cfg.Acks,
	}

	if strings.ToUpper(cfg.SecurityProtocol) != "PLAINTEXT" {
		cm["sasl.mechanisms"] = cfg.SSLMechanism
		cm["sasl.username"] = cfg.Username
		cm["sasl.password"] = cfg.Password
	}

	return ckafka.NewProducer(&cm)
}

func NewClient(l *slog.Logger, cfg *WriterConfig, producer KafkaProducer, now func() time.Time) *Client {
	cm := ckafka.ConfigMap{
		"bootstrap.servers": cfg.Servers,
		"security.protocol": cfg.SecurityProtocol,
		"batch.size":        cfg.BatchSize,
		"linger.ms":         cfg.LingerMs,
		"compression.type":  cfg.CompressionType,
		"acks":              cfg.Acks,
	}

	if strings.ToUpper(cfg.SecurityProtocol) != "PLAINTEXT" {
		cm["sasl.mechanisms"] = cfg.SSLMechanism
		cm["sasl.username"] = cfg.Username
		cm["sasl.password"] = cfg.Password
	}

	return &Client{
		configmap: cm,
		logger:    l,
		producer:  producer,
		quit:      make(chan struct{}, 1),
		now:       now,
		topic:     cfg.Topic,
		timeout:   cfg.FlushTimeout,
	}
}

func (c *Client) Listen(topics []string, queue chan<- worker.Job) (stop func(), err error) {
	ckc := c.configmap
	ckc["enable.auto.commit"] = false
	ckc["group.id"] = "kafka-go-getting-started"
	ckc["auto.offset.reset"] = "latest"
	ckc["session.timeout.ms"] = 20000
	consumer, err := ckafka.NewConsumer(&ckc)
	if err != nil {
		return nil, err
	}

	if err = consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	stop = func() {
		if !consumer.IsClosed() {
			c.logger.Info("consumer already closed")
			return
		}
		if err = consumer.Close(); err != nil {
			c.logger.Error("consumer close", "err", err)
		}
		c.logger.Info("consumer closed")
	}

	go func(consumer *ckafka.Consumer) {
		for {
			select {
			case <-c.quit:
				c.logger.Info("consumer quit")
				return

			default:
				m, err := consumer.ReadMessage(time.Second)
				if err != nil {
					if err.(ckafka.Error).IsTimeout() {
						continue
					}
					c.logger.Error("read message", "err", err)
					continue
				}
				if m == nil || m.TopicPartition.Topic == nil {
					continue
				}

				queue <- worker.Job{
					Topic:   *m.TopicPartition.Topic,
					Payload: m.Value,
					Done: func() error {
						_, err = consumer.CommitMessage(m)
						return err
					},
				}
			}
		}
	}(consumer)

	return stop, nil
}

// PL-53: Added optional topic to parameter
func (p *Client) Produce(ctx context.Context, key []byte, value []byte, optionalTopic ...string) error {
	reportingChan := make(chan ckafka.Event)

	// Use the optional topic if provided, otherwise use the default topic
	topic := p.topic
	if len(optionalTopic) > 0 && optionalTopic[0] != "" {
		topic = optionalTopic[0]
	}

	msg := &ckafka.Message{
		Key:   key,
		Value: value,
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Timestamp: p.now(),
	}

	// Produce the message asynchronously
	if err := p.producer.Produce(msg, reportingChan); err != nil {
		p.logger.Error("Failed to produce message", "err", err)
		flushedMessages := p.producer.Flush(30 * 1000)
		p.logger.Info("Flushed kafka messages. Outstanding events still un-flushed", "flushed_messages", flushedMessages)
		return err
	}

	// RCA-14 & 27: Implement Error Handling to avoid Local Queue Full
	// Listen to all the events on the default events channel
	// Handle and Close the reporting events to avoid filling up the local queue
	// ----------------
	// Not Serving the kafka.Event would lead to memory leak
	go func(reportingChan chan ckafka.Event) {
		defer func() {
			if r := recover(); r != nil {
				p.logger.Error("goroutine panicked", "recover", r)
			}
			close(reportingChan)
		}()

		delivery := <-reportingChan
		// Handle the delivery report
		switch report := delivery.(type) {
		case *ckafka.Message:
			if report.TopicPartition.Error != nil {
				p.logger.Error("Failed to deliver message", "err", report.TopicPartition.Error)
			}
		case ckafka.Error:
			p.logger.Error("Failed to deliver message", "err", report)
		default:
			p.logger.Error("Ignored event:", "event", delivery)
		}

	}(reportingChan)
	// ----------------

	return nil
}

func (p *Client) Close() error {
	p.producer.Flush(p.timeout)
	p.producer.Close()
	p.logger.Info("flushed and closed producer")
	return nil
}
