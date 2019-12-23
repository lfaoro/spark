// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/lfaoro/spark/ent"
	"github.com/lfaoro/spark/ent/apikey"
	"github.com/lfaoro/spark/internal"
	pb "github.com/lfaoro/spark/proto/api/user"
)

func (s Server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Aborted,
			"failed to start transation: %v", err)
	}

	user, err := tx.User.Create().
		SetEmail(r.Email).
		SetPassword(r.Password).
		Save(ctx)
	if err != nil {
		err = rollback(tx, err)
		return nil, status.Errorf(codes.Aborted, "failed to create user: %v", err)
	}

	project, err := tx.Project.Create().
		AddOwner(user).
		SetName(r.Project).
		Save(ctx)
	if err != nil {
		err = rollback(tx, err)
		return nil, status.Errorf(codes.Aborted, "failed to create project: %v", err)
	}

	key, err := tx.ApiKey.Create().
		SetOwner(project).
		SetName(r.Project + "_key").
		SetKey(internal.GenKey(false)).
		SetEnabled(true).
		Save(ctx)
	if err != nil {
		err = rollback(tx, err)
		return nil, status.Errorf(codes.Aborted, "failed to create key: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, status.Errorf(codes.Aborted,
			"failed to commit transaction: %v", err)
	}

	return &pb.CreateUserResponse{
		Disabled:  user.Disabled,
		ApiKey:    key.Key,
		ProjectId: int32(project.ID),
	}, nil
}

func (s Server) AuthUser(ctx context.Context, r *pb.AuthUserRequest) (*pb.AuthUserResponse, error) {
	key := r.GetKey()
	if key == "" {
		return nil, status.Errorf(codes.Unauthenticated,
			"failed to authenticate with empty key")
	}

	project, err := s.db.ApiKey.Query().
		Where(apikey.Key(key)).
		QueryOwner().
		First(ctx)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Unauthenticated,
			"failed to authorize: %v", err)
	}

	if project.Disabled == true {
		return nil, status.Errorf(codes.PermissionDenied,
			"authorized but disabled project %v with ID %v", project.Name, project.ID)
	}

	return &pb.AuthUserResponse{
		Authorized: true,
		ProjectId:  int64(project.ID),
	}, nil
}

// rollback calls to tx.Rollback and wraps the given error
// with the rollback error if occurred.
func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%v: %v", err, rerr)
	}
	log.Println(err)
	return err
}
