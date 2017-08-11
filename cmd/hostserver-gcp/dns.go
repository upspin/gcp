// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"

	"upspin.io/errors"
	"upspin.io/upspin"
)

// TODO(adg): make this configurable?
const (
	dnsProject = "upspin-prod"
	dnsZone    = "upspin-services"
	dnsDomain  = "upspin.services"
)

// userToHost converts an Upspin user name to a fully-qualified domain name
// under the upspin.services domain. The host portion of the name is the
// hex-encoded first 16 bytes of the SHA256 checksum of the user name.
// The security of this service relies on there not being collisions in this
// space, which should be astronomically unlikely.
func userToHost(name upspin.UserName) string {
	hash := sha256.New()
	hash.Write([]byte(name))
	return fmt.Sprintf("%x."+dnsDomain, hash.Sum(nil)[:16])
}

// setupDNSService loads the credentials for accessing the Cloud DNS service
// and sets the server's dnsSvc with a ready-to-use dns.Service.
func (s *server) setupDNSService() error {
	ctx := context.Background()
	var client *http.Client

	// First try to read the serviceaccount.json in the Docker image.
	b, err := ioutil.ReadFile("/upspin/serviceaccount.json")
	if err == nil {
		cfg, err := google.JWTConfigFromJSON(b, dns.CloudPlatformScope)
		if err != nil {
			return err
		}
		client = cfg.Client(ctx)
	} else if os.IsNotExist(err) {
		// Otherwise use the default application credentials,
		// which should work when testing locally.
		client, err = google.DefaultClient(ctx, dns.CloudPlatformScope)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	s.dnsSvc, err = dns.New(client)
	return err
}

// listRecordSets returns the list of record sets for a given host name.
func (s *server) listRecordSets(host string) ([]*dns.ResourceRecordSet, error) {
	resp, err := s.dnsSvc.ResourceRecordSets.List(dnsProject, dnsZone).Name(host + ".").Do()
	if err != nil {
		return nil, err
	}
	return resp.Rrsets, nil
}

// lookupName returns the IP address and host name for a given user, or a
// NotExist error if there is no host name for that user.
func (s *server) lookupName(name upspin.UserName) (ip, host string, err error) {
	host = userToHost(name)

	rrsets, err := s.listRecordSets(host)
	if err != nil {
		return "", "", err
	}
	for _, rrs := range rrsets {
		for _, rrd := range rrs.Rrdatas {
			return rrd, host, nil
		}
	}

	return "", "", errors.E(errors.NotExist)
}

// updateName creates (or replaces) an A record for the given user's host name
// that points to the given IP address, and returns the user's host name.
func (s *server) updateName(name upspin.UserName, ip string) (host string, err error) {
	host = userToHost(name)

	rrsets, err := s.listRecordSets(host)
	if err != nil {
		return "", err
	}
	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{{
			Name:    host + ".",
			Rrdatas: []string{ip},
			Ttl:     3600, // 1 hour
			Type:    "A",
		}},
		Deletions: rrsets,
	}
	change, err = s.dnsSvc.Changes.Create(dnsProject, dnsZone, change).Do()
	if err != nil && !googleapi.IsNotModified(err) {
		return "", err
	}

	return host, nil
}
