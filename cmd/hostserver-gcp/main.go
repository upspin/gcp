// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command hostserver-gcp is a combined DirServer and StoreServer that serves a
// synthetic Upspin tree used to configure sub-domains under the domain
// upspin.services. Each Upspin user can set the IP address of one sub-domain,
// the name of which is a hash of their user name.
//
// Assuming the server is running as host@upspin.io, here's how the user
// user@example.com would configure their sub-domain:
//   $ upspin mkdir host@upspin.io/user@example.com/127.0.0.1
// If that command succeeds, then their host name was created or updated.
// (This can be any command that performs a DirServer.Put to that path; the put
// command would wokr, but mkdir is the simpliest choice in this case.)
//
// To find the host name, the user uses a get request, which returns both the
// host name and configured IP address:
//   $ upspin get host@example.com/user@example.com
//   127.0.0.1
//   b4c9a289323b21a01c3e940f150eb9b8.upspin.services
//
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
