// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/lfaoro/pkg/logger"

	"github.com/lfaoro/spark/internal"
	"github.com/lfaoro/spark/internal/user"
)

var (
	// Version is manually generated using SemVer guidelines.
	// ref: https://semver.org/
	//
	// It is set via build flag using instruction from the Makefile.
	Version = "devel"
	// BuildTime shows the date this program was built,
	//
	// It is set via build flag using instruction from the Makefile.
	BuildTime string
	// BuildHash is the commit hash from which this program was built.
	//
	// It is set via build flag using instruction from the Makefile.
	BuildHash string
)

var (
	log = logger.New("[cmd/user]")

	grpcPort = flag.String("grpcPort", ":50052", "host & port the grpc server listens on.")
	httpPort = flag.String("httpPort", ":3002", "host & port the http server listens on.")

	// use during development only, enforces database migrations.
	debugFlag = flag.Bool("debug", true, "prints debug logs.")

	versionFlag = flag.Bool("version", false, "program version.")
)

func main() {
	flag.Parse()
	internal.LoadEnv()

	if *versionFlag {
		printVersion()
		return
	}

	// if *debugFlag {
	// 	autopprof.Capture(autopprof.CPUProfile{
	// 		Duration: 30 * time.Second,
	// 	})
	// }

	log.Println("starting grpc server on", *grpcPort)
	log.Println("starting http server on", *httpPort)

	ctx := context.Background()
	db := internal.NewDBClient(ctx, os.Getenv("DB_CONNECTION"), *debugFlag)
	srv := user.New(ctx, db, user.WithDebug(*debugFlag))
	defer srv.Close()

	go func() { log.Fatal(srv.ServeGRPC(*grpcPort)) }()
	go func() { log.Fatal(srv.ServeHTTP(*httpPort, *grpcPort)) }()

	internal.GracefulShutdown(time.Now())
}

func printVersion() {
	output := `Hash: %v
Time: %v
Version: %v`
	fmt.Printf(output+"\n", BuildHash, BuildTime, Version)
}
