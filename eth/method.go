package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//-------------------------------------------------------------------------------------------//
// server

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
// tx

func (t *Client) GetTxInfo(txid string) (txInfo *types.Transaction, isPending bool, err error) {
	ethTxHash := common.HexToHash(txid)
	return t.rpc_client.TransactionByHash(context.Background(), ethTxHash)
}

func (t *Client) GetTxReceipt(txid string) (txReceipt *types.Receipt, err error) {
	ethTxHash := common.HexToHash(txid)
	return t.rpc_client.TransactionReceipt(context.Background(), ethTxHash)
}

func (t *Client) SendTx(tx *types.Transaction) (err error) {
	return t.rpc_client.SendTransaction(context.Background(), tx)
}

//-------------------------------------------------------------------------------------------//
// block

func (t *Client) GetBlockMostRecentNumber() (blockNumber uint64, err error) {
	return t.rpc_client.BlockNumber(context.Background())
}

func (t *Client) GetBlockInfo(blockNumber int64) (blockInfo *types.Block, err error) {
	return t.rpc_client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
}

func (t *Client) GetBlockInfoByHash(blockHash string) (blockInfo *types.Block, err error) {
	return t.rpc_client.BlockByHash(context.Background(), common.HexToHash(blockHash))
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
