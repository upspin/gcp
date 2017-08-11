// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"crypto/sha256"
	"fmt"

	"golang.org/x/oauth2/google"
	dns "google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"

	"upspin.io/errors"
	"upspin.io/upspin"
)

const (
	dnsProject = "upspin-prod"
	dnsZone    = "upspin-services"
	dnsDomain  = "upspin.services"
)

// userToHost converts an Upspin user name to a fully-qualified domain name
// under the upspin.services domain. The host portion of the name is the
// hex-encoded first 16 bytes of the SHA256 checksum of the user name.
// The security of this service relies on their not being collisions in this
// space, which should be astronomically unlikely.
func userToHost(name upspin.UserName) string {
	hash := sha256.New()
	hash.Write([]byte(name))
	return fmt.Sprintf("%x."+dnsDomain, hash.Sum(nil)[:16])
}

func (s *server) setupDNSService() error {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	if err != nil {
		return err
	}
	s.dnsSvc, err = dns.New(client)
	return err
}

func (s *server) lookupName(name upspin.UserName) (ip, host string, err error) {
	host = userToHost(name)

	resp, err := s.dnsSvc.ResourceRecordSets.List(dnsProject, dnsZone).Name(host + ".").Do()
	if err != nil {
		// TODO: handle not found
		return "", "", err
	}
	for _, rrs := range resp.Rrsets {
		for _, rrd := range rrs.Rrdatas {
			return rrd, host, nil
		}
	}

	return "", "", errors.E(errors.NotExist)
}

func (s *server) updateName(name upspin.UserName, ip string) (host string, err error) {
	host = userToHost(name)

	resp, err := s.dnsSvc.ResourceRecordSets.List(dnsProject, dnsZone).Name(host + ".").Do()
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
		Deletions: resp.Rrsets,
	}
	change, err = s.dnsSvc.Changes.Create(dnsProject, dnsZone, change).Do()
	if err != nil && !googleapi.IsNotModified(err) {
		return "", err
	}

	return host, nil
}
