package vultr

import "cloud-keeper/pkg/collector/vultr/client"

const (
	key = "MER7QU6JAKTRAFCB4LNUPYDTZF7CXLLBPMHA"
)

//Vultr implements vultr vps interface
type Vultr struct {
	*client.Client
}

//NewVultr create a vultr handler
func NewVultr(key string) *Vultr {
	return &Vultr{
		client.NewClient(key),
	}
}

//GetAccount get account  information
// func (vul *Vultr) GetAccount() (*api.AccountInfoSpec, error) {
// 	return vul.client.GetAccount()
// }
//
// //CreateServer create server
// func (vul *Vultr) CreateServer(server interface{}) error {
// 	return vul.client.CreateServer(server)
// }
//
// //DeleteServer delete server by id
// func (vul *Vultr) DeleteServer(id int64) error {
// 	return vul.client.DeleteServer(id)
// }
//
// //GetServers get account server information
// func (vul *Vultr) GetServers(page pagination.Pager) ([]byte, error) {
// 	return vul.client.GetServers(page)
// }
//
// func (vul *Vultr) Exec(param *api.AccExec) error {
// 	return vul.client.Exec(param)
// }
//
// func (vul *Vultr) GetSSHKey() (*api.AccSSHKey, error) {
// 	return vul.client.GetSSHKey()
// }
//
// func (vul *Vultr) ServerExec(serverid int, cmd string) error {
// 	return vul.client.ServerExec(serverid, cmd)
// }
//
// func (vul *Vultr) GetServer(int id) (*api.AccServerSpec, error) {
//
// }
