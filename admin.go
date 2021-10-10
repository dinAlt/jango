package jango

import (
	"context"
	"fmt"
)

const transactionIDLen = 16

type Request interface {
	Build() interface{}
	Async() bool
}

type PluginRequest interface {
	Plugin() string
	Request
}

// Errer may be implemented by responses to signalize about plugin error
type Errer interface {
	Err() error
}

type Admin struct {
	AdminSecret string
	Transport
}

type request struct {
	Janus         string      `json:"janus,omitempty"`
	TransactionID string      `json:"transaction,omitempty"`
	AdminSecret   string      `json:"admin_secret,omitempty"`
	Plugin        string      `json:"plugin,omitempty"`
	Request       interface{} `json:"request,omitempty"`
}

func (r *request) Transaction() string {
	return r.TransactionID
}

type pluginResponse struct {
	Janus       string      `json:"janus,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	Error       *JanusError `json:"error,omitempty"`
	Response    interface{} `json:"response,omitempty"`
}

type sessionResponse struct {
	Janus       string      `json:"janus,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	Error       *JanusError `json:"error,omitempty"`
	Sessions    []int64     `json:"sessions,omitempty"`
}

type pingResponse struct {
	Janus       string      `json:"janus,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	Error       *JanusError `json:"error,omitempty"`
}

func (c *Admin) PluginRequestCtx(ctx context.Context, req PluginRequest,
	resp interface{}) error {
	wrap := func(err error) error {
		return wrapErr("Admin.PluginRequest", err)
	}
	if c.Transport == nil {
		panic("Transport is nil")
	}
	if resp == nil {
		panic("resp is nil")
	}
	if req.Async() {
		panic(fmt.Sprintf("async requests are not supported (plugin: %s)",
			req.Plugin()))
	}

	aReq := request{
		Janus:         "message_plugin",
		TransactionID: genTransactionID(),
		AdminSecret:   c.AdminSecret,
		Plugin:        req.Plugin(),
		Request:       req.Build(),
	}
	aRes := pluginResponse{
		Response: resp,
	}

	err := c.Transport.Request(ctx, &aReq, &aRes)
	if err != nil {
		return wrap(&TransportError{err})
	}
	if aRes.Error != nil && aRes.Error.Code != 0 {
		return wrap(aRes.Error)
	}
	if ep, ok := resp.(Errer); ok {
		if err := ep.Err(); err != nil {
			return wrap(err)
		}
	}

	return nil
}

func (c *Admin) Sessions(ctx context.Context) ([]int64, error) {
	wrap := func(err error) error {
		return wrapErr("Admin.Sessions", err)
	}
	if c.Transport == nil {
		panic("Transport is nil")
	}

	aReq := request{
		Janus:         "list_sessions",
		TransactionID: genTransactionID(),
		AdminSecret:   c.AdminSecret,
	}
	aRes := sessionResponse{}

	err := c.Transport.Request(ctx, &aReq, &aRes)
	if err != nil {
		return nil, wrap(&TransportError{err})
	}
	if aRes.Error != nil && aRes.Error.Code != 0 {
		return nil, wrap(aRes.Error)
	}

	return aRes.Sessions, nil
}

func (c *Admin) Ping(ctx context.Context) error {
	wrap := func(err error) error {
		return wrapErr("Admin.Ping", err)
	}
	if c.Transport == nil {
		panic("Transport is nil")
	}

	aReq := request{
		Janus:         "ping",
		TransactionID: genTransactionID(),
		AdminSecret:   c.AdminSecret,
	}
	aRes := pingResponse{}

	err := c.Transport.Request(ctx, &aReq, &aRes)
	if err != nil {
		return wrap(&TransportError{err})
	}
	if aRes.Error != nil && aRes.Error.Code != 0 {
		return wrap(aRes.Error)
	}
	if aRes.Janus != "pong" {
		return wrap(fmt.Errorf("unexpected value returned: %s", aRes.Janus))
	}

	return nil
}

func genTransactionID() string {
	return rndString(transactionIDLen)
}
