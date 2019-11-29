package common

import (
	"time"
	"math/rand"
)
var rsid = rand.New(rand.NewSource(time.Now().UnixNano()))

func GetRandInt(n int) int{
	return rsid.Intn(n) % n
}

func GetRandIntRange(lower int, upper int) int {
	rint := GetRandInt(upper)
	ratio := float32(rint) / float32(upper)
	res := float32(upper - lower) * ratio + float32(lower)
	return int(res)
}
