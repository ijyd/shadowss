package config

var (
	maxTCPConnectionPerPort = 300

	defautAPIProxyListenPort = 48888
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
