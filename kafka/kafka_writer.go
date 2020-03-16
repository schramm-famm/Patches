package kafka

import (
	"context"
	"encoding/json"
	"patches/protocol"
	"strconv"
	"time"

	segkafka "github.com/segmentio/kafka-go"
)

type Publisher interface {
	PublishPatch(msg []byte) error
}

type Writer struct {
	*segkafka.Writer
}

func NewWriter(location, topic string) *Writer {
	return &Writer{
		segkafka.NewWriter(segkafka.WriterConfig{
			Brokers:      []string{location},
			Topic:        topic,
			Balancer:     &segkafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
		}),
	}
}

func (k *Writer) PublishPatch(msg protocol.Message, conversationID int64) error {
	pubBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = k.WriteMessages(context.Background(),
		segkafka.Message{
			Key:   []byte(strconv.FormatInt(conversationID, 10)),
			Value: pubBytes,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
