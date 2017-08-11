// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	dns "google.golang.org/api/dns/v1"

	"upspin.io/access"
	"upspin.io/cache"
	"upspin.io/errors"
	"upspin.io/pack"
	"upspin.io/path"
	"upspin.io/upspin"
	"upspin.io/valid"

	_ "upspin.io/pack/ee"
	_ "upspin.io/transports"
)

// state holds the state for a combined directory and store server that serves
// a file system for interacting with the Cloud DNS service that does DNS
// resolution for the upspin.services domain.
// See the package comment for a full description of the file system.
// A single state struct is shared by all instances of the server (see below).
type state struct {
	ep  upspin.Endpoint
	cfg upspin.Config

	accessEntry *upspin.DirEntry
	accessBytes []byte

	// cache is a cache of recently-packed entries keyed by file path.
	// It also acts as a "negative" cache for non-existent entries, in
	// which case the cached value is nil.
	cache *cache.LRU // [filePath]*entry

	dnsSvc *dns.Service
}

// server provides a wrapper around state that provies the upspin.DirServer and
// upspin.StoreServer methods for a particular user.
// A new wrapper is created by each Dial.
type server struct {
	user upspin.UserName // Set by Dial.

	// state is the state struct shared by all instances of the server.
	*state
}

// entry represents a cached DirEntry and its underlying data.
type entry struct {
	de   *upspin.DirEntry
	data []byte
}

type dirServer struct {
	*server
}

type storeServer struct {
	*server
}

func (s *server) DirServer() upspin.DirServer {
	return dirServer{s}
}

func (s *server) StoreServer() upspin.StoreServer {
	return storeServer{s}
}

const (
	accessRef        = upspin.Reference(access.AccessFile)
	packing          = upspin.EEPack
	maxCachedEntries = 1000
)

var accessRefdata = upspin.Refdata{Reference: accessRef}

func newServer(ep upspin.Endpoint, cfg upspin.Config) (*server, error) {
	s := &server{
		state: &state{
			ep:    ep,
			cfg:   cfg,
			cache: cache.NewLRU(maxCachedEntries),
		},
	}

	err := s.setupDNSService()
	if err != nil {
		return nil, err
	}

	// Allow anyone to write, but only the server user to read.
	const accessFile = "read, write: all\n"
	s.accessEntry, s.accessBytes, err = s.pack(access.AccessFile, []byte(accessFile))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// pack uses packer to pack the named file with the given data and returns the
// resulting DirEntry and ciphertext.
func (s *server) pack(filePath string, data []byte) (*upspin.DirEntry, []byte, error) {
	name := path.Join(upspin.PathName(s.cfg.UserName()), filePath)
	de := &upspin.DirEntry{
		Writer:     s.cfg.UserName(),
		Name:       name,
		SignedName: name,
		Packing:    packing,
		Time:       upspin.Now(),
		Sequence:   1,
	}

	packer := pack.Lookup(packing)

	bp, err := packer.Pack(s.cfg, de)
	if err != nil {
		return nil, nil, err
	}
	cipher, err := bp.Pack(data)
	if err != nil {
		return nil, nil, err
	}
	bp.SetLocation(upspin.Location{
		Endpoint:  s.ep,
		Reference: upspin.Reference(filePath),
	})
	if err := bp.Close(); err != nil {
		return nil, nil, err
	}

	// Share with all.
	packer.Share(s.cfg, []upspin.PublicKey{upspin.AllUsersKey}, []*[]byte{&de.Packdata})

	return de, cipher, nil
}

// packHost packs a file for the user containing the host and IP address.
// It then returns the entry after adding it to the cache.
func (s *server) packHost(name upspin.UserName, ip, host string) (e *entry, err error) {
	e = &entry{}
	e.de, e.data, err = s.pack(string(name), []byte(fmt.Sprintf("%s\n%s\n", ip, host)))
	if err != nil {
		return nil, err
	}
	s.cache.Add(name, e)
	return e, nil
}

// lookup returns the entry for the given user. It first consults the cache.
// If the entry is unknown, it invokes lookupName to find the user's host
// information, packs the resulting data, and updates the cache.
func (s *server) lookup(name upspin.UserName) (*entry, error) {
	if err := valid.UserName(name); err != nil {
		return nil, err
	}
	ei, ok := s.cache.Get(name)
	if ok {
		if ei == nil {
			return nil, errors.E(name, errors.NotExist)
		}
		return ei.(*entry), nil
	}

	ip, host, err := s.lookupName(name)
	if errors.Match(errors.E(errors.NotExist), err) {
		s.cache.Add(name, nil)
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return s.packHost(name, ip, host)
}

// These methods implement upspin.Service.

func (s *server) Endpoint() upspin.Endpoint { return s.ep }
func (*server) Close()                      {}

// These methods implement upspin.Dialer.

func (s storeServer) Dial(cfg upspin.Config, _ upspin.Endpoint) (upspin.Service, error) {
	return (&server{
		user:  cfg.UserName(),
		state: s.state,
	}).StoreServer(), nil
}

func (s dirServer) Dial(cfg upspin.Config, _ upspin.Endpoint) (upspin.Service, error) {
	return (&server{
		user:  cfg.UserName(),
		state: s.state,
	}).DirServer(), nil
}

// These methods implement upspin.DirServer.

func (s dirServer) Lookup(name upspin.PathName) (*upspin.DirEntry, error) {
	if name == s.accessEntry.Name {
		return s.accessEntry, nil
	}

	p, err := path.Parse(name)
	if err != nil {
		return nil, err
	}
	if p.User() != s.cfg.UserName() {
		return nil, errors.E(name, errors.NotExist)
	}
	if p.FilePath() == "" {
		return &upspin.DirEntry{
			SignedName: name,
			Attr:       upspin.AttrDirectory,
			Packing:    upspin.PlainPack,
			Time:       upspin.Now(),
			Name:       name,
			Writer:     s.cfg.UserName(),
		}, nil
	}
	e, err := s.lookup(upspin.UserName(p.FilePath()))
	if err != nil {
		return nil, errors.E(name, err)
	}
	return e.de, nil
}

func (s dirServer) Glob(pattern string) ([]*upspin.DirEntry, error) {
	// Nobody has list access.
	return nil, errors.E(errors.Private)
}

func (s dirServer) Put(de *upspin.DirEntry) (*upspin.DirEntry, error) {
	if err := valid.DirEntry(de); err != nil {
		return nil, err
	}

	p, err := path.Parse(de.Name)
	if err != nil {
		return nil, err
	}
	if p.User() != s.cfg.UserName() {
		return nil, errors.E(de.Name, errors.NotExist)
	}
	parts := strings.Split(p.FilePath(), "/")
	if len(parts) != 2 {
		return nil, errors.E(errors.Permission, de.Name, errors.Str("file names must be of the form user@example.com/ip"))
	}
	user, ip := upspin.UserName(parts[0]), parts[1]
	if user != s.user {
		return nil, errors.E(errors.Permission, de.Name)
	}

	host, err := s.updateName(user, ip)
	if err != nil {
		return nil, err
	}
	_, err = s.packHost(user, ip, host)
	return nil, err
}

func (s dirServer) WhichAccess(name upspin.PathName) (*upspin.DirEntry, error) {
	return s.accessEntry, nil
}

// This method implements upspin.StoreServer.

func (s storeServer) Get(ref upspin.Reference) ([]byte, *upspin.Refdata, []upspin.Location, error) {
	if ref == accessRef {
		return s.accessBytes, &accessRefdata, nil, nil
	}
	e, err := s.lookup(upspin.UserName(ref))
	if err != nil {
		return nil, nil, nil, err
	}
	return e.data, &upspin.Refdata{Reference: ref, Volatile: true}, nil, nil
}

// The DirServer and StoreServer methods below are not implemented.

var errNotImplemented = errors.E(errors.Invalid, errors.Str("method not implemented"))

func (dirServer) Delete(name upspin.PathName) (*upspin.DirEntry, error) {
	return nil, errNotImplemented
}

func (dirServer) Watch(_ upspin.PathName, _ int64, _ <-chan struct{}) (<-chan upspin.Event, error) {
	return nil, upspin.ErrNotSupported
}

func (storeServer) Put(data []byte) (*upspin.Refdata, error) {
	return nil, errNotImplemented
}

func (storeServer) Delete(ref upspin.Reference) error {
	return errNotImplemented
}
