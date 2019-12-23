// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/status"
)

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request, err error) {

	const fallback = `{"error": "failed to marshal error message"}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(status.Code(err)))

	var message string
	wrap := errors.Unwrap(err)
	if wrap != nil {
		message = wrap.Error()
	}

	jErr := json.NewEncoder(w).Encode(errorBody{
		Code:    runtime.HTTPStatusFromCode(status.Code(err)),
		Error:   status.Convert(err).Message(),
		Message: message,
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

type errorBody struct {
	Error string `protobuf:"bytes,100,name=error" json:"error"`
	// This is to make the error more compatible with users that expect errors to be Status objects:
	// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
	// It should be the exact same message as the Error field.
	Code    int    `protobuf:"varint,1,name=code" json:"code"`
	Message string `protobuf:"bytes,2,name=message" json:"message"`
}
