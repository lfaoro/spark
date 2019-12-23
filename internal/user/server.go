// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"net"
	"os"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/lfaoro/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/lfaoro/spark/internal"
	pb "github.com/lfaoro/spark/proto/api/user"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"

	"github.com/lfaoro/spark/ent"
)

var log = logger.New("[internal/user]")
var grpcLog = grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)

type Option func(*Server)

func WithDebug(v bool) Option {
	return func(s *Server) {
		s.debug = v
	}
}

// Server implements proto/card.proto.
type Server struct {
	ctx context.Context
	db  *ent.Client

	err   error
	debug bool
}

// New setups a new Server.
func New(ctx context.Context, db *ent.Client, opts ...Option) Server {
	s := &Server{
		ctx: ctx,
		db:  db,
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.debug {
		grpclog.SetLoggerV2(grpcLog)
	}

	return *s
}

func (s Server) ServeGRPC(grpcPort string) error {
	l, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	sdLog := logrus.New()
	sdLog.SetFormatter(stackdriver.NewFormatter(
		stackdriver.WithService("user"),
		stackdriver.WithVersion("1.0.0"),
	))

	// lmt := newVaultLimiter(time.Minute, 5, 1, time.Second*10)
	srv := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(&logrus.Entry{
				Logger: sdLog,
			}),
			// grpc_ratelimt.UnaryServerInterceptor(lmt),
			grpc_validator.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
		)))

	pb.RegisterUserServer(srv, s)

	return srv.Serve(l)
}

func (s Server) ServeHTTP(httpPort, grpcPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	runtime.HTTPError = CustomHTTPError
	opts := []grpc.DialOption{grpc.WithInsecure()}
	mux := runtime.NewServeMux()

	err := pb.RegisterUserHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		return err
	}

	srv := internal.SecureServer(httpPort, mux)

	return srv.ListenAndServe()
}

func (s Server) Close() {
	err := s.db.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("database connection has been closed")
	// close grpc
	// close http
}
