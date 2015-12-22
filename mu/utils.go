package main

import (
	"fmt"
	"github.com/orvice/shadowsocks-go/mu/user"
)

func genUserInfoKey(user user.UserInfo) string {
	return fmt.Sprintf("userinfo:%v", user.Port)
}

func genUserFlowKey(user user.UserInfo) string {
	return fmt.Sprintf("userflow:%v", user.Port)
}
