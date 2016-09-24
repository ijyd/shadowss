package client

import (
	"cloud-keeper/pkg/api"
	"encoding/json"
)

func (c *Client) GetAccount() ([]byte, error) {
	info, _, err := c.client.Account.Get()
	if err != nil {
		return nil, err
	}

	information := make(map[string]interface{}, 1)
	information[api.OperatorDigitalOcean] = info

	apiInfo := api.AccountInfo{
		TypeMeta:    api.AccountInfoType,
		Information: information,
	}

	return json.Marshal(&apiInfo)
}
