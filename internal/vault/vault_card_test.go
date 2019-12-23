// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

import (
	"context"
	"testing"

	"github.com/lfaoro/spark/internal"
	"github.com/lfaoro/spark/pkg/kms/aesgcm"
	pb "github.com/lfaoro/spark/proto/api/vault"
)

func TestServer_PutCard(t *testing.T) {
	ctx := context.Background()
	ctx = internal.StoreValueIn(ctx, "userid", 12884901889)
	conn := "host=localhost port=5432 user=postgres dbname=vault_test password=develop sslmode=disable"
	kms, _ := aesgcm.New("")
	s := Server{
		db:  internal.NewDBClient(context.Background(), conn, true),
		kms: kms,
	}

	req := &pb.PutCardRequest{
		Holder:     "leonardo",
		Number:     "4590260001666286",
		ExpMonth:   1,
		ExpYear:    2022,
		Cvc:        123,
		AutoDelete: pb.PutCardRequest_NEVER,
	}

	_, err := s.PutCard(ctx, req)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkServer_PutCard(b *testing.B) {
	t := &testing.T{}
	for i := 0; i < b.N; i++ {
		TestServer_PutCard(t)
	}
}
