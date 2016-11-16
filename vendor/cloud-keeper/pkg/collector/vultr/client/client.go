package client

import (
	vultr "cloud-keeper/pkg/collector/vultr/client/lib"
	"gofreezer/pkg/api/unversioned"
)

var AccountInfoType = unversioned.TypeMeta{
	Kind:       "AccountInfo",
	APIVersion: "v1",
}

var AccServerType = unversioned.TypeMeta{
	Kind:       "AccServer",
	APIVersion: "v1",
}

var AccServerSSHKeyType = unversioned.TypeMeta{
	Kind:       "AccServerSSHKey",
	APIVersion: "v1",
}

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
