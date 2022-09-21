package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type RawTx struct {
	client *Client
	param  *params.ChainConfig

	gasFeeCap *big.Int
	gasTipCap *big.Int
	gasLimit  uint64

	fromPrivKey string
	fromAddr    string
	nonce       uint64

	tokenAddr string
	decimal   uint8
	toAddr    string
	toAmount  string
}

func (t *RawTx) Init(client *Client, gasPrice, gasTip *big.Int) {

}

func (t *RawTx) SendTx() (txid string, err error) {
	tx, err := t.make()
	if err != nil {
		return "", err
	}
	txSigned, err := t.sign(tx, t.fromPrivKey)
	if err != nil {
		return "", err
	}
	txid, err = t.send(txSigned)
	if err != nil {
		return "", err
	}
	return txid, nil
}

//--------------------------------------------------------------------------------//
// method

func (t *RawTx) make() (tx *types.Transaction, err error) {
	var byteCodeErc20Transfer []byte
	var bigAmountWei *big.Int
	var toAddress common.Address
	if t.tokenAddr == "" {
		// eth transfer
		toAmountWei, err := Conv_EthToWei(t.toAmount)
		if err != nil {
			return nil, err
		}

		bigAmountWei, _ = big.NewInt(0).SetString(toAmountWei, 10)
		toAddress = common.HexToAddress(t.toAddr)
	} else {
		// erc20 transfer
		toAmountWei, err := Conv_UnitToWei(t.toAmount, t.decimal)
		if err != nil {
			return nil, err
		}

		bigAmountWei, _ = big.NewInt(0).SetString(toAmountWei, 10)
		byteCodeErc20Transfer = t.client.MakeErc20TransferBytecode(t.toAddr, bigAmountWei)
		bigAmountWei.SetInt64(0)                     // to amount = 0
		toAddress = common.HexToAddress(t.tokenAddr) // to address = token address
	}

	// make dynamic fee transaction
	tx = types.NewTx(types.DynamicFeeTx{
		Nonce:     t.nonce,
		Value:     bigAmountWei,
		To:        &toAddress,
		Gas:       t.gasLimit,
		GasTipCap: t.gasTipCap,
		GasFeeCap: t.gasFeeCap,
		Data:      byteCodeErc20Transfer,
	})

	return tx, nil
}

func (t *RawTx) sign(tx *types.Transaction, privKey string) (txSigned *types.Transaction, err error) {
	ecdsaPrivKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}

	return types.SignTx(tx, types.LatestSigner(t.param), ecdsaPrivKey)
}

func (t *RawTx) send(txSigned *types.Transaction) (txid string, err error) {
	err = t.client.SendTx(txSigned)
	if err != nil {
		return "", err
	}
	txid = txSigned.Hash().Hex()
	return txid, nil
}
