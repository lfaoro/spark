// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/lfaoro/spark/proto/api/vault"
)

func (s Server) PutPII(context.Context, *pb.PutPIIRequest) (*pb.PutPIIResponse, error) {
	panic("implement me")
}

func (s Server) GetPII(context.Context, *pb.GetPIIRequest) (*pb.GetPIIResponse, error) {
	panic("implement me")
}

func (s Server) DelPII(context.Context, *pb.DelPIIRequest) (*empty.Empty, error) {
	panic("implement me")
}
