// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"google.golang.org/grpc"

	pb "github.com/lfaoro/spark/proto/api/user"
)

func NewClient(host string) pb.UserClient {
	if host == "" {
		log.Fatal("environment variable FB_AUTH_HOST is not set. Cannot continue without authorization server.")
	}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
	}

	return pb.NewUserClient(conn)
}
