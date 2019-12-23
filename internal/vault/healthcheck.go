// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) HealthCheck(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s Server) HealthCheckVerbose(context.Context, *empty.Empty) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Database: true,
		Kms:      true,
		Mpi:      true,
		Risk:     true,
		Iin:      true,
	}, nil
}
