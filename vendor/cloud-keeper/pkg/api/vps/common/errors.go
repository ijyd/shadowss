package common

import (
	apierr "cloud-keeper/pkg/api/errors"

	"github.com/golang/glog"
)

func EncodeError(status interface{}) []byte {
	var output []byte
	internalErr, ok := status.(*apierr.StatusError)
	if ok {
		output = internalErr.ErrStatus.Encode()
	} else {
		glog.Errorln("status type error")
	}

	return output
}

// func EncodeError(status interface{}) []byte {
// 	return EncodeError(status)
// }

func IsNotfoundErr(err error) bool {
	if err.Error() == string("not found") {
		return true
	}
	return false
}
