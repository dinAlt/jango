package jango

import (
	"context"

	"github.com/dinalt/jango/logging"
)

type Logger = logging.Logger

type Transport interface {
	Request(ctx context.Context, req interface{}, resp interface{}) error
}
