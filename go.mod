module github.com/lfaoro/spark

go 1.13

require (
	cloud.google.com/go v0.50.0
	cloud.google.com/go/pubsub v1.1.0
	github.com/TV4/logrus-stackdriver-formatter v0.1.0
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/facebookincubator/ent v0.0.0-20191223124427-99510a458dda
	github.com/golang/groupcache v0.0.0-20191027212112-611e8accdfc9 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/joho/godotenv v1.3.0
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lfaoro/pkg v0.5.0
	github.com/lib/pq v1.3.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.3.0
	github.com/rakyll/autopprof v0.1.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	go.opencensus.io v0.22.2 // indirect
	go.uber.org/zap v1.10.0
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6 // indirect
	golang.org/x/sys v0.0.0-20191220220014-0732a990476f // indirect
	google.golang.org/api v0.15.0
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/genproto v0.0.0-20191220175831-5c49e3ecc1c1
	google.golang.org/grpc v1.26.0
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.26.0
