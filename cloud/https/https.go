// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package https provides a helper function to set up an HTTPS server derived
// from the current environment.
package https

import (
	"log"

	"gcp.upspin.io/cloud/autocert"

	"cloud.google.com/go/compute/metadata"
	"upspin.io/cloud/https"
)

func ListenAndServe(ready chan struct{}, serverName string) {
	opt := https.OptionsFromFlags()
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
