// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dirserver is a wrapper for a directory implementation that presents it as an
// HTTP interface that stores its data in a GCS implementation
// of the Store service.
package main // import "gcp.upspin.io/cmd/dirserver-gcp"

import (
	cloudLog "upspin.io/cloud/log"
	"upspin.io/flags"
	"upspin.io/log"
	"upspin.io/metric"
	"upspin.io/serverutil/dirserver"

	"gcp.upspin.io/cloud/gcpmetric"
	"gcp.upspin.io/cloud/https"

	// TODO: Which of these are actually needed?

	// Load useful packers
	_ "upspin.io/pack/ee"
	_ "upspin.io/pack/eeintegrity"
	_ "upspin.io/pack/plain"

	// Load required transports
	_ "upspin.io/transports"
)

const (
	serverName    = "dirserver"
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

	ready := dirserver.Main()
	https.ListenAndServe(ready, serverName)
}
