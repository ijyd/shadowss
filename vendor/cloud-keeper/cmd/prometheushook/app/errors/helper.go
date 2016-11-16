package errors

import (
	"github.com/golang/glog"
)

func EncodeError(status interface{}) []byte {
	var output []byte
	internalErr, ok := status.(*StatusError)
	if ok {
		output = internalErr.ErrStatus.Encode()
	} else {
		glog.Errorln("status type error")
	}

	return output
}

func IsNotfoundErr(err error) bool {
	if err.Error() == string("not found") {
		return true
	}
	return false
}
