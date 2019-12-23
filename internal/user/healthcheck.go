// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s Server) HealthCheck(context.Context, *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}