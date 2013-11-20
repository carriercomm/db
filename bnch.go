// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *

package db

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// const counts
const (
	COUNT = 100000
)

// benchmarker struct
type Bench struct {
	host  string
	count int
	conns []net.Conn
}

// initialize and return a benchmarker
func InitBench(host string) *Bench {
	return &Bench{
		host: host,
	}
}

// run n number of benchmark clients
func (self *Bench) Run(n int) {
	for i := 0; i < n; i++ {
		self.conns = append(self.conns, dial(self.host))
	}
	var wg sync.WaitGroup
	wg.Add(len(self.conns))
	t := time.Now().Unix()
	for n, conn := range self.conns {
		go run(&wg, conn, n, COUNT)
	}
	wg.Wait()
	ts := time.Now().Unix() - t
	fmt.Printf("Server took %d seconds to complete %d requests (%d/rps)\n", ts, len(self.conns)*COUNT, int64(len(self.conns)*COUNT)/ts)
}

// dial up to the host
func dial(host string) net.Conn {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// internal single runner
func run(wg *sync.WaitGroup, conn net.Conn, start, count int) {
	dec, enc := json.NewDecoder(conn), json.NewEncoder(conn)
	for i := start; i < start+count; i++ {
		enc.Encode(Msg{Cmd: "set", Key: fmt.Sprintf("%d", i), Val: 123})
		var v interface{}
		dec.Decode(&v)
		if v == false {
			log.Panicln("Invalid command issued!")
		}
	}
	enc.Encode(Msg{Cmd: "bye"})
	conn.Close()
	wg.Done()
}
