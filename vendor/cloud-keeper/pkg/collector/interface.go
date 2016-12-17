package collector

import (
	"cloud-keeper/pkg/api"
	"gofreezer/pkg/pages"
)

//Account server acc information
type Account interface {
	GetAccount() (*api.AccountInfoSpec, error)
}

//Server account server operate interface
type Server interface {
	CreateServer(server interface{}) error
	DeleteServer(id int64) error
	GetServers(page pages.Selector) ([]api.AccServer, error)
	GetServer(id int) (*api.AccServer, error)
	ServerExec(serverid int, cmd string) error
	Exec(exec *api.AccExec) error
	GetSSHKey() (*api.AccSSHKey, error)
}

//Collector aggregation operate interface
type Collector interface {
	Account
	Server
}
