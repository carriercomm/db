// *
// * Copyright 2013, Scott Cagno, All rights reserved.
// * BSD Licensed - sites.google.com/site/bsdc3license
// *
// * Db :: Database
// *

package db

import (
	"flag"
	"time"
)

// cmd line options
var (
	BACKUP = flag.Duration("backup_rate", time.Minute*15, "archive interval")
	LISTEN = flag.Int("listen_port", 9000, "listening port")
)

// run flag parser
func init() {
	flag.Parse()
}
