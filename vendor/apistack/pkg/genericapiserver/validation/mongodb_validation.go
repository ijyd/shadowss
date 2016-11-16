package validation

import (
	"apistack/pkg/genericapiserver/options"

	"github.com/golang/glog"
)

func VerifyMongodbConfig(options *options.ServerRunOptions) {
	if len(options.StorageConfig.Mongodb.ServerList) == 0 {
		glog.Fatalf("--mongo-servers must be specified")
	}
}
