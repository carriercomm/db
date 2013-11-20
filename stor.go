// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
)

// data store
type Store struct {
	data map[string]interface{}
	sync.RWMutex
}

// initialize store
func InitStore() *Store {
	return &Store{
		data: make(map[string]interface{}),
	}
}

// store has method
func (self *Store) Has(m *Msg) Msg {
	self.RLock()
	_, ok := self.data[m.Key]
	self.RUnlock()
	return Msg{Val: ok}
}

// store set method
func (self *Store) Set(m *Msg) Msg {
	self.Lock()
	self.data[m.Key] = m.Val
	_, ok := self.data[m.Key]
	self.Unlock()
	return Msg{Val: ok}
}

// store get method
func (self *Store) Get(m *Msg) Msg {
	self.RLock()
	v := Match(self.data, MatchCleaner(m.Val), 0)
	self.RUnlock()
	return Msg{Val: v}
}

// store del method
func (self *Store) Del(m *Msg) Msg {
	self.Lock()
	delete(self.data, m.Key)
	_, ok := self.data[m.Key]
	self.Unlock()
	return Msg{Val: !ok}
}

// store query method
// requires map[string]interface{} (aka json object) as a query term
func (self *Store) Qry(m *Msg) Msg {
	var nm Msg
	switch m.Val.(type) {
	case map[string]interface{}:
		var set []interface{}
		for _, im := range self.data {
			switch im.(type) {
			case map[string]interface{}:
				ok := true
				for k, v := range m.Val.(map[string]interface{}) {
					switch im.(map[string]interface{})[k].(type) {
					case []interface{}:
						if !reflect.DeepEqual(im.(map[string]interface{})[k], v) {
							ok = false
							break
						}
					case map[string]interface{}:
						if !reflect.DeepEqual(im.(map[string]interface{})[k], v) {
							ok = false
							break
						}
					default:
						if im.(map[string]interface{})[k] != v {
							ok = false
							break
						}
					}
				}
				if ok {
					set = append(set, im.(map[string]interface{}))
				}
			}
		}
		if len(set) < 1 {
			nm.Val = nil
		} else {
			nm.Val = set
		}
	default:
		nm.Val = nil
	}
	return nm
}

// store dsk (save/load) calls
func (self *Store) Dsk(m *Msg) Msg {
	switch m.Key {
	case "save":
		fd, err := os.Create("data.json")
		if err != nil {
			return Msg{Val: false}
		}
		self.RLock()
		json.NewEncoder(fd).Encode(self.data)
		self.RUnlock()
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
		debug.FreeOSMemory()
		return Msg{Val: true}
	case "load":
		fd, err := os.Open("data.json")
		if err != nil {
			return Msg{Val: false}
		}
		self.Lock()
		self.data = map[string]interface{}{}
		json.NewDecoder(fd).Decode(&self.data)
		self.Unlock()
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
		debug.FreeOSMemory()
		return Msg{Val: true}
	default:
		return Msg{Val: false}
	}
}
