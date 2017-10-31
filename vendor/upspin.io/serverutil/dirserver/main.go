// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Dirserver is a wrapper for a directory implementation that presents it as an
// HTTP interface. It provides the common code used by all dirserver commands.
package dirserver // import "upspin.io/serverutil/dirserver"

import (
	"flag"
	"net/http"

	"upspin.io/config"
	"upspin.io/dir/inprocess"
	"upspin.io/dir/server"
	"upspin.io/errors"
	"upspin.io/flags"
	"upspin.io/log"
	"upspin.io/rpc/dirserver"
	"upspin.io/serverutil/perm"
	"upspin.io/upspin"

	// TODO: Which of these are actually needed?

	// Load useful packers
	_ "upspin.io/pack/ee"
	_ "upspin.io/pack/eeintegrity"
	_ "upspin.io/pack/plain"

	// Load required transports
	_ "upspin.io/transports"
)

var storeServerUser = flag.String("storeserveruser", "", "`user name` of the StoreServer")

func Main() (ready chan<- struct{}) {
	flags.Parse(flags.Server, "kind", "serverconfig")

	// Load configuration and keys for this server. It needs a real upspin username and keys.
	cfg, err := config.FromFile(flags.Config)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new store implementation.
	var dir upspin.DirServer
	err = nil
	switch flags.ServerKind {
	case "inprocess":
		dir = inprocess.New(cfg)
	case "server":
		dir, err = server.New(cfg, flags.ServerConfig...)
	default:
		err = errors.Errorf("bad -kind %q", flags.ServerKind)
	}
	if err != nil {
		log.Fatalf("Setting up DirServer: %v", err)
	}

	// Wrap with permission checks, if requested.
	if *storeServerUser != "" {
		readyCh := make(chan struct{})
		ready = readyCh
		dir = perm.WrapDir(cfg, readyCh, upspin.UserName(*storeServerUser), dir)
	} else {
		log.Printf("Warning: no Writers Group file protection -- all access permitted")
	}

	httpDir := dirserver.New(cfg, dir, upspin.NetAddr(flags.NetAddr))
	http.Handle("/api/Dir/", httpDir)

	return ready
}
