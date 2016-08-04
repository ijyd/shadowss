package config

// ConnectionInfo description connection base information
type ConnectionInfo struct {
	Host          string `json:"host"`
	Port          int    `json:"port"`
	EncryptMethod string `json:"encrypt"`
	Password      string `json:"password"`
	EnableOTA     bool   `json:"enableOTA"`
	Timeout       int    `json:"timeout"`
}

// Config as a interface for configure file implement
type Config interface {
	Parse(file string) error
}
