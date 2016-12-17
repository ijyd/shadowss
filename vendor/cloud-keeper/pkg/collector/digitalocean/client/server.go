package client

import (
	"fmt"
	"gofreezer/pkg/pages"
	"net/url"
	"strconv"

	"cloud-keeper/pkg/ansible"
	api "cloud-keeper/pkg/api"
	"cloud-keeper/pkg/collector/common"

	"github.com/digitalocean/godo"
)

func pageForURL(urlText string) (int, error) {
	u, err := url.ParseRequestURI(urlText)
	if err != nil {
		return 0, err
	}

	pageStr := u.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}

	return page, nil
}

func getItemsCount(p *godo.Pages, requirePerPage, actualPageNum int) (int, error) {

	if p != nil {
		firstPage, err := pageForURL(p.First)
		if err != nil {
			return 0, err
		}

		LastPage, err := pageForURL(p.Last)
		if err != nil {
			return 0, err
		}

		if firstPage == LastPage {

		} else {
			pages := LastPage - firstPage
			return pages * requirePerPage, nil
		}

		return firstPage + 1, nil
	}

	return 1 * actualPageNum, nil
}

func (c *Client) newAccserverSpec(droplet *godo.Droplet) *api.AccServer {
	accsrv := &api.AccServer{}

	accsrv.Spec.DigitalOcean = api.DGServerInfo{
		Location:    droplet.Region.Slug,
		CreatedTime: droplet.Created,
		Status:      droplet.Status,
		Name:        droplet.Name,

		IPV4Addr:     droplet.Networks.V4[0].IPAddress,
		IPV4NetMask:  droplet.Networks.V4[0].Netmask,
		IPV4Gateway:  droplet.Networks.V4[0].Gateway,
		PriceMonthly: droplet.Size.PriceMonthly,
		PriceHourly:  droplet.Size.PriceHourly,
	}
	accsrv.Spec.ID = fmt.Sprintf("%d", droplet.ID)
	accsrv.Spec.Size = fmt.Sprintf("%d", droplet.Memory)

	return accsrv
}

func (c *Client) GetServer(id int) (*api.AccServer, error) {
	droplet, _, err := c.client.Droplets.Get(id)
	if err != nil {
		return nil, err
	}

	accsrv := c.newAccserverSpec(droplet)

	return accsrv, nil
}

//GetServers get account server information
func (c *Client) GetServers(page pages.Selector) ([]api.AccServer, error) {

	//first get all
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	droplets, _, err := c.client.Droplets.List(opt)
	if err != nil {
		return nil, err
	}
	dropletsSize := len(droplets)

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	if !notPage {
		pagenum, perPage := page.RequirePage()
		opt.Page = int(pagenum)
		opt.PerPage = int(perPage)
	}

	var response *godo.Response
	var accsrvs []api.AccServer

	for {
		droplets, resp, err := c.client.Droplets.List(opt)
		response = resp
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, d := range droplets {
			list = append(list, d)
		}

		// if we are at the last page, break out the for loop
		if response.Links == nil || response.Links.IsLastPage() {
			break
		}

		page, err := response.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		//only get our need page? strange condition??
		if page == opt.Page {
			break
		}

	}

	//need to update our list total
	if !notPage {
		page.SetItemTotal(uint64(dropletsSize))
	}

	for _, v := range list {
		accsrv := c.newAccserverSpec(&v)
		accsrvs = append(accsrvs, *accsrv)
	}

	return accsrvs, nil
}

//CreateServer create server
func (c *Client) CreateServer(server interface{}) error {
	srv, ok := server.(*api.AccServer)
	if !ok {
		return fmt.Errorf("invalid obj type")
	}

	// request := &godo.DropletCreateRequest{
	// 	Name:   srv.Spec.Name,
	// 	Region: srv.Spec.Region,
	// 	Size:   srv.Spec.Size,
	// 	Image: godo.DropletCreateImage{
	// 		Slug: srv.Spec.Image,
	// 	},
	//
	// }

	//return c.client.Droplets.Create(*godo.DropletCreateRequest)
	return ansible.DeployVPS(api.OperatorDigitalOcean, srv, c.key)
}

//DeleteServer delete server by id
func (c *Client) DeleteServer(id int64) error {

	_, err := c.client.Droplets.Delete(int(id))
	return err
	//return ansible.DeleteVPS(api.OperatorDigitalOcean, id, c.key)
}

func (c *Client) Exec(param *api.AccExec) error {
	execParam := param

	var sshkey string
	switch execParam.Spec.Command {
	case "deployss":
		return ansible.DeployShadowss(api.OperatorDigitalOcean, execParam.Spec.Deploy.HostList, sshkey, execParam.Spec.Deploy.Attribute)
	case "restartSS":
		return ansible.RestartShadowss(api.OperatorDigitalOcean, execParam.Spec.Deploy.HostList, sshkey)
	case "reboot":
		_, _, err := c.client.DropletActions.Reboot(int(execParam.Spec.ID))
		return err
	default:
		return fmt.Errorf("not support command %s", execParam.Spec.Command)
	}

}

func (c *Client) GetSSHKey() (*api.AccSSHKey, error) {
	accKey := &api.AccSSHKey{}
	opt := &godo.ListOptions{}

	for {
		keys, resp, err := c.client.Keys.List(opt)
		if err != nil {
			return nil, err
		}

		// append the current page's droplets to our list
		for _, v := range keys {
			sshKey := api.SSHKey{
				KeyID:       strconv.FormatInt(int64(v.ID), 10),
				Key:         v.PublicKey,
				Name:        v.Name,
				FingerPrint: v.Fingerprint,
			}

			accKey.Spec.Keys = append(accKey.Spec.Keys, sshKey)
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}

	return accKey, nil
}

func (c *Client) ServerExec(serverid int, cmd string) error {
	switch cmd {
	case common.ServerExecRestart:
		_, _, err := c.client.DropletActions.Reboot(int(serverid))
		return err
	}

	return fmt.Errorf("not support command %v", cmd)
}
