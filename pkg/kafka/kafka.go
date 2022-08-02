package kafka

import (
	"context"
	"errors"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	Writer *kafka.Writer
}

// New создает и инициализирует клиента Kafka.
func New(brokers []string, topic string, groupId string) (*Client, error) {
	if len(brokers) == 0 || brokers[0] == "" || topic == "" {
		return nil, errors.New("не указаны параметры подключения к Kafka")
	}

	c := Client{}

	c.Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers[0]),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &c, nil
}

func (c *Client) SendMessages(messages []kafka.Message) error {
	err := c.Writer.WriteMessages(context.Background(), messages...)
	return err
}
