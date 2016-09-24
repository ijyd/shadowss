package client

import (
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

//Client wrap a digitalocean client
type Client struct {
	client *godo.Client
	key    string
}

type tokenSource struct {
	accessToken string
}

func (c *Client) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: c.key,
	}
	return token, nil
}

//NewClient create a DigitalOcean client
func NewClient(apiKey string) *Client {
	c := &Client{
		key: apiKey,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, c)
	c.client = godo.NewClient(oauthClient)
	return c
}
