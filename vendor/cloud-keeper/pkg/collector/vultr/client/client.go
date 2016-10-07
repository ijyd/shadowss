package client

import vultr "cloud-keeper/pkg/collector/vultr/client/lib"

type Client struct {
	vultrClient *vultr.Client
	key         string
}

func NewClient(apiKey string) *Client {
	return &Client{
		vultrClient: vultr.NewClient(apiKey, nil),
		key:         apiKey,
	}
}
