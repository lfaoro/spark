.PHONY: proto clean install schema run helm

#export GOFLAGS="-mod=vendor"

run:
	@killall -9 gin &>/dev/null || :
	@killall -9 user || :
	@go install ./cmd/user
	user&
	gin -i

run-db:
	@bash hack/postgres.sh

# todo: https://gist.github.com/c4milo/01e344c369ed2f9e0253ca168047197d
SRC::="api/user"
DST:="./proto"
INCLUDE:=-I. \
			-I/usr/local/include \
			-I . \
			-I${GOPATH}/src \
			-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			-I${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
			 ${SRC}/*.proto
proto:
	@echo "Building proto definitions for ${SRC}"
	# @mkdir -p ${DST}/go
	# @mkdir -p ${DST}/js
	# @mkdir -p ${DST}/java
	# @mkdir -p ${DST}/swagger
	@protoc $(INCLUDE) --go_out=plugins=grpc,paths=source_relative:${DST} \
					--validate_out="lang=go:${DST}" \
					--grpc-gateway_out=logtostderr=true:${DST} \
					--swagger_out=logtostderr=true,use_go_templates=true:${DST}

	@protoc $(INCLUDE)	--js_out=import_style=commonjs:${DST}
	#	  protoc $(INCLUDE) --grpc-web_out=import_style=commonjs,mode=grpcwebtext:${DST}/js

	@$(shell bash hack/docs.bash)

schema: clean-ent
	go run github.com/facebookincubator/ent/cmd/entc generate ./ent/schema

generate: schema proto

clean-proto:
	@rm -fr proto/api

clean-ent:
	@pushd ent;ls -1 | grep -v "schema" | xargs rm -r;popd

clean: clean-ent clean-proto

install: install-grpcweb
	@GO111MODULE=off go get -u \
	google.golang.org/grpc \
	github.com/golang/protobuf/{proto,protoc-gen-go} \
	github.com/grpc-ecosystem/grpc-gateway/{protoc-gen-grpc-gateway,protoc-gen-swagger} \
	github.com/envoyproxy/protoc-gen-validate \
	github.com/facebookincubator/ent/cmd/entc \
	github.com/codegangsta/gin

	which protoc || brew install protobuf
	which swagger || brew tap go-swagger/go-swagger && brew install go-swagger

release:="1.0.7" # update from https://github.com/grpc/grpc-web/releases
install-grpcweb:
	@which protoc-gen-grpc-web || hack/grpcweb.bash ${release}

SITE?=doc.fireblaze.io
doc-deploy:
	@gcloud config set project fireblaze
	@gsutil mb gs://$(SITE) || :
	@gsutil -m rsync -R -d web/doc/. gs://${SITE}/
	@gsutil iam ch allUsers:objectViewer gs://${SITE}
	@gsutil web set -m index.html -e index.html gs://${SITE}
	@echo "Visit https://${SITE} to preview the documentation."

LDFLAGS += -X "main.Version=$(shell git describe --abbrev=0 --tags)"
LDFLAGS += -X "main.BuildHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "main.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
release:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go install -ldflags '$(LDFLAGS)' -gcflags "-N -l" .

release-user:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go install -ldflags '$(LDFLAGS)' -gcflags "-N -l" ./cmd/user



docker: docker-user docker-vault

docker-vault:
	docker build -t eu.gcr.io/fireblaze/vault:latest .
	docker push eu.gcr.io/fireblaze/vault:latest

docker-user:
	docker build -t eu.gcr.io/fireblaze/user:latest \
		-f ./cmd/user/Dockerfile .
	docker push eu.gcr.io/fireblaze/user:latest

run-docker:
	docker run -it --rm eu.gcr.io/fireblaze/user
	docker run -it --rm eu.gcr.io/fireblaze/vault

helm:
	helm template kube/spark
	helm upgrade --install vault "kube/spark" \
	--namespace spark --cleanup-on-fail --reset-values --force