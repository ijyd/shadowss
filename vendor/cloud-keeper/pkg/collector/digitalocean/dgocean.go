package digitalocean

import "cloud-keeper/pkg/collector/digitalocean/client"

const (
	key = "3584c3b9ae910c10bca2d93af64d2a6897a6f20e5b0c54215ee6a5d154c76a3e"
)

//DigitalOcean implements DigitalOcean vps interface
type DigitalOcean struct {
	*client.Client
}

//NewVultr create a vultr handler
func NewDigitalOcean(key string) *DigitalOcean {
	return &DigitalOcean{
		client.NewClient(key),
	}
}

//GetAccount get account  information
// func (dgoc *DigitalOcean) GetAccount() (*api.AccountInfoSpec, error) {
// 	return dgoc.client.GetAccount()
// }
//
// //CreateServer create server
// func (dgoc *DigitalOcean) CreateServer(server interface{}) error {
// 	return dgoc.client.CreateServer(server)
// }
//
// //DeleteServer delete server by id
// func (dgoc *DigitalOcean) DeleteServer(id int64) error {
// 	return dgoc.client.DeleteServer(id)
// }
//
// //GetServers get account server information
// func (dgoc *DigitalOcean) GetServers(page pagination.Pager) ([]byte, error) {
// 	return dgoc.client.GetServers(page)
// }
//
// func (dgoc *DigitalOcean) Exec(param *api.AccExec) error {
// 	return dgoc.client.Exec(param)
// }
//
// func (dgoc *DigitalOcean) GetSSHKey() (*api.AccSSHKey, error) {
// 	return dgoc.client.GetSSHKey()
// }
//
// func (dgoc *DigitalOcean) ServerExec(serverid int, cmd string) error {
// 	return dgoc.client.ServerExec(serverid, cmd)
// }
