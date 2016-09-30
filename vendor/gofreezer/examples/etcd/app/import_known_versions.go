package app

import (
	"fmt"

	"github.com/golang/glog"

	_ "gofreezer/examples/etcd/app/api/install"
)

func init() {
	fmt.Printf("Call install api")
	glog.Infof("Call install")
}
