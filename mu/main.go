package main

import (
	"os"
	_ "net/http/pprof"
	"log"
	"net/http"
)

func main() {
	var err error
	InitFlag()
	err = InitConfig()
	if err != nil {
		Log.Error(err)
		os.Exit(0)
	}
	InitLog()
	err = InitRedis()
	if err != nil {
		Log.Error("boot redis fail: ", err)
		os.Exit(0)
	}

	if debug {
		go func() {
			log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
		}()
	}

	boot()
	waitSignal()
}
