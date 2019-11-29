package strutil

import (
	"fmt"
	"testing"
)

func TestM(t *testing.T) {
	m := &columnShieldWords{
		Words:  []string{"aa", "bb", "ccc"},
		Number: 2,
	}
	fmt.Println(FromObject(m))
}

type columnShieldWords struct {
	Words  []string `json:"words"`
	Number int64    `json:"number"`
}
