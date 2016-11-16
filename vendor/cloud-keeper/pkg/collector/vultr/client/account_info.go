package client

import "cloud-keeper/pkg/api"

func (c *Client) GetAccount() (*api.AccountInfoSpec, error) {
	info, err := c.vultrClient.GetAccountInfo()
	if err != nil {
		return nil, err
	}

	apiInfo := &api.AccountInfoSpec{
		Vultr: api.VultrAccountInfo{
			Balance:           info.Balance,
			PendingCharges:    info.PendingCharges,
			LastPaymentDate:   info.LastPaymentDate,
			LastPaymentAmount: info.LastPaymentAmount,
		},
	}

	return apiInfo, nil
}
