package kafka

import (
	"context"
	"encoding/json"
	"patches/protocol"
	"strconv"
	"time"

	segkafka "github.com/segmentio/kafka-go"
)

// Publisher defines the publishing methods for a messaging system.
type Publisher interface {
	PublishPatch(msg []byte) error
}

// Writer represents an entity for writing to one or more Kafka topics.
type Writer struct {
	patchesWriter *segkafka.Writer
}

// NewWriter initializes a new Writer.
func NewWriter(location, topic string) *Writer {
	return &Writer{
		patchesWriter: segkafka.NewWriter(segkafka.WriterConfig{
			Brokers:      []string{location},
			Topic:        topic,
			Balancer:     &segkafka.LeastBytes{},
			BatchTimeout: 10 * time.Millisecond,
		}),
	}
}

// PublishPatch publishes a message to a Kafka topic using a conversationID as
// the key and a struct representation of a WebSocket Update message (msg) as
// the value.
func (k *Writer) PublishPatch(msg protocol.Message, conversationID int64) error {
	pubBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = k.patchesWriter.WriteMessages(context.Background(),
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
