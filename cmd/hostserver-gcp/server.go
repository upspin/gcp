// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"time"

	dns "google.golang.org/api/dns/v1"

	"upspin.io/access"
	"upspin.io/cache"
	"upspin.io/client/clientutil"
	"upspin.io/errors"
	"upspin.io/pack"
	"upspin.io/path"
	"upspin.io/upspin"
	"upspin.io/valid"

	_ "upspin.io/pack/ee"
	_ "upspin.io/transports"
)

// server provides implementations of upspin.DirServer and upspin.StoreServer
// ...
type server struct {
	user upspin.UserName // Set by Dial.

	*state
}

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

	if err := s.setupDNSService(); err != nil {
		return nil, err
	}

	// Allow anyone to write, but only the server user to read.
	accessFile := "read, write: all\n"

	var err error
	s.accessEntry, s.accessBytes, err = s.pack(access.AccessFile, []byte(accessFile))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *server) pack(filePath string, data []byte) (*upspin.DirEntry, []byte, error) {
	name := upspin.PathName(s.cfg.UserName()) + "/" + upspin.PathName(filePath)
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

func (s *server) packHost(name upspin.UserName, ip, host string) (e *entry, err error) {
	e = &entry{}
	e.de, e.data, err = s.pack(string(name), []byte(fmt.Sprintf("%s\n%s\n", ip, host)))
	if err != nil {
		return nil, err
	}
	s.cache.Add(name, e)
	return e, nil
}

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
	if err != nil {
		return nil, err
	}
	return s.packHost(name, ip, host)
}

// These methods implement upspin.Service.

func (s *server) Endpoint() upspin.Endpoint { return s.ep }
func (*server) Ping() bool                  { return true }
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
	e, err := s.lookup(upspin.UserName(p.FilePath()))
	if err != nil {
		return nil, err
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
	if upspin.UserName(p.FilePath()) != s.user {
		return nil, errors.E(errors.Permission, de.Name)
	}

	// Try to read the file a few times, because it could be stuck in the
	// user's cacheserver.
	var b []byte
	for attempt := 0; attempt < 4; attempt++ {
		b, err = clientutil.ReadAll(s.cfg, de)
		if errors.Match(errors.E(errors.NotExist), err) {
			time.Sleep(5 * time.Second)
			continue
		}
		if err != nil {
			return nil, err
		}
	}
	if len(b) == 0 {
		return nil, errors.E(errors.NotExist, de.Name, errors.Str("could not read data for entry"))
	}
	ip := string(bytes.TrimSpace(b))

	host, err := s.updateName(s.user, ip)
	if err != nil {
		return nil, err
	}
	e, err := s.packHost(s.user, ip, host)
	if err != nil {
		return nil, err
	}
	return e.de, nil
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
