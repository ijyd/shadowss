package collector

import "cloud-keeper/pkg/pagination"

//Account server acc information
type Account interface {
	GetAccount() ([]byte, error)
}

//Server account server operate interface
type Server interface {
	CreateServer(server interface{}) error
	DeleteServer(id int64) error
	GetServers(page pagination.Pager) ([]byte, error)
	Exec(param interface{}) error
	GetSSHKey() ([]byte, error)
}

//Collector aggregation operate interface
type Collector interface {
	Account
	Server
}
