package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
)

type Client struct {
	rpc_client *ethclient.Client
}

func (t *Client) Open(url string) (err error) {
	t.rpc_client, err = ethclient.Dial(url)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) Close() {
	if t.rpc_client == nil {
		return
	}
	t.rpc_client.Close()
	t.rpc_client = nil
}
