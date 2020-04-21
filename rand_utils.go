package request

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandInt returns a pseudo-random number in [0,max)
func RandInt(max int) int {
	return rand.Intn(max)
}
