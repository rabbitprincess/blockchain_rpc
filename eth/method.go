package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
)

func (t *Client) GetServerInfo() (status *ethereum.SyncProgress, err error) {
	return t.rpc_client.SyncProgress(context.Background())
}

//-------------------------------------------------------------------------------------------//
// address

func (t *Client) GetAddressBalance() (status *ethereum.SyncProgress, err error) {

	return
}

//-------------------------------------------------------------------------------------------//
// fee

func (t *Client) SuggestGasInfo() (gasPriceWei, gasTipCapWei *big.Int, err error) {
	context := context.Background()

	gasPriceWei, err = t.rpc_client.SuggestGasPrice(context)
	if err != nil {
		return nil, nil, err
	}
	gasTipCapWei, err = t.rpc_client.SuggestGasTipCap(context)
	if err != nil {
		return nil, nil, err
	}

	return gasPriceWei, gasTipCapWei, nil

}
