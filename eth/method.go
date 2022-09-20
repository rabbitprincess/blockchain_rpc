package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func (t *Client) GetServerInfo() (status *ethereum.SyncProgress, err error) {
	return t.rpc_client.SyncProgress(context.Background())
}

//-------------------------------------------------------------------------------------------//
// address

func (t *Client) GetAddressBalance(address string) (balance string, err error) {
	wei, err := t.rpc_client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	return Conv_WeiToEth(wei.String())
}

func (t *Client) GetAddressNonce(address string) (nonce uint64, err error) {
	return t.rpc_client.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func (t *Client) GetAddressCode(address string, blockNumber uint64) (byteCode []byte, err error) {
	return t.rpc_client.CodeAt(context.Background(), common.HexToAddress(address), big.NewInt(int64(blockNumber)))
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
