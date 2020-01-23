// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"
	"time"

	"cloud.google.com/go/errorreporting"
	"github.com/lfaoro/pkg/logger"

	"github.com/lfaoro/spark/internal"
	"github.com/lfaoro/spark/internal/user"
	"github.com/lfaoro/spark/internal/vault"
	"github.com/lfaoro/spark/pkg/iin/binlist"

	"github.com/lfaoro/spark/pkg/kms"
	"github.com/lfaoro/spark/pkg/kms/aesgcm"
	"github.com/lfaoro/spark/pkg/kms/gcpkms"
)

var (
	log      = logger.New("[vault]")
	grpcPort = flag.String("grpcPort", ":50051", "host & port the grpc server listens on.")
	httpPort = flag.String("httpPort", ":3001", "host & port the http server listens on.")
	// use during development only, enforces database migrations.
	debugFlag = flag.Bool("debug", false, "prints debug logs.")
)

func main() {
	flag.Parse()
	internal.LoadEnv()

	envDatabase := os.Getenv("DB_CONNECTION")
	envProjectID := os.Getenv("GCP_PROJECTID")
	envLocationID := os.Getenv("GCP_LOCATIONID")
	envAuthHost := os.Getenv("FB_AUTH_HOST")
	envServiceName := os.Getenv("FB_SERVICE_NAME")
	envServiceVersion := os.Getenv("FB_SERVICE_VERSION")

	log.Println("starting grpc server on", *grpcPort)
	log.Println("starting http server on", *httpPort)

	var _kms kms.Service
	var err error
	if *debugFlag {
		_kms, err = aesgcm.New("")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_kms = gcpkms.New(envProjectID, envLocationID, gcpkms.Debug(*debugFlag))
	}

	_ctx := context.Background()

	cfg := errorreporting.Config{
		ServiceName:    envServiceName,
		ServiceVersion: envServiceVersion,
	}
	_reporting, err := errorreporting.NewClient(_ctx, envProjectID, cfg)
	if err != nil {
		log.Println(err)
	}
	defer _reporting.Close()

	_db := internal.NewDBClient(_ctx, envDatabase, *debugFlag)
	defer _db.Close()

	_iin := binlist.New("")

	_auth := user.NewClient(envAuthHost)

	srv := vault.New(_ctx, _db, _kms,
		vault.WithIIN(_iin),
		vault.WithAuth(_auth),
		vault.WithReporting(_reporting),
		vault.WithDebug(*debugFlag),
	)

	go func() { log.Fatal(srv.ServeGRPC(*grpcPort)) }()
	go func() { log.Fatal(srv.ServeHTTP(*httpPort, *grpcPort)) }()

	internal.GracefulShutdown(time.Now())

	// todo: use mux to run http/grpc on same port 80
}
