// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The upspin-setupstorage-gcp command is an external upspin subcommand that
// executes the second step in establishing an upspinserver for GCP.
// Run upspin setupstorage-gcp -help for more information.
package main // import "gcp.upspin.io/cmd/upspin-setupstorage-gcp"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	iam "google.golang.org/api/iam/v1"
	storage "google.golang.org/api/storage/v1"

	"upspin.io/subcmd"
)

type state struct {
	*subcmd.State
}

const help = `
Setupstorage-gcp is the second step in establishing an upspinserver,
It sets up GCS storage for your Upspin installation. You may skip this step
if you wish to store Upspin data on your server's local disk.
The first step is 'setupdomain' and the final step is 'setupserver'.

Setupstorage-gcp creates a Google Cloud Storage bucket and a service account for
accessing that bucket. It then writes the service account private key to
$where/$domain/serviceaccount.json and updates the server configuration files
in that directory to use the specified bucket.

Before running this command, you must create a Google Cloud Project and
associated Billing Account using the Cloud Console:
	https://cloud.google.com/console
The project ID can be any available string, but for clarity it's helpful to
pick something that resembles your domain name.

You must also install the Google Cloud SDK:
	https://cloud.google.com/sdk/downloads
Authenticate and enable the necessary APIs:
	$ gcloud auth login
	$ gcloud --project <project> beta service-management enable iam.googleapis.com storage_api
And, finally, authenticate again in a different way:
	$ gcloud auth application-default login

Running this command when the service account or bucket exists is a no-op.
`

func main() {
	const name = "setupstorage-gcp"

	log.SetFlags(0)
	log.SetPrefix("upspin setupstorage-gcp: ")

	s := &state{
		State: subcmd.NewState(name),
	}

	where := flag.String("where", filepath.Join(os.Getenv("HOME"), "upspin", "deploy"), "`directory` to store private configuration files")
	domain := flag.String("domain", "", "domain `name` for this Upspin installation")
	project := flag.String("project", "", "GCP `project` name")

	s.ParseFlags(flag.CommandLine, os.Args[1:], help,
		"setupstorage-gcp -domain=<name> -project=<gcp_project_name> <bucket_name>")
	if flag.NArg() != 1 {
		s.Exitf("a single bucket name must be provided")
	}
	if *domain == "" || *project == "" {
		s.Exitf("the -domain and -project flags must be provided")
	}

	bucket := flag.Arg(0)

	cfgPath := filepath.Join(*where, *domain)
	cfg := s.ReadServerConfig(cfgPath)

	email, privateKeyData := s.createServiceAccount(*project, cfgPath)

	s.createBucket(*project, email, bucket)

	cfg.StoreConfig = []string{
		"backend=GCS",
		"defaultACL=publicRead",
		"gcpBucketName=" + bucket,
		"privateKeyData=" + privateKeyData,
	}
	s.WriteServerConfig(cfgPath, cfg)

	fmt.Fprintf(os.Stderr, "You should now deploy the upspinserver binary and run 'upspin setupserver'.\n")

	s.ExitNow()
}

func (s *state) createServiceAccount(project, cfgPath string) (email, privateKeyData string) {
	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		// TODO: ask the user to run 'gcloud auth application-default login'
		s.Exit(err)
	}
	svc, err := iam.New(client)
	if err != nil {
		s.Exit(err)
	}

	name := "projects/" + project
	req := &iam.CreateServiceAccountRequest{
		AccountId: "upspinstorage", // TODO(adg): flag?
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: "Upspin Storage",
		},
	}
	created := true
	acct, err := svc.Projects.ServiceAccounts.Create(name, req).Do()
	if isExists(err) {
		// This should be the name we need to get.
		// TODO(adg): make this more robust by listing instead.
		guess := name + "/serviceAccounts/upspinstorage@" + project + ".iam.gserviceaccount.com"
		acct, err = svc.Projects.ServiceAccounts.Get(guess).Do()
		if err != nil {
			s.Exit(err)
		}
		created = false
	} else if err != nil {
		s.Exit(err)
	}

	name += "/serviceAccounts/" + acct.Email
	req2 := &iam.CreateServiceAccountKeyRequest{}
	key, err := svc.Projects.ServiceAccounts.Keys.Create(name, req2).Do()
	if err != nil {
		s.Exit(err)
	}
	if created {
		fmt.Fprintf(os.Stderr, "Service account %q created.\n", acct.Email)
	} else {
		fmt.Fprintf(os.Stderr, "A new key for the service account %q was created.\n", acct.Email)
	}

	return acct.Email, key.PrivateKeyData
}

func (s *state) createBucket(project, email, bucket string) {
	client, err := google.DefaultClient(context.Background(), storage.DevstorageFullControlScope)
	if err != nil {
		// TODO: ask the user to run 'gcloud auth application-default login'
		s.Exit(err)
	}
	svc, err := storage.New(client)
	if err != nil {
		s.Exit(err)
	}

	_, err = svc.Buckets.Insert(project, &storage.Bucket{
		Acl: []*storage.BucketAccessControl{{
			Bucket: bucket,
			Entity: "user-" + email,
			Email:  email,
			Role:   "OWNER",
		}},
		Name: bucket,
		// TODO(adg): flag for location
	}).Do()
	if isExists(err) {
		// TODO(adg): update bucket ACL to make sure the service
		// account has access. (For now, we assume that the user
		// created the bucket using this command and that the bucket
		// has the correct permissions.)
		fmt.Fprintf(os.Stderr, "Bucket %q already exists; re-using it.\n", bucket)
	} else if err != nil {
		s.Exit(err)
	} else {
		fmt.Fprintf(os.Stderr, "Bucket %q created.\n", bucket)
	}
}

func isExists(err error) bool {
	if e, ok := err.(*googleapi.Error); ok && len(e.Errors) > 0 {
		for _, e := range e.Errors {
			if e.Reason != "alreadyExists" && e.Reason != "conflict" {
				return false
			}
		}
		return true
	}
	return false
}
