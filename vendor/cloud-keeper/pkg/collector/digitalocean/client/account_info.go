package client

import "cloud-keeper/pkg/api"

func (c *Client) GetAccount() (*api.AccountInfoSpec, error) {
	info, _, err := c.client.Account.Get()
	if err != nil {
		return nil, err
	}

	apiInfo := &api.AccountInfoSpec{
		DigitalOcean: api.DGAccountInfo{
			DropletLimit:  info.DropletLimit,
			Email:         info.Email,
			EmailVerified: info.EmailVerified,
			UUID:          info.UUID,
		},
	}

	return apiInfo, nil

}
