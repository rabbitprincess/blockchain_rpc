package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewClient(url string) (*Client, error) {
	client := &Client{}
	err := client.Open(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

type Client struct {
	rpc *ethclient.Client
}

func (t *Client) Open(url string) (err error) {
	t.rpc, err = ethclient.Dial(url)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) Close() {
	if t.rpc == nil {
		return
	}
	t.rpc.Close()
	t.rpc = nil
}
