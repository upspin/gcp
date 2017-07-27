// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Keyserver is a wrapper for a key implementation that presents it as an HTTP
// interface that stores the keys on Google Cloud Storage.
package main // import "gcp.upspin.io/cmd/keyserver-gcp"

import (
	"flag"

	cloudLog "gcp.upspin.io/cloud/log"
	"upspin.io/log"
	"upspin.io/metric"
	"upspin.io/serverutil/keyserver"

	"gcp.upspin.io/cloud/gcpmetric"
	"gcp.upspin.io/cloud/https"

	// Load required transports
	_ "upspin.io/key/transports"

	// Storage on GCS.
	_ "gcp.upspin.io/cloud/storage/gcs"
)

const (
	// serverName is the name of this server.
	serverName = "keyserver"

	// metricSampleSize is the size of the sample from which pick one metric
	// to save.
	metricSampleSize = 100

	// metricMaxQPS is the maximum number of metric batches to save per
	// second.
	metricMaxQPS = 5
)

func main() {
	project := flag.String("project", "", "GCP `project` name")

	keyserver.Main(nil)

	if *project != "" {
		cloudLog.Connect(*project, serverName)
		// Disable logging locally so we don't pay the price of local
		// unbuffered writes on a busy server.
		log.SetOutput(nil)
		svr, err := gcpmetric.NewSaver(*project, metricSampleSize, metricMaxQPS, "serverName", serverName)
		if err != nil {
			log.Fatalf("Can't start a metric saver for GCP project %q: %s", *project, err)
		}
		metric.RegisterSaver(svr)
	}

	https.ListenAndServe(nil, serverName)
}
