// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"context"

	"github.com/lfaoro/pkg/logger"

	"github.com/lfaoro/spark/ent"
	"github.com/lfaoro/spark/ent/migrate"

	_ "github.com/lib/pq"
)

var log = logger.New("[internal]")

func NewDBClient(ctx context.Context, connection string, debug bool) *ent.Client {
	var err error
	db, err := ent.Open("postgres", connection)
	log.FatalIfErr(err)
	if debug {
		// drop columns during debug (local development)
		err = db.Schema.Create(ctx,
			migrate.WithGlobalUniqueID(true),
			migrate.WithDropColumn(true),
		)
	} else {
		err = db.Schema.Create(ctx,
			migrate.WithGlobalUniqueID(true),
		)
	}
	if err != nil {
		log.Fatalf("failed creating or migrating schema resources: %v", err)
	}
	return db
}


