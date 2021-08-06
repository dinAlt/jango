package jango

import (
	"context"
)

type Transport interface {
	Request(ctx context.Context, req interface{}, resp interface{}) error
}
