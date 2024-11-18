package rabbitmq

import (
	"context"

	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrNilChannel = errors.New("channel is nil")

// consumeOptions defines options for consuming messages.
type consumeOptions struct {
	Queue     string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
}

// ConsumeOption defines a functional option for Consume.
type ConsumeOption func(*consumeOptions)

// WithConsumer sets the Consumer option.
func WithQueueName(queue string) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.Queue = queue
	}
}

// WithAutoAck sets the AutoAck option.
func WithAutoAck(autoAck bool) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.AutoAck = autoAck
	}
}

// WithExclusive sets the Exclusive option.
func WithExclusive(exclusive bool) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.Exclusive = exclusive
	}
}

// WithNoLocal sets the NoLocal option.
func WithNoLocal(noLocal bool) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.NoLocal = noLocal
	}
}

// WithNoWait sets the NoWait option.
func WithNoWait(noWait bool) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.NoWait = noWait
	}
}

// WithArgs sets the Args option.
func WithArgs(args amqp.Table) ConsumeOption {
	return func(opts *consumeOptions) {
		opts.Args = args
	}
}

// Consumer consumes messages from a queue.
type Consumer struct {
	channel *amqp.Channel
	opts    *consumeOptions
}

// NewConsumer initializes a Consumer with the provided options.
func NewConsumer(channel *amqp.Channel, opts ...ConsumeOption) (*Consumer, error) {
	if channel == nil {
		return nil, ErrNilChannel
	}

	options := &consumeOptions{
		Queue:     "",
		AutoAck:   true,
		Exclusive: false,
		NoLocal:   false,
		NoWait:    false,
		Args:      nil,
	}
	for _, opt := range opts {
		opt(options)
	}

	return &Consumer{
		channel: channel,
		opts:    options,
	}, nil
}

// Consume starts consuming messages with the provided handler.
func (c *Consumer) Consume(ctx context.Context, handler func(msg []byte) error) error {
	msgs, err := c.channel.Consume(
		"queue_name",     // queue
		c.opts.Queue,     // consumer
		c.opts.AutoAck,   // auto-ack
		c.opts.Exclusive, // exclusive
		c.opts.NoLocal,   // no-local
		c.opts.NoWait,    // no-wait
		c.opts.Args,      // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgs:
			if err := handler(msg.Body); err != nil {
				return err
			}
		}
	}
}
