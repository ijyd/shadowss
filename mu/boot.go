package main

import (
	//"encoding/binary"
	//"errors"
	// "flag"
	// "fmt"
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	//"io"
	//"net"
	//"os"
	log "github.com/Sirupsen/logrus"
	//"os/signal"
	// "runtime"
	//"strconv"
	//"sync"
	//"syscall"
	"strconv"
)

var configFile string
var config *ss.Config

func boot() {
	var err error

	// log.SetOutput(os.Stdout)

	InitClient()
	client := user.GetClient()
	users, err := client.GetUsers()
	if err != nil {

	}

	for _, user := range users {
		log.Info(user)
		port := strconv.Itoa(user.GetPort())
		password := user.GetPasswd()
		log.WithFields(log.Fields{
			"Port":     port,
			"Password": password,
		}).Info("Start new user")
		go run(port, password)
	}

	waitSignal()
}
