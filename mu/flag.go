package main

import (
	"flag"
)

var (
	debug      bool
	logPath    string
	configPath string
)

func InitFlag() {
	// flag.IntVar(&flagvar, "flagname", 1234, "help message for flagname")
	flag.StringVar(&logPath, "log_path", "./ss.log", "log file path")
	flag.StringVar(&configPath, "config_path", "./config.conf", "log file path")
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.Parse()
}
