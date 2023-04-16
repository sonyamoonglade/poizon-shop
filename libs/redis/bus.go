package redis

import (
	"context"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type Bus[PayloadType any] struct {
	client *Client
}

func NewBus[PayloadType any](c *Client) *Bus[PayloadType] {
	return &Bus[PayloadType]{
		client: c,
	}
}

type CallbackFunc[PayloadType any] func(payload PayloadType)
type ErrorHandlerFunc func(topicName string, err error)

func (b *Bus[PayloadType]) SubscribeToTopicWithCallback(
	ctx context.Context,
	topicName string,
	cb CallbackFunc[PayloadType],
	errorHandler ErrorHandlerFunc,
) {
	topic := b.client.c.Subscribe(ctx, topicName)
	defer topic.Close()
	updates := topic.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-updates:
			var out PayloadType
			err := msgpack.Unmarshal(stringToBytes(msg.Payload), &out)
			if err != nil {
				errorHandler(topicName, fmt.Errorf("msgpack.Unmarshal: %w", err))
				continue
			}
			cb(out)
		}
	}
}

func (b *Bus[PayloadType]) SendToTopic(ctx context.Context, topicName string, payload PayloadType) error {
	packed, err := msgpack.Marshal(payload)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal: %w", err)
	}
	return b.client.c.Publish(ctx, topicName, packed).Err()
}
