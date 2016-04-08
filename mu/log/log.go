package log

import (
	"github.com/Sirupsen/logrus"
)

var (
	Log = logrus.New()
)

func SetLogClient(client *logrus.Logger) {
	Log = client
}
