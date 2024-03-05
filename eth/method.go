package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//-------------------------------------------------------------------------------------------//
// server

func (t *Client) GetServerInfo() (status *ethereum.SyncProgress, err error) {
	return t.rpc.SyncProgress(context.Background())
}

func (t *Client) SuggestGasInfo() (gasPriceWei, gasTipCapWei *big.Int, err error) {
	context := context.Background()

	gasPriceWei, err = t.rpc.SuggestGasPrice(context)
	if err != nil {
		return nil, nil, err
	}
	gasTipCapWei, err = t.rpc.SuggestGasTipCap(context)
	if err != nil {
		return nil, nil, err
	}

	return gasPriceWei, gasTipCapWei, nil

}

//-------------------------------------------------------------------------------------------//
// address

func (t *Client) GetNewAddress() (privKey, address string, err error) {
	// get private key
	priv, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	// get public key
	pub, ok := priv.Public().(*ecdsa.PublicKey)
	if ok == false {
		return "", "", fmt.Errorf("public key ECDSA casting failed | key info - %v", pub)
	}

	// get wallet address
	addr := crypto.PubkeyToAddress(*pub)

	// encode
	privKey = hex.EncodeToString(crypto.FromECDSA(priv))
	address = addr.Hex()
	return privKey, address, nil
}

func (t *Client) GetAddressBalance(address string) (balance string, err error) {
	wei, err := t.rpc.BalanceAt(context.Background(), common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	return Conv_WeiToEth(wei.String())
}

func (t *Client) GetAddressNonce(address string) (nonce uint64, err error) {
	return t.rpc.PendingNonceAt(context.Background(), common.HexToAddress(address))
}

func (t *Client) GetAddressCode(address string, blockNumber uint64) (byteCode []byte, err error) {
	var blockNumberBig *big.Int
	if blockNumber > 0 {
		blockNumberBig = big.NewInt(int64(blockNumber))
	}
	return t.rpc.CodeAt(context.Background(), common.HexToAddress(address), blockNumberBig)
}

func (t *Client) ValidAddress(address string) (isContract bool, err error) {
	valid := AddressValid(address)
	if valid == false {
		return valid, fmt.Errorf("invalid address | %s", address)
	}
	code, err := t.GetAddressCode(address, 0)
	if err != nil {
		return false, err
	}
	isContract = len(code) > 0
	return isContract, nil
}

//-------------------------------------------------------------------------------------------//
// tx

func (t *Client) GetTxInfo(txid string) (txInfo *types.Transaction, isPending bool, err error) {
	ethTxHash := common.HexToHash(txid)
	return t.rpc.TransactionByHash(context.Background(), ethTxHash)
}

func (t *Client) GetTxReceipt(txid string) (txReceipt *types.Receipt, err error) {
	ethTxHash := common.HexToHash(txid)
	return t.rpc.TransactionReceipt(context.Background(), ethTxHash)
}

func (t *Client) SendTx(tx *types.Transaction) (err error) {
	return t.rpc.SendTransaction(context.Background(), tx)
}

//-------------------------------------------------------------------------------------------//
// block

func (t *Client) GetBlockMostRecent() (blockNumber uint64, err error) {
	return t.rpc.BlockNumber(context.Background())
}

func (t *Client) GetBlockInfo(blockNumber uint64) (blockInfo *types.Block, err error) {
	return t.rpc.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
}

func (t *Client) GetBlockInfoByHash(blockHash string) (blockInfo *types.Block, err error) {
	return t.rpc.BlockByHash(context.Background(), common.HexToHash(blockHash))
}
