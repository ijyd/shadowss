package config

var (
	maxTCPConnectionPerPort = 50
)

//SetMaxTCPConnectionPerPort set max tcp connection for per port
func SetMaxTCPConnectionPerPort(limit int) {
	maxTCPConnectionPerPort = limit
}

//GetMaxTCPConnectionPerPort set max tcp connection for per port
func GetMaxTCPConnectionPerPort() int {
	return maxTCPConnectionPerPort
}
