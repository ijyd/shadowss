package web

import (
	"github.com/orvice/shadowsocks-go/mu/user"
)

var (
	client = new(Client)
)

type Client struct {
	baseUrl string
	key     string
	nodeId  int
}

func (c *Client) setBaseUrl(baseUrl string) {
	c.baseUrl = baseUrl
}

func (c *Client) setKey(key string) {
	c.key = key
}

func (c *Client) setNodeId(id int) {
	c.nodeId = id
}

func (c *Client) GetUsers() ([]user.User, error) {
	var datas []*User

	return datas, nil
}
