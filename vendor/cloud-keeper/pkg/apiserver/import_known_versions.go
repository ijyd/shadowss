package apiserver

import (
	"fmt"

	_ "cloud-keeper/pkg/api/install"

	"github.com/golang/glog"
)

func init() {
	fmt.Printf("Call install api")
	glog.Infof("Call install")
}
