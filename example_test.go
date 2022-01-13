package jango_test

import (
	"context"
	"fmt"

	"github.com/dinalt/jango"
	"github.com/dinalt/jango/transport"
	"github.com/dinalt/jango/videoroom"
)

func Example() {
	admin := jango.Admin{
		AdminSecret: "janusoverlord",
		Transport: &transport.HTTP{
			URL: "http://localhost:7088/admin",
		},
	}

	ctx := context.Background()

	// Ping Janus server
	if err := admin.PingCtx(ctx); err != nil {
		panic(err)
	}

	// Get sessions list
	sessions, err := admin.SessionsCtx(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sessions count: %d\n", len(sessions))

	// Plugin request
	req := videoroom.ListRequest{}
	res := videoroom.ListResponse{}
	if err := admin.PluginRequestCtx(ctx, &req, &res); err != nil {
		panic(err)
	}
	// We don't need to check res.Err() here, PluginRequestCtx do it for us.
	fmt.Printf("rooms count: %d\n", len(res.List))
}
