#!/bin/bash

# This script deploys the Upspin servers running under upspin.io.

go install && upspin-deploy-gcp -project=upspin-prod -domain=upspin.io -zone=us-central1-c -keyserver="" "$@"
