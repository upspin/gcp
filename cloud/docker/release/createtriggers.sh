#!/bin/bash

# This script creates the Google Container Builder trigger to build release
# binaries on each push to the upspin repo's master branch.
# It should be called by an administrator.

project=upspin-test
auth="$(gcloud config config-helper --format='value(credential.access_token)')"
sed "s/PROJECT_ID/$project/" trigger.yaml | curl -X POST -T - \
	-H "Authorization: Bearer $auth" \
	https://cloudbuild.googleapis.com/v1/projects/upspin-test/triggers
