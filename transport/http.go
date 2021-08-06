package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dinalt/jango"
)

var _ jango.Transport = (*HTTP)(nil)

type HTTP struct {
	*http.Client
	UseHTTPS bool
	URL      string
}

func (c *HTTP) Request(ctx context.Context, req interface{}, resp interface{}) error {
	if req == nil {
		panic("req is nil")
	}
	if resp == nil {
		panic("resp is nil")
	}
	cli := c.setup()

	bb := &bytes.Buffer{}
	err := json.NewEncoder(bb).Encode(req)
	if err != nil {
		return fmt.Errorf("request encoding: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.URL, bb)
	if err != nil {
		return fmt.Errorf("http.NewRequestWithContext(...): %w", err)
	}
	hResp, err := cli.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request send: %w", err)
	}
	defer func() {
		if hResp.Body != nil {
			_ = hResp.Body.Close()
		}
	}()

	if hResp.Body != nil {
		err := json.NewDecoder(hResp.Body).Decode(resp)
		if err != nil {
			return fmt.Errorf("response decode: %w", err)
		}
	}
	if hResp.StatusCode > 399 {
		return fmt.Errorf("bad response: %s", hResp.Status)
	}

	return nil
}

func (c *HTTP) setup() *http.Client {
	if c.URL == "" {
		panic("url is not set")
	}

	cli := c.Client
	if cli == nil {
		cli = http.DefaultClient
	}

	return cli
}
