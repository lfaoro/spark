// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb_user "github.com/lfaoro/spark/proto/api/user"
)

func (s Server) authFunc(ctx context.Context) (context.Context, error) {
	// skip authorization for health check
	m, ok := grpc.Method(ctx)
	if m == "/fireblaze.vault.v1.Card/HealthCheck" && ok {
		return ctx, nil
	}

	key, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	// let demo key pass
	if key == "key_demo" {
		return ctx, nil
	}

	res, err := s.auth.AuthUser(ctx, &pb_user.AuthUserRequest{
		Key: key,
	})
	if err != nil {
		if status.Code(err) == codes.Unavailable {
			log.Println(err)
			const errUnavailable = "failed to contact authorization service, try again in a few seconds"
			return nil, status.Errorf(status.Code(err), errUnavailable)
		}
		return nil, err
	}

	grpc_ctxtags.Extract(ctx).Set("auth.sub", res.ProjectId)
	ctx = context.WithValue(ctx, "project_id", int(res.ProjectId))
	return ctx, nil
}
