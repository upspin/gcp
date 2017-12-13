#!/bin/bash -e

# vendor-update.sh updates the vendored copy of the upspin.io repository and
# stages the result, ready for "git commit".

# The dep command can be obtained with "go get github.com/golang/dep/cmd/dep".

# Update the upspin.io package.
dep ensure -update upspin.io
# Remove any vendored packages we don't use.
dep prune
# Delete test files.
find vendor -name '*_test.go' -delete
# Delete Google Cloud JSON API schemas.
find vendor -name '*-api.json' -delete

git add vendor Gopkg.lock
git gofmt
