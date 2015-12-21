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
	"github.com/orvice/shadowsocks-go/mu/log"
	//"os/signal"
	// "runtime"
	//"strconv"
	//"sync"
	//"syscall"
	"fmt"
	"os"
	"strconv"
)

var configFile string
var config *ss.Config

func boot() {
	var err error

	// log.SetOutput(os.Stdout)

	err = InitMySqlClient()
	if err != nil {
		log.Log.Error(err)
		os.Exit(0)
	}
	client := user.GetClient()
	users, err := client.GetUsers()
	if err != nil {
		log.Log.Error(err)
		os.Exit(0)
	}
	log.Log.Info(len(users))
	bootUsers(users)
	waitSignal()
}

// 第一次启动
func bootUsers(users []user.User) {
	for _, user := range users {
		log.Log.Info(user)
		port := strconv.Itoa(user.GetPort())
		password := user.GetPasswd()
		log.Log.Info(port)
		cipher, err := user.GetCipher()
		if err != nil {
			log.Log.Error(fmt.Sprintf("error on boot port %d,skip.", user.GetPort()), err)
			continue
		}
		go runWithCustomMethod(port, password, cipher)
	}
}

// check users
func checkUsers(users []user.User) {

}
