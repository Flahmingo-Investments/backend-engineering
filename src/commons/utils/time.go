package commons

import (
	"flag"
	"time"
)

var TestingNow time.Time

func Now() time.Time {
	if flag.Lookup("test.v") == nil {
		return time.Now()
	} else {
		return TestingNow
	}
}
