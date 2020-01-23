// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/lfaoro/spark/pkg/kms/gcpkms"
	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) HealthCheck(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}

func (s Server) HealthCheckVerbose(ctx context.Context, empty *empty.Empty) (*pb.HealthCheckResponse, error) {
	var kmsState string
	if kms, ok := s.kms.(*gcpkms.GCPKMS); ok {
		kmsState = kms.State().String()
	} else {
		log.Println("vault: unable to get KMS state")
	}

	return &pb.HealthCheckResponse{
		Database: true,
		Kms:      kmsState,
		Mpi:      false,
		Risk:     false,
		Iin:      s.iin.State().String(),
	}, nil
}
