package helper

import (
	"math/rand"
	"strconv"
)

// RandomHex generates a random hex.
func RandomHex() string {
	return strconv.FormatInt(RandomInt(), 16)
}

// RandomInt generates a random int64.
func RandomInt() int64 {
	//return rand.Int63() & 0x001fffffffffffff
	return rand.Int63()
}
