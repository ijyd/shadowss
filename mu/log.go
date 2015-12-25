package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"io"
	"os"
)

var (
	Log     = logrus.New()
	logPath = flag.String("log_path", "./ss.log", "log file path")
	debug   = flag.Bool("debug", false, "debug")
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)

	// Log.Formatter = new(logrus.JSONFormatter)
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		writer := io.MultiWriter(os.Stdout, f)
		Log.Out = writer
	}
	if *debug {
		Log.Level = logrus.DebugLevel
	}

}
