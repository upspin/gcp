#!/bin/bash

# This script creates the Google Container Builder trigger to build
# the upspinserver image whenever the gcp master branch is updated.
# It should be run only by an administrator.

# By default it deploys to upspin-test.
# Use the -prod flag to deploy to production.

project="upspin-test"
suffix="-test"
if [[ "$@" == "-prod" ]]; then
	project="upspin-prod"
	suffix=""
elif [[ "$@" != "" ]]; then
	echo "usage: $0 [-prod]"
	exit 2
fi
auth="$(gcloud config config-helper --format='value(credential.access_token)')"
<trigger.yaml sed "s/PROJECT_ID/$project/" | sed "s/SUFFIX/$suffix/" | curl -X POST -T - \
	-H "Authorization: Bearer $auth" \
	https://cloudbuild.googleapis.com/v1/projects/$project/triggers
