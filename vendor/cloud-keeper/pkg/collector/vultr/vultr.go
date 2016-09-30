package vultr

import (
	"cloud-keeper/pkg/collector/vultr/client"
	"cloud-keeper/pkg/pagination"
)

const (
	key = "MER7QU6JAKTRAFCB4LNUPYDTZF7CXLLBPMHA"
)

//Vultr implements vultr vps interface
type Vultr struct {
	client *client.Client
}

//NewVultr create a vultr handler
func NewVultr(key string) *Vultr {
	return &Vultr{
		client: client.NewClient(key),
	}
}

//GetAccount get account  information
func (vul *Vultr) GetAccount() ([]byte, error) {
	return vul.client.GetAccount()
}

//CreateServer create server
func (vul *Vultr) CreateServer(server interface{}) error {
	return vul.client.CreateServer(server)
}

//DeleteServer delete server by id
func (vul *Vultr) DeleteServer(id int64) error {
	return vul.client.DeleteServer(id)
}

//GetServers get account server information
func (vul *Vultr) GetServers(page pagination.Pager) ([]byte, error) {
	return vul.client.GetServers(page)
}

func (vul *Vultr) Exec(param interface{}) error {
	return vul.client.Exec(param)
}

func (vul *Vultr) GetSSHKey() ([]byte, error) {
	return vul.client.GetSSHKey()
}
