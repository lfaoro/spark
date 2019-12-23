# build stage
FROM golang:1.13-alpine as app-builder
WORKDIR /build

# environment
ENV GOFLAGS="-mod=vendor"

# add necessary packages
RUN apk add --update --no-cache \
    git gcc make

# cache dependencies
#COPY go.mod .
#COPY go.sum .
#RUN go mod download

COPY . .
RUN make release

# release stage
FROM alpine:latest
WORKDIR /usr/local/bin/
EXPOSE 3001
EXPOSE 50051

RUN apk add --update --no-cache \
    ca-certificates && \
    update-ca-certificates

COPY --from=app-builder /go/bin/ /usr/local/bin/
CMD ["/usr/local/bin/spark", "-debug=false"]
