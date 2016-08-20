package config

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"os"

	"github.com/golang/glog"
)

// ServerConfig for server configure
type ServerConfig struct {
	Clients []ConnectionInfo `json:"clients"`
}

const (
	generateRandID = 100000000
)

// ServerCfg server configure handle
var ServerCfg = NewServerConfig()

// NewServerConfig new a ClientConfig handler
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Clients: make([]ConnectionInfo, 1, 30),
	}
}

func (s *ServerConfig) verifyConfig() error {
	for idx, v := range s.Clients {
		glog.Infof("verify config %v", v)
		if v.ID == 0 {
			//generate default id for this clients
			id := rand.Int63n(generateRandID)
			s.Clients[idx].ID = (id + generateRandID) * 100
		}
	}

	return nil
}

// Parse input a config file for parse
func (s *ServerConfig) Parse(file string) error {
	fileHandle, err := os.Open(file) // For read access.
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	data, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return err
	}

	//config := &ServerConfig{}
	if err = json.Unmarshal(data, s); err != nil {
		return err
	}
	//copy(s, config)
	glog.V(5).Infoln("Got Configure clients:%+v", s)

	return s.verifyConfig()
}
