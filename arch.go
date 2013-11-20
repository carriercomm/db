// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"encoding/gob"
	"log"
	"os"
	"runtime/debug"
	"time"
)

// archive struct
type Archiver struct {
	data  map[string]interface{}
	count int
	store *Store
}

// init archive
func InitArchiver(st *Store) *Archiver {
	return &Archiver{
		data:  make(map[string]interface{}),
		count: 0,
		store: st,
	}
}

// run wrapper
func (self *Archiver) RunArchiver(rate time.Duration) {
	self.Load()
	go self.Run(rate)
}

// run method
func (self *Archiver) Run(rate time.Duration) {
	time.AfterFunc(rate, func() {
		if len(self.store.data) != self.count {
			self.Save()
		}
		self.Run(rate)
	})
}

// try to load archive
func (self *Archiver) Load() {
	fd, err := os.Open("archive.gob")
	if err != nil || len(self.store.data) > 0 {
		return
	}
	gob.NewDecoder(fd).Decode(&self.data)
	self.store.Lock()
	for k, v := range self.data {
		self.store.data[k] = v
	}
	self.store.Unlock()
	self.data = map[string]interface{}{}
	debug.FreeOSMemory()
}

// try to save archive
func (self *Archiver) Save() {
	fd, err := os.Create("archive.gob")
	if err != nil {
		log.Fatal(err)
	}
	self.store.RLock()
	for k, v := range self.store.data {
		self.data[k] = v
	}
	self.store.RUnlock()
	gob.NewEncoder(fd).Encode(&self.data)
	if err := fd.Close(); err != nil {
		log.Fatal(err)
	}
	self.data = map[string]interface{}{}
	self.count = len(self.store.data)
	debug.FreeOSMemory()
}
