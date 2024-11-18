package mq

import "context"

type Publisher interface {
	Publish(ctx context.Context, subject string, message []byte) error
}
