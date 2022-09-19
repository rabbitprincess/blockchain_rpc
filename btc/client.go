package btc

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/rpcclient"
)

type Client struct {
	params     *chaincfg.Params
	rpc_client *rpcclient.Client
}

func (t *Client) Open(params *chaincfg.Params, host, id, pw string) (err error) {
	t.params = params
	t.rpc_client, err = rpcclient.New(&rpcclient.ConnConfig{
		Host:         host,
		User:         id,
		Pass:         pw,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) Close() {
	if t.rpc_client == nil {
		return
	}
	t.rpc_client.Shutdown()
	t.rpc_client = nil
}
