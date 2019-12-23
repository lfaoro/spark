// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcp_pubsub

import (
	"context"
	"fmt"

	gcp_pubsub "cloud.google.com/go/pubsub"

	"github.com/lfaoro/spark/pkg/pubsub"
)

type GCPPubSub struct {
	projectID string
	client    *gcp_pubsub.Client

	ctx   context.Context
	debug bool
}

type Option func(sub *GCPPubSub)

func Debug(v bool) Option {
	return func(c *GCPPubSub) {
		c.debug = v
	}
}

func New(projectID string, opts ...Option) pubsub.Service {
	ctx := context.Background()
	client, _ := gcp_pubsub.NewClient(ctx, projectID)

	service := &GCPPubSub{
		client: client,
		ctx:    ctx,
	}

	for _, opt := range opts {
		opt(service)
	}

	return service
}

func (g GCPPubSub) Put(msg pubsub.Message) error {
	t := g.client.Topic(msg.Topic)
	res := t.Publish(g.ctx, &gcp_pubsub.Message{
		Data: msg.Data,
	})
	// Block until the result is returned.
	id, err := res.Get(g.ctx)
	if g.debug {
		fmt.Printf("published message with ID: %v\n", id)
	}

	return err
}

func (g GCPPubSub) Pull(topic string, msgC chan<- string) error {
	cctx, cancel := context.WithCancel(g.ctx)
	defer cancel()

	sub := g.client.Subscription(topic)

	err := sub.Receive(cctx, func(ctx context.Context, msg *gcp_pubsub.Message) {
		data := string(msg.Data)
		if data == "" {
			// send Nack if message wasn't successfully processed, so
			// the system will retry.
			msg.Nack()
		}

		if g.debug {
			fmt.Printf("got message: %q\n", data)
		}
		msgC <- data
		// send Ack to successfully process the message.
		msg.Ack()
	})

	return err
}
