package validation

import (
	"apistack/pkg/genericapiserver/options"

	"github.com/golang/glog"
)

func VerifyDynamodbConfig(options *options.ServerRunOptions) {
	if len(options.StorageConfig.AWSDynamoDB.Region) == 0 {
		glog.Fatalf("--aws-region must be specified")
	}
}
