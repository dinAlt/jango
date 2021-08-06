package jango

import (
	"fmt"
)

func wrapErr(where string, err error) error {
	return fmt.Errorf("[jango.%s] %w", where, err)
}

type TransportError struct {
	Inner error
}

func (e *TransportError) Error() string {
	return fmt.Sprintf("transport error: %s", e.Inner.Error())
}

func (e *TransportError) Unwrap() error {
	return e.Inner
}

type JanusError struct {
	Code   int    `json:"code,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (e *JanusError) Error() string {
	return fmt.Sprintf("janus: [%d] %s", e.Code, e.Reason)
}
