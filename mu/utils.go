package main

import (
	"fmt"
	"github.com/orvice/shadowsocks-go/mu/user"
)

func genUserInfoKey(user user.UserInfo) string {
	return fmt.Sprintf("userinfo:%v", user.Port)
}
