package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	channel *amqp.Channel
}

// NewPublisher initializes a new Publisher with the given amqp.Channel.
func NewPublisher(channel *amqp.Channel) *Publisher {
	return &Publisher{
		channel: channel,
	}
}

// PublishOption defines a functional option for Publish.
type PublishOption func(*publishOptions)

// publishOptions holds optional parameters for Publish.
type publishOptions struct {
	routingKey  string
	mandatory   bool
	immediate   bool
	contentType string
	publishing  amqp.Publishing
}

// WithRoutingKey sets the routing key.
func WithRoutingKey(key string) PublishOption {
	return func(opts *publishOptions) {
		opts.routingKey = key
	}
}

// WithMandatory sets the mandatory flag.
func WithMandatory(mandatory bool) PublishOption {
	return func(opts *publishOptions) {
		opts.mandatory = mandatory
	}
}

// WithImmediate sets the immediate flag.
func WithImmediate(immediate bool) PublishOption {
	return func(opts *publishOptions) {
		opts.immediate = immediate
	}
}

// WithContentType sets the content type.
func WithContentType(contentType string) PublishOption {
	return func(opts *publishOptions) {
		opts.contentType = contentType
	}
}

// Publish sends a message to the specified subject (exchange) with options.
func (p *Publisher) Publish(ctx context.Context, subject string, message []byte, opts ...PublishOption) error {
	// Initialize default publish options
	pOpts := &publishOptions{
		routingKey:  "",
		mandatory:   false,
		immediate:   false,
		contentType: "text/plain",
		publishing: amqp.Publishing{
			Body: message,
		},
	}

	// Apply provided options
	for _, opt := range opts {
		opt(pOpts)
	}

	// Set content type if provided
	if pOpts.contentType != "" {
		pOpts.publishing.ContentType = pOpts.contentType
	}

	err := p.channel.Publish(
		subject,          // exchange
		pOpts.routingKey, // routing key
		pOpts.mandatory,  // mandatory
		pOpts.immediate,  // immediate
		pOpts.publishing, // publishing
	)
	if err != nil {
		return err
	}
	return nil
}
