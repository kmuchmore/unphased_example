package mq

import "context"

type Consumer interface {
	Consume(ctx context.Context, handler func(msg []byte) error) error
}
