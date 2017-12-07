// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcs

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"upspin.io/cloud/storage"
	"upspin.io/log"
	"upspin.io/upspin"
)

const defaultTestBucketName = "upspin-test-scratch"

var (
	client      storage.Storage
	testDataStr = fmt.Sprintf("This is test at %v", time.Now())
	testData    = []byte(testDataStr)
	fileName    = fmt.Sprintf("test-file-%d", time.Now().Second())

	testBucket = flag.String("test_bucket", defaultTestBucketName, "bucket name to use for testing")
	useGcloud  = flag.Bool("use_gcloud", false, "enable to run google cloud tests; requires gcloud auth login")
)

// This is more of a regression test as it uses the running cloud
// storage in prod. However, since GCP is always available, we accept
// to rely on it.
func TestPutGetAndDownload(t *testing.T) {
	err := client.Put(fileName, testData)
	if err != nil {
		t.Fatalf("Can't put: %v", err)
	}
	data, err := client.Download(fileName)
	if err != nil {
		t.Fatalf("Can't Download: %v", err)
	}
	if string(data) != testDataStr {
		t.Errorf("Expected %q got %q", testDataStr, string(data))
	}
	// Check that Download yields the same data
	bytes, err := client.Download(fileName)
	if err != nil {
		t.Fatal(err)
	}
	if string(bytes) != testDataStr {
		t.Errorf("Expected %q got %q", testDataStr, string(bytes))
	}
}

func TestDelete(t *testing.T) {
	err := client.Put(fileName, testData)
	if err != nil {
		t.Fatal(err)
	}
	err = client.Delete(fileName)
	if err != nil {
		t.Fatalf("Expected no errors, got %v", err)
	}
	// Test the side effect after Delete.
	_, err = client.Download(fileName)
	if err == nil {
		t.Fatal("Expected an error, but got none")
	}
}

func TestList(t *testing.T) {
	ls, ok := client.(storage.Lister)
	if !ok {
		t.Fatal("impl does not provide List method")
	}

	if err := client.(*gcsImpl).emptyBucket(false); err != nil {
		t.Fatal(err)
	}

	refs, next, err := ls.List("")
	if err != nil {
		t.Fatal(err)
	}
	if len(refs) != 0 {
		t.Errorf("list of empty bucket returned %d refs", len(refs))
	}
	if next != "" {
		t.Errorf("list of empty bucket returned non-empty page token %q", next)
	}

	// Test pagination by reducing the results per page to 2.
	oldMaxResults := maxResults
	defer func() { maxResults = oldMaxResults }()
	maxResults = 2

	const nFiles = 6 // Must be evenly divisible by maxResults.
	for i := 0; i < nFiles; i++ {
		err = client.Put(fmt.Sprintf("test-%d", i), testData)
		if err != nil {
			t.Fatal(err)
		}
	}

	seen := make(map[upspin.Reference]bool)
	for i := 0; i < nFiles/2; i++ {
		refs, next, err = ls.List(next)
		if err != nil {
			t.Fatal(err)
		}
		if len(refs) != 2 {
			t.Errorf("got %d refs, want 2", len(refs))
		}
		if i == nFiles/2-1 {
			if next != "" {
				t.Errorf("got page token %q, want empty", next)
			}
		} else if next == "" {
			t.Error("got empty page token, want non-empty")
		}
		for _, ref := range refs {
			if seen[ref.Ref] {
				t.Errorf("saw duplicate ref %q", ref.Ref)
			}
			seen[ref.Ref] = true
			if got, want := ref.Size, int64(len(testData)); got != want {
				t.Errorf("ref %q has size %d, want %d", ref.Ref, got, want)
			}
		}
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	if !*useGcloud {
		log.Printf(`

cloud/storage/gcs: skipping test as it requires GCS access. To enable this test,
ensure you are authenticated to a GCP project that has editor permissions to a
GCS bucket named by flag -test_bucket and then set this test's flag -use_gcloud.

`)
		os.Exit(0)
	}

	// Create client that writes to test bucket.
	var err error
	client, err = storage.Dial("GCS",
		storage.WithKeyValue("gcpBucketName", *testBucket),
		storage.WithKeyValue("defaultACL", PublicRead))
	if err != nil {
		log.Fatalf("cloud/storage/gcs: couldn't set up client: %v", err)
	}

	code := m.Run()

	// Clean up.
	const verbose = true
	if err := client.(*gcsImpl).emptyBucket(verbose); err != nil {
		log.Printf("cloud/storage/gcs: emptyBucket failed: %v", err)
	}

	os.Exit(code)
}
