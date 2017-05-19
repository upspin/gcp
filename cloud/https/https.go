// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package https provides a helper function to set up an HTTPS server derived
// from the current environment.
package https // import "gcp.upspin.io/cloud/https"

import (
	"log"

	"gcp.upspin.io/cloud/autocert"

	"cloud.google.com/go/compute/metadata"
	"upspin.io/cloud/https"
)

// ListenAndServe serves the http.DefaultServeMux by HTTPS (and HTTP,
// redirecting to HTTPS), configured using the given options (or the
// command-line flags if opt is nil).
//
// If running on GCE and the -letscache flag is *not* specified, ListenAndServe
// will configure a Let's Encrypt cache that stores its data in the Cloud
// Storage bucket specified by the "letsencrypt-bucket" instance attribute in
// the Compute Engine Metadata server.
//
// The given channel, if any, is closed when the TCP listener has succeeded. It
// may be used to signal that the server is ready to start serving requests.
//
// ListenAndServe does not return. It exits the program when the server is shut
// down (via SIGTERM or due to an error) and calls shutdown.Shutdown.
func ListenAndServe(ready chan<- struct{}, opt *https.Options, serverName string) {
	if opt == nil {
		opt = https.OptionsFromFlags()
	}
	if metadata.OnGCE() && opt.LetsEncryptCache == "" {
		const key = "letsencrypt-bucket"
		bucket, err := metadata.InstanceAttributeValue(key)
		if err != nil {
			log.Fatalf("https: couldn't read %q metadata value: %v", key, err)
		}
		cache, err := autocert.NewCache(bucket, serverName)
		if err != nil {
			log.Fatalf("https: couldn't set up letsencrypt cache: %v", err)
		}
		opt.AutocertCache = cache
	}
	https.ListenAndServe(ready, opt)
}
