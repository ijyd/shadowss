package main

import (
	"flag"
	"shadowsocks-go/cmd/shadowsslicense/app"

	"github.com/golang/glog"
)

func main() {

	flag.Set("alsologtostderr", "true")
	flag.Set("v", "5")
	flag.Parse()

	glog.V(5).Infoln("tset")
	app.Run()
}
