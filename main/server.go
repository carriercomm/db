// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *

package main

import "db"

// main server implementation
func main() {
	db.InitServer(*db.BACKUP).ListenAndServe(*db.LISTEN)
}
