package main

import (
	"github.com/orvice/shadowsocks-go/mu/log"
	"os"
)

func main() {
	var err error

	err = InitConfig()
	if err != nil {
		log.Log.Error(err)
		os.Exit(0)
	}

	err = InitRedis()
	if err != nil {
		log.Log.Error("boot redis fail: ", err)
		os.Exit(0)
	}
	boot()
}
