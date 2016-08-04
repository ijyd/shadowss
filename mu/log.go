package main

import (
	"io"
	"os"

	"github.com/Sirupsen/logrus"
	"shadowsocks/shadowsocks-go/mu/log"
)

var Log = logrus.New()

func InitLog() {
	// Log.Formatter = new(logrus.JSONFormatter)
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		writer := io.MultiWriter(os.Stdout, f)
		Log.Out = writer
	}
	if debug {
		Log.Level = logrus.DebugLevel
		Log.Debug("debug on")
	}
	log.SetLogClient(Log)
}
