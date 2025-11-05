// path: pkg/kafka/producer.go
package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/proto"

	pricev1 "github.com/binari-digital/price-engine/proto/price/v1"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			BatchTimeout: 5 * time.Millisecond,
		},
	}
}

func (p *Producer) Close() error { return p.writer.Close() }

func (p *Producer) Publish(ctx context.Context, c *pricev1.Candle) error {
	b, err := proto.Marshal(c)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(c.Symbol),
		Value: b,
		Time:  time.UnixMilli(c.GetLastTradeTs()),
	}
	return p.writer.WriteMessages(ctx, msg)
}
