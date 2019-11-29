package toolkit

import (
	"reflect"

	std "github.com/PKUJohnson/solar/std"
)

type StdGoRoutine interface {
	Start()
}

func GO(goroutine StdGoRoutine) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				value := reflect.ValueOf(goroutine)
				std.LogErrorLn("GO Error..xxxxxxxxxxx..FBI WARRING XXXXXX:", r, value)
			}
		}()
		goroutine.Start()
	}()
}
