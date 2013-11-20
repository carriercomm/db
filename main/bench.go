// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *

package main

import "db"
import "fmt"

// main benchmarker implementation
func main() {
	db.InitBench(fmt.Sprintf("localhost:%d", *db.LISTEN)).Run(10)
}
