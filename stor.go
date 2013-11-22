// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"reflect"
	"runtime/debug"
	"sync"
	"time"
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
func (self *Store) Has(m *Msg) *Msg {
	self.RLock()
	_, ok := self.data[m.Key]
	self.RUnlock()
	m, m.Val = &Msg{}, ok
	return m
}

// store set method
func (self *Store) Set(m *Msg) *Msg {
	self.Lock()
	self.data[m.Key] = m.Val
	_, ok := self.data[m.Key]
	self.Unlock()
	m, m.Val = &Msg{}, ok
	return m
}

// store get method
func (self *Store) Get(m *Msg) *Msg {
	self.RLock()
	v := Match(self.data, MatchCleaner(m.Val), 0)
	self.RUnlock()
	m, m.Val = &Msg{}, v
	return m
}

// store del method
func (self *Store) Del(m *Msg) *Msg {
	self.Lock()
	delete(self.data, m.Key)
	_, ok := self.data[m.Key]
	self.Unlock()
	m, m.Val = &Msg{}, ok
	return m
}

// store query method
// requires map[string]interface{} (aka json object) as a query term
func (self *Store) Qry(m *Msg) *Msg {
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
			m, m.Val = &Msg{}, nil
		} else {
			m, m.Val = &Msg{}, set
		}
	default:
		m, m.Val = &Msg{}, nil
	}
	return m
}

// store dsk (save/load) calls
func (self *Store) Dsk(m *Msg) *Msg {
	switch m.Key {
	case "save":
		fd, err := os.Create("data.json")
		if err != nil {
			m, m.Val = &Msg{}, false
			return m
		}
		self.RLock()
		json.NewEncoder(fd).Encode(self.data)
		self.RUnlock()
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
		debug.FreeOSMemory()
		m, m.Val = &Msg{}, true
		return m
	case "load":
		fd, err := os.Open("data.json")
		if err != nil {
			m, m.Val = &Msg{}, false
			return m
		}
		self.Lock()
		self.data = map[string]interface{}{}
		json.NewDecoder(fd).Decode(&self.data)
		self.Unlock()
		if err := fd.Close(); err != nil {
			log.Fatal(err)
		}
		debug.FreeOSMemory()
		m, m.Val = &Msg{}, true
		return m
	default:
		m, m.Val = &Msg{}, false
		return m
	}
}

// handle "micro-threaded" connection
func (self *Store) HandleConn(conn net.Conn) {
	dec, enc := json.NewDecoder(conn), json.NewEncoder(conn)
	for {
		var m Msg
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			conn.Close()
			return
		} else {
			conn.SetDeadline(time.Now().Add(time.Minute * 15))
		}
		switch m.Cmd {
		case "has":
			enc.Encode(self.Has(&m))
		case "set":
			enc.Encode(self.Set(&m))
		case "get":
			enc.Encode(self.Get(&m))
		case "del":
			enc.Encode(self.Del(&m))
		case "qry":
			enc.Encode(self.Qry(&m))
		case "dsk":
			enc.Encode(self.Dsk(&m))
		case "bye":
			conn.SetDeadline(time.Now())
		default:
			enc.Encode(Msg{Val: false})
		}
	}
}
