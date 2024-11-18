package rabbitmq

import (
	"context"
	"crypto/tls"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	tlsCfg  *tls.Config
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func New(url string, tlsCfg *tls.Config) *RabbitMQ {
	return &RabbitMQ{
		tlsCfg: tlsCfg,
		url:    url,
	}
}

func (r *RabbitMQ) Connect(ctx context.Context) error {
	var err error
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			r.conn, err = amqp.DialTLS(r.url, r.tlsCfg)
			if err == nil {
				r.channel, err = r.conn.Channel()
				if err == nil {
					return nil
				}
			}
			log.Printf("Failed to connect to RabbitMQ, retrying in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitMQ) Reconnect(ctx context.Context) error {
	r.Close()
	return r.Connect(ctx)
}

type ExchangeType string

const (
	ExchangeDirect  ExchangeType = "direct"
	ExchangeFanout  ExchangeType = "fanout"
	ExchangeTopic   ExchangeType = "topic"
	ExchangeHeaders ExchangeType = "headers"
)

// ExchangeOptions defines options for declaring an exchange.
type ExchangeOptions struct {
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// ExchangeOptionFn defines a functional option for DeclareExchange.
type ExchangeOptionFn func(*ExchangeOptions)

// WithExchangeDurable sets the Durable option.
func WithExchangeDurable(durable bool) ExchangeOptionFn {
	return func(opts *ExchangeOptions) {
		opts.Durable = durable
	}
}

// WithExchangeAutoDelete sets the AutoDelete option.
func WithExchangeAutoDelete(autoDelete bool) ExchangeOptionFn {
	return func(opts *ExchangeOptions) {
		opts.AutoDelete = autoDelete
	}
}

// WithExchangeInternal sets the Internal option.
func WithExchangeInternal(internal bool) ExchangeOptionFn {
	return func(opts *ExchangeOptions) {
		opts.Internal = internal
	}
}

// WithExchangeNoWait sets the NoWait option.
func WithExchangeNoWait(noWait bool) ExchangeOptionFn {
	return func(opts *ExchangeOptions) {
		opts.NoWait = noWait
	}
}

// WithExchangeArgs sets the Args option.
func WithExchangeArgs(args amqp.Table) ExchangeOptionFn {
	return func(opts *ExchangeOptions) {
		opts.Args = args
	}
}

// NewExchangeOptions initializes ExchangeOptions with default values and applies any provided options.
func NewExchangeOptions(options ...ExchangeOptionFn) *ExchangeOptions {
	opts := &ExchangeOptions{
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

// DeclareExchange declares an exchange with the provided options.
func (r *RabbitMQ) DeclareExchange(exchangeName string, exchangeType ExchangeType, options ...ExchangeOptionFn) error {
	opts := NewExchangeOptions(options...)

	return r.channel.ExchangeDeclare(
		exchangeName,         // name of the exchange
		string(exchangeType), // type
		opts.Durable,         // durable
		opts.AutoDelete,      // auto-deleted
		opts.Internal,        // internal
		opts.NoWait,          // no-wait
		opts.Args,            // arguments
	)
}

// QueueOptions defines options for declaring a queue.
type QueueOptions struct {
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

// QueueOption defines a functional option for DeclareQueue.
type QueueOption func(*QueueOptions)

// WithQueueDurable sets the Durable option.
func WithQueueDurable(durable bool) QueueOption {
	return func(opts *QueueOptions) {
		opts.Durable = durable
	}
}

// WithQueueAutoDelete sets the AutoDelete option.
func WithQueueAutoDelete(autoDelete bool) QueueOption {
	return func(opts *QueueOptions) {
		opts.AutoDelete = autoDelete
	}
}

// WithQueueExclusive sets the Exclusive option.
func WithQueueExclusive(exclusive bool) QueueOption {
	return func(opts *QueueOptions) {
		opts.Exclusive = exclusive
	}
}

// WithQueueNoWait sets the NoWait option.
func WithQueueNoWait(noWait bool) QueueOption {
	return func(opts *QueueOptions) {
		opts.NoWait = noWait
	}
}

// WithQueueArgs sets the Args option.
func WithQueueArgs(args amqp.Table) QueueOption {
	return func(opts *QueueOptions) {
		opts.Args = args
	}
}

// NewQueue initializes QueueOptions with default values and applies any provided options.
func NewQueue(options ...QueueOption) *QueueOptions {
	opts := &QueueOptions{
		Durable:    true,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false,
		Args:       nil,
	}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

// DeclareQueue declares a queue with the provided options.
func (r *RabbitMQ) DeclareQueue(queueName string, options ...QueueOption) (amqp.Queue, error) {
	opts := NewQueue(options...)

	return r.channel.QueueDeclare(
		queueName,       // name of the queue
		opts.Durable,    // durable
		opts.AutoDelete, // delete when unused
		opts.Exclusive,  // exclusive
		opts.NoWait,     // no-wait
		opts.Args,       // arguments
	)
}

// BindOptions defines options for binding a queue.
type BindOptions struct {
	Args amqp.Table
}

// BindOption defines a functional option for BindQueue.
type BindOption func(*BindOptions)

// WithBindArgs sets the Args option.
func WithBindArgs(args amqp.Table) BindOption {
	return func(opts *BindOptions) {
		opts.Args = args
	}
}

// NewBind initializes BindOptions with default values and applies any provided options.
func NewBind(options ...BindOption) *BindOptions {
	opts := &BindOptions{
		Args: nil,
	}
	for _, opt := range options {
		opt(opts)
	}
	return opts
}

// BindQueue binds a queue with the provided options.
func (r *RabbitMQ) BindQueue(queueName, exchange, routingKey string, options ...BindOption) error {
	opts := NewBind(options...)

	return r.channel.QueueBind(
		queueName,  // name of the queue
		routingKey, // binding key
		exchange,   // source exchange
		false,      // no-wait
		opts.Args,  // arguments
	)
}
