// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * NARDb :: Not A Relational Database
// *

package main

import (
    "encoding/json"
    "fmt"
    "log"
    "db"
    "net"
    "sync"
    "time"
)

const (
    COUNT = 100000
    VAL_1 = `{"cmd":"set","key":"%d","val":{"name":"greg","age":25,"orders":["one","two","three"],"married":true,"inner":{"name2":"scott"}}}`
    VAL_2 = `{"cmd":"set","key":"%d","val":["foo","bar"]}`
)

type Bench struct {
    name         string
    start, count int
    conn         net.Conn
}

func InitBench(name, host string, start, count int) *Bench {
    conn, err := net.Dial("tcp", host)
    if err != nil {
        log.Fatal(err)
    }
    return &Bench{
        name:  name,
        conn:  conn,
        start: start,
        count: count,
    }
}

func (self *Bench) Run(wg *sync.WaitGroup) {
    dec, enc := json.NewDecoder(self.conn), json.NewEncoder(self.conn)
    for i := self.start; i < self.start+self.count; i++ {
        //self.conn.Write([]byte(fmt.Sprintf(VAL_2, i)))
        enc.Encode(db.Msg{Cmd: "set", Key: fmt.Sprintf("%d", i), Val: "bar"})
        var v interface{}
        dec.Decode(&v)
        if v == nil {
            log.Panicln("Empty response!")
        }
    }
    enc.Encode(db.Msg{Cmd: "bye"})
    self.conn.Close()
    wg.Done()
}

func main() {

    // heading
    fmt.Println("Benchmarking server...")

    // initialize 10 benchmark clients
    // each client sending 100,000 requests
    b0 := InitBench("b0", "localhost:1234", 0*COUNT, COUNT)
    b1 := InitBench("b1", "localhost:1234", 1*COUNT, COUNT)
    b2 := InitBench("b2", "localhost:1234", 2*COUNT, COUNT)
    b3 := InitBench("b3", "localhost:1234", 3*COUNT, COUNT)
    b4 := InitBench("b4", "localhost:1234", 4*COUNT, COUNT)
    b5 := InitBench("b5", "localhost:1234", 5*COUNT, COUNT)
    b6 := InitBench("b6", "localhost:1234", 6*COUNT, COUNT)
    b7 := InitBench("b7", "localhost:1234", 7*COUNT, COUNT)
    b8 := InitBench("b8", "localhost:1234", 8*COUNT, COUNT)
    b9 := InitBench("b9", "localhost:1234", 9*COUNT, COUNT)
    b10 := InitBench("b10", "localhost:1234", 10*COUNT, COUNT)
    b11 := InitBench("b11", "localhost:1234", 11*COUNT, COUNT)
    b12 := InitBench("b12", "localhost:1234", 12*COUNT, COUNT)
    b13 := InitBench("b13", "localhost:1234", 13*COUNT, COUNT)
    b14 := InitBench("b14", "localhost:1234", 14*COUNT, COUNT)
    b15 := InitBench("b15", "localhost:1234", 15*COUNT, COUNT)
    b16 := InitBench("b16", "localhost:1234", 16*COUNT, COUNT)
    b17 := InitBench("b17", "localhost:1234", 17*COUNT, COUNT)
    b18 := InitBench("b18", "localhost:1234", 18*COUNT, COUNT)
    b19 := InitBench("b19", "localhost:1234", 19*COUNT, COUNT)

    // run each benchmark client on a sperate thread
    // 10 concurrent clients, totaling 1 million requests
    var wg sync.WaitGroup
    wg.Add(20)

    t := time.Now().Unix()
    go b0.Run(&wg)
    go b1.Run(&wg)
    go b2.Run(&wg)
    go b3.Run(&wg)
    go b4.Run(&wg)
    go b5.Run(&wg)
    go b6.Run(&wg)
    go b7.Run(&wg)
    go b8.Run(&wg)
    go b9.Run(&wg)
    go b10.Run(&wg)
    go b11.Run(&wg)
    go b12.Run(&wg)
    go b13.Run(&wg)
    go b14.Run(&wg)
    go b15.Run(&wg)
    go b16.Run(&wg)
    go b17.Run(&wg)
    go b18.Run(&wg)
    go b19.Run(&wg)
    wg.Wait()
    ts := time.Now().Unix() - t

    // server request stats
    fmt.Printf("Server took %d seconds to complete %d requests (%d/rps)\n", ts, 20*COUNT, int64(20*COUNT)/ts)
}
