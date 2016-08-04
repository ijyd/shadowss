package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ClientConfig manage client configure file
type ClientConfig struct {
	ListenPort int              `json:"listenPort"`
	Server     []ConnectionInfo `json:"servers"`
}

// NewClientConfig new a ClientConfig handler
func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		ListenPort: 1080,
		Server:     make([]ConnectionInfo, 1, 10),
	}
}

// Parse input a config file for parse
func (c *ClientConfig) Parse(file string) error {
	fileHandle, err := os.Open(file) // For read access.
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	data, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return err
	}

	config := &ClientConfig{}
	if err = json.Unmarshal(data, config); err != nil {
		return err
	}

	return nil
}
