#!/bin/bash -e

# This script runs inside the upspinserver-build Docker container.
# It copies the gcp.uspin.io tree out of /workspace (a volume provided by
# Google Container Builder) and into GOPATH, builds upspinserver-gcp, and puts
# the resulting binary in /workspace/bin.

echo "Copying /workspace to $GOPATH/src/gcp.upspin.io"
cp -R /workspace /go/src/gcp.upspin.io
echo "Emptying workspace"
rm -r /workspace/*

echo GOPATH=$GOPATH
find $GOPATH -type f

echo "Copying upspinserver artifacts to workspace"
cp $GOPATH/src/gcp.upspin.io/cloud/docker/upspinserver/* /workspace
mkdir /workspace/bin

cd $GOPATH/src/gcp.upspin.io

# TODO(adg): restore this functionality
#echo "Generating version package"
#/usr/local/go/bin/go generate -run make_version gcp.upspin.io/vendor/upspin.io/version

echo "Building upspinserver-gcp"
/usr/local/go/bin/go build -o /workspace/bin/upspinserver-gcp gcp.upspin.io/cmd/upspinserver-gcp
