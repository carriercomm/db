// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"encoding/json"
	"net"
	"sync"
)

// client struct
type Client struct {
	conn net.Conn
	enc  *json.Encoder
	dec  *json.Decoder
	host string
	sync.RWMutex
}

// initialize a client struct
func InitClient(host string) *Client {
	return &Client{
		host: host,
	}
}

// open a connection to the server and initialize the json decoder/encoder
func (self *Client) Open() bool {
	if conn, err := net.Dial("tcp", self.host); err == nil {
		self.enc, self.dec = json.NewEncoder(conn), json.NewDecoder(conn)
		self.conn = conn
		return true
	}
	return false
}

// close the connection to the server
func (self *Client) Close() bool {
	if err := self.conn.Close(); err != nil {
		return false
	}
	return true
}

// send request to the server
func (self *Client) Call(cmd, key string, val interface{}) interface{} {
	self.enc.Encode(Msg{Cmd: cmd, Key: key, Val: val})
	var m Msg
	self.dec.Decode(&m)
	return m.Val
}

// has call wrapper
func (self *Client) Has(key string) bool {
	return self.Call("has", key, nil).(bool)
}

// set call wrapper
func (self *Client) Set(key string, val interface{}) bool {
	return self.Call("set", key, val).(bool)
}

// get call wrapper
func (self *Client) Get(val ...string) interface{} {
	return self.Call("get", "", val)
}

// del call wrapper
func (self *Client) Del(key string) bool {
	return self.Call("get", key, nil).(bool)
}

// qry call wrapper
func (self *Client) Qry(val interface{}) interface{} {
	return self.Call("qry", "", val)
}

// dsk call wrapper
func (self *Client) Dsk(key string) bool {
	return self.Call("dsk", key, nil).(bool)
}
