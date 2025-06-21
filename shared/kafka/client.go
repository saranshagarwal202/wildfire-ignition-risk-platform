package kafka

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	brokers []string
}

func NewClient(brokers string) *Client {
	return &Client{
		brokers: strings.Split(brokers, ","),
	}
}

func (c *Client) NewWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(c.brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func (c *Client) NewReader(topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func (c *Client) CreateTopics(ctx context.Context, topics []string) error {
	conn, err := kafka.DialContext(ctx, "tcp", c.brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfigs := make([]kafka.TopicConfig, len(topics))
	for i, topic := range topics {
		topicConfigs[i] = kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		}
	}

	return conn.CreateTopics(topicConfigs...)
}
