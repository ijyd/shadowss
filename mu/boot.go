package main

import (
	//"encoding/binary"
	//"errors"
	// "flag"
	// "fmt"
	"github.com/orvice/shadowsocks-go/mu/user"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	//"io"
	"log"
	//"net"
	"os"
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

	log.SetOutput(os.Stdout)

	client := user.GetClient()
	users, err := client.GetUsers()
	if err != nil {

	}

	for _, user := range users {
		port := strconv.Itoa(user.GetPort())
		password := user.GetPasswd()
		go run(port, password)
	}

	waitSignal()
}
