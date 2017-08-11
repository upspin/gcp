// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command hostserver-gcp ...
package main // import "gcp.upspin.io/cmd/hostserver-gcp"

import (
	"log"
	"net/http"

	"gcp.upspin.io/cloud/https"

	"upspin.io/config"
	"upspin.io/flags"
	"upspin.io/rpc/dirserver"
	"upspin.io/rpc/storeserver"
	"upspin.io/upspin"
)

func main() {
	flags.Parse(flags.Server)

	addr := upspin.NetAddr(flags.NetAddr)
	ep := upspin.Endpoint{
		Transport: upspin.Remote,
		NetAddr:   addr,
	}
	cfg, err := config.FromFile(flags.Config)
	if err != nil {
		log.Fatal(err)
	}

	s, err := newServer(ep, cfg)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/api/Store/", storeserver.New(cfg, s.StoreServer(), addr))
	http.Handle("/api/Dir/", dirserver.New(cfg, s.DirServer(), addr))

	https.ListenAndServe(nil, "hostserver")
}
