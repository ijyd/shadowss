package config

import (
	"fmt"
	"net/http"
)

var (
	maxTCPConnectionPerPort = 300

	defautAPIProxyListenPort = 12345
	defaultAPIProxyPassword  = "123456790"

	token = "Bearer 455151fsfjkkdakllds1111a"
)

//SetMaxTCPConnectionPerPort set max tcp connection for per port
func SetMaxTCPConnectionPerPort(limit int) {
	maxTCPConnectionPerPort = limit
}

//GetMaxTCPConnectionPerPort set max tcp connection for per port
func GetMaxTCPConnectionPerPort() int {
	return maxTCPConnectionPerPort
}

//SetDefaultAPIProxyListenPort set default api proxy listen port
func SetDefaultAPIProxyListenPort(port int) {
	defautAPIProxyListenPort = port
}

//GetDefaultAPIProxyListenPort get default api proxy listen port
func GetDefaultAPIProxyListenPort() int {
	return defautAPIProxyListenPort
}

//SetDefaultAPIProxyPassword set default api proxy listen port
func SetDefaultAPIProxyPassword(pwd string) {
	defaultAPIProxyPassword = pwd
}

//GetDefaultAPIProxyPassword get default api proxy password
func GetDefaultAPIProxyPassword() string {
	return defaultAPIProxyPassword
}

//AddAuthHTTPHeader add Authorization header into http request
func AddAuthHTTPHeader(req *http.Request) {
	req.Header.Add("Authorization", token)
}

//SetToken get common token
func SetToken(tk string) {
	token = fmt.Sprintf("Bearer %s", tk)
}

//GetToken get common token
func GetToken() string {
	return token
}
