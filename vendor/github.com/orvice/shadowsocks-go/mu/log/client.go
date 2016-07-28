package log

import (
	"github.com/orvice/shadowsocks-go/mu/user"
)

type Client interface {
	Info(user user.User, args ...interface{})
	Error(user user.User, args ...interface{})
	Debug(user user.User, args ...interface{})
}
