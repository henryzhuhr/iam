// Package queue provides Kafka producer wrapper for the IAM application.
package queue

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"
)

// KafkaConfig holds Kafka connection configuration.
type KafkaConfig struct {
	Brokers []string `json:"Brokers"`
	Topic   string   `json:"Topic,optional"`
}

// KafkaProducer wraps Kafka producer for audit log publishing.
type KafkaProducer struct {
	config  KafkaConfig
	stopped bool
}

// NewKafkaProducer creates a Kafka producer (stub implementation).
func NewKafkaProducer(cfg KafkaConfig) (*KafkaProducer, error) {
	logx.Infof("kafka producer initialized (stub): brokers=%v", cfg.Brokers)
	return &KafkaProducer{config: cfg}, nil
}

// SendMessage publishes a message to Kafka (stub implementation).
func (p *KafkaProducer) SendMessage(ctx context.Context, topic string, key, value []byte) error {
	if p.stopped {
		return fmt.Errorf("producer is stopped")
	}
	logx.Infof("kafka send message (stub): topic=%s key=%s", topic, string(key))
	return nil
}

// Close releases producer resources.
func (p *KafkaProducer) Close() {
	p.stopped = true
	logx.Info("kafka producer closed")
}
