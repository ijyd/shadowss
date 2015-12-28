package main

import (
	"os"
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
	boot()
}
