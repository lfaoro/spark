// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

type Service interface {
	Put(msg Message) error
	Pull(topic string, msgC chan<- string) error
}

type Message struct {
	Topic string
	Data  []byte
}
