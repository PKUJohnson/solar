package std

import (
	"testing"
)

func TestError(t *testing.T) {
	err := &Err{1, "error message"}
	err = ErrFromString(err.Error())
	t.Log(err)
}
