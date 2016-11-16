package app

import (
	"fmt"

	"github.com/golang/glog"

	_ "gofreezer/examples/common/apiext/install"
)

func init() {
	fmt.Printf("Call install api")
	glog.Infof("Call install")
}
