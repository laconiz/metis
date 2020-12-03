package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func New(addr string) (*Client, error) {

	client, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		return nil, fmt.Errorf("new consul client error: %w", err)
	}

	if _, err = client.Catalog().Datacenters(); err != nil {
		return nil, fmt.Errorf("check consul connection error: %w", err)
	}

	return &Client{addr: addr, Client: client}, nil
}

type Client struct {
	addr string
	*api.Client
}

func (client *Client) Address() string {
	return client.addr
}

func (client *Client) KV() *KV {
	return &KV{kv: client.Client.KV()}
}

func (client *Client) Watcher() *Watcher {
	return &Watcher{client: client}
}

func (client *Client) Register() *Register {
	return &Register{client: client}
}
