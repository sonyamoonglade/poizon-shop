package redis

import (
	"context"
	"encoding/json"
	"fmt"
)

type Bus[PayloadType any] struct {
	client *Client
}

func NewBus[PayloadType any](c *Client) *Bus[PayloadType] {
	return &Bus[PayloadType]{
		client: c,
	}
}

type CallbackFunc[PayloadType any] func(payload PayloadType) error
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
			if msg.Payload == "" {
				runCallback[PayloadType](topicName, cb, *new(PayloadType), errorHandler)
				continue
			}

			var out PayloadType
			err := json.Unmarshal(stringToBytes(msg.Payload), &out)
			if err != nil {
				errorHandler(topicName, fmt.Errorf("json.Unmarshal: %w", err))
				continue
			}

			runCallback[PayloadType](topicName, cb, out, errorHandler)
		}
	}
}

func (b *Bus[PayloadType]) SendToTopic(ctx context.Context, topicName string, payload any) error {
	if payload == nil {
		fmt.Println("sending empty")
		return b.client.c.Publish(ctx, topicName, nil).Err()
	}

	packed, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	return b.client.c.Publish(ctx, topicName, packed).Err()
}

func runCallback[T any](topicName string, cb CallbackFunc[T], payload T, errHandler ErrorHandlerFunc) {
	if err := cb(payload); err != nil {
		errHandler(topicName, fmt.Errorf("callback: %w", err))
	}
}
