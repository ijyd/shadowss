package collector

import "cloud-keeper/pkg/api"

//Account server acc information
type Account interface {
	GetAccount() (*api.AccountInfoSpec, error)
}

//Server account server operate interface
type Server interface {
	CreateServer(server interface{}) error
	DeleteServer(id int64) error
	//GetServers(page pagination.Pager) ([]byte, error)
	GetServer(id int) (*api.AccServerSpec, error)
	ServerExec(serverid int, cmd string) error
	Exec(exec *api.AccExec) error
	GetSSHKey() (*api.AccSSHKey, error)
}

//Collector aggregation operate interface
type Collector interface {
	Account
	Server
}
