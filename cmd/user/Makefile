.PHONY: run release docker

run:
	@killall -9 gin &>/dev/null || :
	@killall -9 gin-bin &>/dev/null || :
	gin -p 3001 -a 3002 -i

LDFLAGS += -X "main.Version=$(shell git describe --abbrev=0 --tags)"
LDFLAGS += -X "main.BuildHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "main.BuildTime=$(shell date -u '+%Y-%m-%d %I:%M:%S %Z')"
release:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	go install -ldflags '$(LDFLAGS)' -gcflags "-N -l" .

docker:
	docker build -t eu.gcr.io/fireblaze/user:latest .
	docker push eu.gcr.io/fireblaze/user:latest
	docker run -it --rm eu.gcr.io/fireblaze/user