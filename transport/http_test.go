package transport_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/dinalt/jango"
	"github.com/dinalt/jango/transport"
)

func TestHTTP_Request(t *testing.T) {
	type fields struct {
		Client *http.Client
		URL    string
	}
	type args struct {
		ctx  context.Context
		req  interface{}
		resp interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"valid request",
			fields{
				URL: os.Getenv("JANGOTEST_REQUESTS_URL"),
			},
			args{
				ctx: context.Background(),
				req: &request{
					Janus:         "ping",
					TransactionID: "1",
				},
				resp: &response{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &transport.HTTP{
				Client: tt.fields.Client,
				URL:    tt.fields.URL,
			}
			if err := c.Request(tt.args.ctx, tt.args.req, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("HTTP.Request() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.resp.(*response).Janus != "pong" {
				t.Errorf("HTTP.Request() unexpected response: %+v", *tt.args.resp.(*response))
			}
		})
	}
}

type request struct {
	Janus         string `json:"janus,omitempty"`
	TransactionID string `json:"transaction,omitempty"`
	AdminSecret   string `json:"admin_secret,omitempty"`
}

func (r *request) Transaction() string {
	return r.TransactionID
}

type response struct {
	Janus       string            `json:"janus,omitempty"`
	Transaction string            `json:"transaction,omitempty"`
	Error       *jango.JanusError `json:"error,omitempty"`
}
