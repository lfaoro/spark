// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"net/http"
	"time"
)

// SecureServer returns an http.Server with sane timeouts.
func SecureServer(hostPort string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              hostPort,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
