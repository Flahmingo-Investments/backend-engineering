package commons

import (
	"math/rand"
	"time"
)

func RandOTP() string {
	rand.Seed(time.Now().UnixNano())
	source := "012345789"
	digits := 4
	b := make([]byte, digits)
	for i := range b {
		b[i] = source[rand.Intn(len(source))]
	}
	return string(b)
}
