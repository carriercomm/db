// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"fmt"
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
		go self.st.HandleConn(conn)
	}
}
