package validation

import (
	"apistack/pkg/genericapiserver/options"

	"github.com/golang/glog"
)

func VerifyMysqlServersList(options *options.ServerRunOptions) {
	if len(options.StorageConfig.Mysql.ServerList) == 0 {
		glog.Fatalf("--mysql-servers must be specified")
	}
}
