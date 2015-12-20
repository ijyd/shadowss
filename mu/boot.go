package main

import (
	//"encoding/binary"
	//"errors"
	// "flag"
	// "fmt"
	// "github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	//"io"
	//"net"
	//"os"
	"github.com/orvice/shadowsocks-go/mu/log"
	//"os/signal"
	// "runtime"
	//"strconv"
	//"sync"
	//"syscall"
	"strconv"
	"time"
)

var configFile string
var config *ss.Config

func boot() {
	var err error

	// log.SetOutput(os.Stdout)

	err = InitMySqlClient()
	if err != nil {
		log.Log.Panic(err)
	}
	users, err := Client.GetUsers()
	if err != nil {
		log.Log.Panic(err)
	}
	log.Log.Info(len(users))
	for _, user := range users {
		log.Log.Info(user)
		port := strconv.Itoa(user.GetPort())
		password := user.GetPasswd()
		log.Log.Info(port)
		go run(port, password)
		time.Sleep(30)
	}

	// waitSignal()
}
