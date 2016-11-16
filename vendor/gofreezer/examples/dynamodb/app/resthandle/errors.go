package resthandle

import (
	apierr "gofreezer/examples/common/apiext/errors"

	"github.com/golang/glog"
)

func encodeError(status interface{}) []byte {
	var output []byte
	internalErr, ok := status.(*apierr.StatusError)
	if ok {
		output = internalErr.ErrStatus.Encode()
	} else {
		glog.Errorln("status type error")
	}

	return output
}

func isNotfoundErr(err error) bool {
	if err.Error() == string("not found") {
		return true
	}
	return false
}
