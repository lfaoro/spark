// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func GracefulShutdown(startTime time.Time) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT)
	log.Printf("caught signal: %+v", <-signalChan)
	// Log server shutdown
	host, _ := os.Hostname()
	message := fmt.Sprintf("server on %v terminated after %v", host, time.Since(startTime))
	log.Println(message)
}
