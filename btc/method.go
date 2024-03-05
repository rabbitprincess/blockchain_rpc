package btc

import (
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

//---------------------------------------------------------------------------//
// wallet

func (t *Client) UnlockPassphrase(passphrase string, timeoutSec int64) (err error) {
	if timeoutSec <= 0 {
		timeoutSec = 60 // default sec
	}
	return t.rpc.WalletPassphrase(passphrase, timeoutSec)
}

func (t *Client) ImportPrivKey(privKey string) (err error) {
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return err
	}
	return t.rpc.ImportPrivKey(wif)
}

func (t *Client) GetBalanceTotal() (balance btcutil.Amount, err error) {
	return t.GetBalance("*")
}

func (t *Client) GetBalance(address string) (balance btcutil.Amount, err error) {
	return t.rpc.GetBalance(address)
}

func (t *Client) GetListUnspent(minBlock, maxBlock int, addresses ...btcutil.Address) (unspents []btcjson.ListUnspentResult, err error) {
	if minBlock <= 0 {
		minBlock = 1
	}
	if maxBlock <= 0 {
		maxBlock = 99999999
	}

	if len(addresses) == 0 { // get all unspent
		return t.rpc.ListUnspentMinMax(minBlock, maxBlock)
	} else {
		return t.rpc.ListUnspentMinMaxAddresses(minBlock, maxBlock, addresses)
	}
}

//---------------------------------------------------------------------------//
// address

func (t *Client) GetNewAddress() (privkey, address string, err error) {
	btcAddr, err := t.rpc.GetNewAddress("")
	if err != nil {
		return "", "", err
	}
	btcPrivKey, err := t.rpc.DumpPrivKey(btcAddr)
	if err != nil {
		return "", "", err
	}

	return btcPrivKey.String(), btcAddr.EncodeAddress(), nil
}

func (t *Client) GetAddressInfo(address string) (addrInfo *btcjson.GetAddressInfoResult, err error) {
	return t.rpc.GetAddressInfo(address)
}

func (t *Client) ValidateAddress(address string) (validateAddr *btcjson.ValidateAddressWalletResult, err error) {
	btcAddr, err := btcutil.DecodeAddress(address, t.params)
	if err != nil {
		return nil, err
	}
	return t.rpc.ValidateAddress(btcAddr)
}

func (t *Client) DumpPrivKey(address string) (privKey string, err error) {
	btcAddr, err := btcutil.DecodeAddress(address, t.params)
	if err != nil {
		return "", err
	}
	wif, err := t.rpc.DumpPrivKey(btcAddr)
	if err != nil {
		return "", err
	}
	return wif.String(), nil
}

//---------------------------------------------------------------------------//
// tx

func (t *Client) GetTxInfo(txid string) (txInfo *btcjson.GetTransactionResult, err error) {
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	txInfo, err = t.rpc.GetTransaction(hash)
	if err != nil {
		return nil, err
	}
	return txInfo, nil
}

func (t *Client) GetRawTxInfo(txid string) (rawTxInfo *btcjson.TxRawResult, err error) {
	hash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return nil, err
	}
	rawTxInfo, err = t.rpc.GetRawTransactionVerbose(hash)
	if err != nil {
		return nil, err
	}
	return rawTxInfo, nil
}

//---------------------------------------------------------------------------//
// block

func (t *Client) GetBestBlock() (blockHash string, err error) {
	btcBlockHash, err := t.rpc.GetBestBlockHash()
	if err != nil {
		return "", err
	}
	return btcBlockHash.String(), nil
}

func (t *Client) GetBlockHash(blockNumber int64) (blockHash string, err error) {
	btcBlockHash, err := t.rpc.GetBlockHash(blockNumber)
	if err != nil {
		return "", err
	}
	return btcBlockHash.String(), nil
}

func (t *Client) GetBlockInfo(blockHash string) (blockInfo *btcjson.GetBlockVerboseResult, err error) {
	btcBlockHash, err := chainhash.NewHashFromStr(blockHash)
	if err != nil {
		return nil, err
	}
	return t.rpc.GetBlockVerbose(btcBlockHash)
}

func (t *Client) GetBlockInfoWithTx(blockHash string) (blockInfo *btcjson.GetBlockVerboseTxResult, err error) {
	btcBlockHash, err := chainhash.NewHashFromStr(blockHash)
	if err != nil {
		return nil, err
	}
	return t.rpc.GetBlockVerboseTx(btcBlockHash)
}

//---------------------------------------------------------------------------//
// fee

func (t *Client) GetSmartFee(confTargetBlock int64, feeEstimateMode *btcjson.EstimateSmartFeeMode) (smartFee btcutil.Amount, err error) {
	if confTargetBlock <= 0 {
		confTargetBlock = 10
	}
	if feeEstimateMode == nil {
		feeEstimateMode = &btcjson.EstimateModeUnset
	}
	result, err := t.rpc.EstimateSmartFee(confTargetBlock, feeEstimateMode)
	if err != nil {
		return 0, err
	} else if len(result.Errors) != 0 {
		err = fmt.Errorf("%v", result.Errors)
		return 0, err
	}
	smartFee, err = btcutil.NewAmount(*result.FeeRate)
	if err != nil {
		return 0, err
	}
	return smartFee, nil
}

func (t *Client) SetFee(fee btcutil.Amount) (err error) {
	return t.rpc.SetTxFee(fee)
}

//---------------------------------------------------------------------------//
// transfer

func (t *Client) SendCoin(addrTo string, amount btcutil.Amount) (txid string, err error) {
	to, err := btcutil.DecodeAddress(addrTo, t.params)
	if err != nil {
		return "", err
	}

	chainHash, err := t.rpc.SendToAddress(to, amount)
	if err != nil {
		return "", err
	}
	return chainHash.String(), nil
}

func (t *Client) SendCoinMany(mapAmounts map[string]btcutil.Amount) (txid string, err error) {
	mapAddrAmounts := make(map[btcutil.Address]btcutil.Amount)
	for address, amount := range mapAmounts {
		btcAddr, err := btcutil.DecodeAddress(address, t.params)
		if err != nil {
			return "", err
		}
		mapAddrAmounts[btcAddr] = amount
	}
	chainHash, err := t.rpc.SendMany("", mapAddrAmounts)
	if err != nil {
		return "", err
	}
	return chainHash.String(), nil
}
