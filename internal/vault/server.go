// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// todo: needs refactor
package vault

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/errorreporting"
	sdformat "github.com/TV4/logrus-stackdriver-formatter"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/juju/ratelimit"
	"github.com/lfaoro/pkg/logger"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	grpc_logrus "github.com/lfaoro/spark/pkg/interceptor/logging/logrus"
	"github.com/lfaoro/spark/pkg/mpi"

	"github.com/lfaoro/spark/ent"
	"github.com/lfaoro/spark/internal"
	"github.com/lfaoro/spark/pkg/iin"
	grpc_auth "github.com/lfaoro/spark/pkg/interceptor/auth"
	"github.com/lfaoro/spark/pkg/kms"
	"github.com/lfaoro/spark/pkg/risk"

	pb_billing "github.com/lfaoro/spark/proto/api/billing"
	pb_user "github.com/lfaoro/spark/proto/api/user"
	pb "github.com/lfaoro/spark/proto/api/vault"
)

var log = logger.New("[vault]")

type Option func(*Server)

func WithDebug(v bool) Option {
	return func(s *Server) {
		s.debug = v
	}
}

func WithIIN(service iin.Service) Option {
	return func(s *Server) {
		s.iin = service
	}
}

func WithRisk(engine risk.Service) Option {
	return func(s *Server) {
		s.risk = engine
	}
}

func WithReporting(client *errorreporting.Client) Option {
	return func(s *Server) {
		s.reporting = client
	}
}

func WithAuth(client pb_user.UserClient) Option {
	return func(s *Server) {
		s.auth = client
	}
}

// Server implements proto/card.proto.
type Server struct {
	ctx context.Context

	db        *ent.Client
	auth      pb_user.UserClient
	billing   pb_billing.BillingClient
	reporting *errorreporting.Client

	kms  kms.Service
	iin  iin.Service  // todo: get api key from service
	risk risk.Service // todo: implement risk engine
	mpi  mpi.Service  // todo: implement

	err   error
	debug bool
}

// New configures a new server.
func New(ctx context.Context, db *ent.Client, kms kms.Service, opts ...Option) Server {
	srv := &Server{
		ctx: ctx,
		db:  db,
		kms: kms,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.debug {
		grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
		grpclog.SetLoggerV2(grpcLog)
	}

	return *srv
}

func (s Server) ServeGRPC(grpcPort string) error {
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return err
	}

	sdLog := logrus.New()
	sdLog.SetFormatter(sdformat.NewFormatter(
		sdformat.WithService("vault"),
		sdformat.WithVersion("1.0.0"),
	))
	// lmt := newVaultLimiter(time.Minute, 5, 1, time.Second*10)
	srv := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			// grpc_ratelimt.UnaryServerInterceptor(lmt),
			grpc_auth.UnaryServerInterceptor(s.authFunc),
			grpc_validator.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(&logrus.Entry{Logger: sdLog},
				grpc_logrus.WithReporting(s.reporting)),
			grpc_opentracing.UnaryServerInterceptor(),
		)))

	pb.RegisterVaultCardServer(srv, s)

	return srv.Serve(listener)
}

func (s Server) ServeHTTP(httpPort, grpcPort string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			"application/json+pretty",
			&runtime.JSONPb{Indent: "  "}),
	)

	wrapper := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Powered-By", "The Fireblaze Engineers")
			h.ServeHTTP(w, r)
		})
	}

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterVaultCardHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		return err
	}

	srv := internal.SecureServer(httpPort, wrapper(mux))
	return srv.ListenAndServe()
}

// todo: add rate limiting
type vaultLimiter struct {
	limiter         *ratelimit.Bucket
	maxWaitDuration time.Duration
}

func newVaultLimiter(fillInterval time.Duration, capacity, quantum int64, maxWaitDuration time.Duration) *vaultLimiter {
	return &vaultLimiter{
		limiter:         ratelimit.NewBucketWithQuantum(fillInterval, capacity, quantum),
		maxWaitDuration: maxWaitDuration,
	}
}

func (v *vaultLimiter) Limit() bool {
	return !v.limiter.WaitMaxDuration(1, v.maxWaitDuration)
}
