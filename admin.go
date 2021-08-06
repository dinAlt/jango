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

type response struct {
	Janus       string      `json:"janus,omitempty"`
	Transaction string      `json:"transaction,omitempty"`
	Error       *JanusError `json:"error,omitempty"`
	Response    interface{} `json:"response,omitempty"`
}

type adminRequest struct {
	Janus         string      `json:"janus,omitempty"`
	TransactionID string      `json:"transaction,omitempty"`
	AdminSecret   string      `json:"admin_secret,omitempty"`
	Plugin        string      `json:"plugin,omitempty"`
	Request       interface{} `json:"request,omitempty"`
}

func (r *adminRequest) Transaction() string {
	return r.TransactionID
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

	jError := &JanusError{}
	trID := genTransactionID()
	aReq := adminRequest{"message_plugin", trID, c.AdminSecret, req.Plugin(),
		req.Build()}
	aRes := response{
		Response: resp,
		Error:    jError,
	}

	err := c.Transport.Request(ctx, &aReq, &aRes)
	if err != nil {
		return wrap(&TransportError{err})
	}
	if jError.Code != 0 {
		return wrap(jError)
	}
	if ep, ok := resp.(Errer); ok {
		if err := ep.Err(); err != nil {
			return wrap(err)
		}
	}

	aRes.Response = nil
	aRes.Error = nil

	return nil
}

func genTransactionID() string {
	return rndString(transactionIDLen)
}
