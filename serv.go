// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// transport schema
type Msg struct {
	Cmd string      `json:"cmd,omitempty"`
	Key string      `json:"key,omitempty"`
	Val interface{} `json:"val,omitempty"`
}

// data store
type Server struct {
	st   *Store
	arch *Archiver
	rate time.Duration
}

// initialize data store
func InitServer(backup_rate time.Duration) *Server {
	self := &Server{
		st:   InitStore(),
		rate: backup_rate,
	}
	self.arch = InitArchiver(self.st)
	self.CatchSignals()
	return self
}

// catch interrupt signals
func (self *Server) CatchSignals() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSTOP)
	go func() {
		fmt.Printf("\nCaught %v\n", <-sig)
		self.arch.Save()
		os.Exit(0)
	}()
}

// listen on given port, and serve client connections
func (self *Server) ListenAndServe(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	self.arch.RunArchiver(self.rate)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go self.HandleConn(conn)
	}
}

// handle "micro-threaded" connection
func (self *Server) HandleConn(conn net.Conn) {
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
			enc.Encode(self.st.Has(&m))
		case "set":
			enc.Encode(self.st.Set(&m))
		case "get":
			enc.Encode(self.st.Get(&m))
		case "del":
			enc.Encode(self.st.Del(&m))
		case "qry":
			enc.Encode(self.st.Qry(&m))
		case "dsk":
			enc.Encode(self.st.Dsk(&m))
		case "bye":
			conn.SetDeadline(time.Now())
		default:
			enc.Encode(Msg{Val: false})
		}
	}
}
