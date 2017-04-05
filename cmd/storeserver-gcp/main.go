// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Storeserver-gcp is a wrapper for a store implementation that presents it as an
// HTTP interface to Google Cloud Storage.
package main // import "gcp.upspin.io/cmd/storeserver-gcp"

import (
	cloudLog "upspin.io/cloud/log"
	"upspin.io/flags"
	"upspin.io/log"
	"upspin.io/metric"
	"upspin.io/serverutil/storeserver"

	"gcp.upspin.io/cloud/gcpmetric"

	// Storage on GCS.
	_ "gcp.upspin.io/cloud/storage/gcs"
)

const (
	serverName    = "storeserver"
	samplingRatio = 1    // report all metrics
	maxQPS        = 1000 // unlimited metric reports per second
)

func main() {
	flags.Register("project")

	if flags.Project != "" {
		cloudLog.Connect(flags.Project, serverName)
		svr, err := gcpmetric.NewSaver(flags.Project, samplingRatio, maxQPS, "serverName", serverName)
		if err != nil {
			log.Fatalf("Can't start a metric saver for GCP project %q: %s", flags.Project, err)
		} else {
			metric.RegisterSaver(svr)
		}
	}

	storeserver.Main()
}
