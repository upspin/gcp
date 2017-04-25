// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Frontend-gcp provides a web server that serves documentation and meta
// tags to instruct "go get" where to find the upspin source repository.
package main // import "gcp.upspin.io/cmd/frontend-gcp"

import (
	"gcp.upspin.io/cloud/https"
	"upspin.io/serverutil/frontend"
)

func main() {
	frontend.Main()
	https.ListenAndServe(nil, "frontend")
}
