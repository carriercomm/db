// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *

package main

import (
	"db"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	REQUESTS = 100000
	WORKERS  = 10
)

func client() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		log.Fatal(err)
	}
	enc, dec := json.NewEncoder(conn), json.NewDecoder(conn)
	for i := 0; i < REQUESTS; i++ {
		enc.Encode(db.Msg{Cmd: "set", Key: fmt.Sprintf("%d", i), Val: true})
		var m db.Msg
		dec.Decode(&m)
		if m.Val == false {
			log.Fatal("woops")
		}
	}
	conn.Close()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(WORKERS)
	t := time.Now().Unix()
	for i := 0; i < WORKERS; i++ {
		go func() {
			client()
			wg.Done()
		}()
	}
	wg.Wait()
	t = time.Now().Unix() - t
	fmt.Printf("%v rps\n", float64((WORKERS*REQUESTS)/t))
}
