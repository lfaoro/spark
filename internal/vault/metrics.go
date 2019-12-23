// Copyright (c) 2019 Leonardo Faoro. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vault

// WIP

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "vault"

var metricStoredCards = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "stored_cards",
	})

var metricQueryCard = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "queried_cards",
	})

var metricDeleteCard = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "deleted_cards",
	})

var metricHighRisk = promauto.NewCounter(
	prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "highrisk_cards",
	})

const promPort = "9090"

func handleMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":"+promPort, nil)
	if err != nil {
		log.Fatal(err)
	}
}
