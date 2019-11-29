package service

import (
	"log"
	"net/http"
	_ "net/http/pprof" // import the pprof package
	"os"
)

func StartPprof() {
	if os.Getenv("OPEN_PPROF") == "1" {
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()
	}
}
