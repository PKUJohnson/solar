package toolkit

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	std "github.com/PKUJohnson/solar/std"
)

func Death(will func()) {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	std.LogErrorLn("DEATH triggered, die in 1 second !!!")
	if will != nil {
		will()
	}
	time.Sleep(time.Second * 1)
}

func SyncMapToMap(m *sync.Map) map[interface{}]interface{} {
	ret := make(map[interface{}]interface{})
	m.Range(func(key, value interface{}) bool {
		ret[key] = value
		return true
	})
	return ret
}

func MapToSyncMap(m map[interface{}]interface{}) *sync.Map {
	ret := new(sync.Map)
	for k, v := range m {
		ret.Store(k, v)
	}
	return ret
}
