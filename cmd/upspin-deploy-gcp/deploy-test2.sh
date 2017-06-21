#!/bin/bash

# This script deploys the Upspin servers running under test.upspin.io.

go install && upspin-deploy -project=upspin-test2 -domain=test2.upspin.io -zone=us-central1-c -keyserver="" "$@"
