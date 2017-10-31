#!/bin/bash -e

# This script runs inside the upspinserver-build Docker container.
# It copies the gcp.uspin.io tree out of /workspace (a volume provided by
# Google Container Builder) and into GOPATH, builds upspinserver-gcp, and puts
# the resulting binary in /workspace/bin.

cp -R /workspace /go/src/gcp.upspin.io
rm -r /workspace/*
cp /go/src/gcp.upspin.io/cloud/docker/upspinserver/* /workspace
mkdir /workspace/bin
/usr/local/go/bin/go build -o /workspace/bin/upspinserver-gcp gcp.upspin.io/cmd/upspinserver-gcp
