package client

import (
	"encoding/json"

	"cloud-keeper/pkg/api"
)

func (c *Client) GetAccount() ([]byte, error) {
	info, err := c.vultrClient.GetAccountInfo()
	if err != nil {
		return nil, err
	}

	information := make(map[string]interface{}, 1)
	information[api.OperatorVultr] = info

	apiInfo := api.AccountInfo{
		TypeMeta:    api.AccountInfoType,
		Information: information,
	}

	return json.Marshal(&apiInfo)
}
