package main

import (
	"flag"
)

var (
	debug   bool
	logPath string
)

func InitFlag() {
	// flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	flag.StringVar(&logPath, "log_path", "./ss.log", "log file path")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()
}
