package main

import (
	"os"
)

func main() {
	var err error

	err = InitConfig()
	if err != nil {
		Log.Error(err)
		os.Exit(0)
	}

	err = InitRedis()
	if err != nil {
		Log.Error("boot redis fail: ", err)
		os.Exit(0)
	}
	boot()
}
