package btc

import (
	"fmt"
	"math"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// transfer without wallet support ( use privkey only )
type RawTx struct {
	client      *Client
	balanceAddr btcutil.Address
	fee         btcutil.Amount

	fromPrivKeys []string
	fromAddrs    []btcutil.Address
	toAmounts    map[btcutil.Address]btcutil.Amount
}

func (t *RawTx) Init(client *Client, balanceAddr string, fee float64) (err error) {
	t.client = client
	t.fromPrivKeys = make([]string, 0, 10)
	t.fromAddrs = make([]btcutil.Address, 0, 10)
	t.toAmounts = make(map[btcutil.Address]btcutil.Amount)

	t.balanceAddr, err = btcutil.DecodeAddress(balanceAddr, t.client.params)
	if err != nil {
		return err
	}
	t.fee, err = btcutil.NewAmount(fee)
	if err != nil {
		return err
	}
	return nil
}

func (t *RawTx) AddFrom(privKey, address string) (err error) {
	btcAddr, err := btcutil.DecodeAddress(address, t.client.params)
	if err != nil {
		return err
	}

	t.fromPrivKeys = append(t.fromPrivKeys, privKey)
	t.fromAddrs = append(t.fromAddrs, btcAddr)
	return nil
}

func (t *RawTx) AddTo(address string, amount btcutil.Amount) (err error) {
	addr, err := btcutil.DecodeAddress(address, t.client.params)
	if err != nil {
		return err
	}
	t.toAmounts[addr] = amount
	return nil
}

func (t *RawTx) SendTx() (txid string, err error) {
	utxo, err := t.utxoGet()
	if err != nil {
		return "", err
	}

	msgTx, leftAmount, err := t.make(utxo)
	if err != nil {
		return "", err
	}
	msgTxFunded, err := t.fund(msgTx, utxo, leftAmount)
	if err != nil {
		return "", err
	}
	msgTxSigned, err := t.sign(msgTxFunded, utxo)
	if err != nil {
		return "", err
	}

	return t.send(msgTxSigned)
}

//--------------------------------------------------------------------------------//
// method

type utxo struct {
	Txid         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	FromAddr     string  `json:"from_address"`
	FromAmount   float64 `json:"from_amount"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript"`
}

func (t *RawTx) utxoGet() (utxos []*utxo, err error) {
	// 1. separate address in wallet / out wallet
	var addressesInWallet, addressesOutWallet []btcutil.Address
	for _, address := range t.fromAddrs {
		addrInfo, err := t.client.GetAddressInfo(address.String())
		if err != nil {
			return nil, err
		}
		isInWallet := addrInfo.IsMine
		if isInWallet == true {
			addressesInWallet = append(addressesInWallet, address)
		} else {
			addressesOutWallet = append(addressesOutWallet, address)
		}
	}

	// 2. get utxo in wallet
	var utxosInWallet []btcjson.ListUnspentResult
	if len(addressesInWallet) > 0 {
		utxosInWallet, err = t.client.GetListUnspent(0, 0, addressesInWallet...)
		if err != nil {
			return nil, err
		}
	}

	// 3. get utxo out wallet
	var utxosOutWallet []UnSpents
	if len(addressesOutWallet) > 0 {
		scanTxOutSet, err := t.client.ScanTxOutSet(addressesOutWallet...)
		if err != nil {
			return nil, err
		}
		utxosOutWallet = scanTxOutSet.Unspents
	}

	// 4. make utxo
	utxos = make([]*utxo, 0, len(utxosInWallet))
	for _, utxoInWallet := range utxosInWallet {
		utxo := &utxo{
			Txid:         utxoInWallet.TxID,
			Vout:         utxoInWallet.Vout,
			FromAddr:     utxoInWallet.Address,
			FromAmount:   utxoInWallet.Amount,
			ScriptPubKey: utxoInWallet.ScriptPubKey,
			RedeemScript: utxoInWallet.RedeemScript,
		}
		utxos = append(utxos, utxo)
	}

	for _, utxoOutWallet := range utxosOutWallet {
		walletAddr, err := decodeDescToAddr(utxoOutWallet.Desc)
		if err != nil {
			return nil, err
		}
		utxo := &utxo{
			Txid:         utxoOutWallet.TxID,
			Vout:         utxoOutWallet.Vout,
			FromAddr:     walletAddr,
			FromAmount:   utxoOutWallet.Amount,
			ScriptPubKey: utxoOutWallet.ScriptPubKey,
		}
		utxos = append(utxos, utxo)
	}

	return utxos, nil
}

func decodeDescToAddr(desc string) (addr string, err error) {
	posFront := strings.Index(desc, "addr(")
	posEnd := strings.Index(desc, ")")

	if posFront == -1 || posEnd == -1 || posFront+5 >= posEnd {
		return "", fmt.Errorf("invalid desc format | %s", desc)
	}
	return desc[posFront+5 : posEnd], nil
}

func (t *RawTx) make(utxos []*utxo) (msgTx *wire.MsgTx, leftAmount btcutil.Amount, err error) {
	// conv utxo to input
	txsInput := make([]btcjson.TransactionInput, 0, len(utxos))
	for _, utxo := range utxos {
		txInput := btcjson.TransactionInput{}
		txInput.Txid = utxo.Txid
		txInput.Vout = utxo.Vout
		txsInput = append(txsInput, txInput)
	}

	msgTx, err = t.client.rpc.CreateRawTransaction(txsInput, t.toAmounts, nil)
	if err != nil {
		return nil, 0, err
	}

	// set amount left ( from amounts total - to amounts total )
	leftAmount, err = btcutil.NewAmount(0)
	if err != nil {
		return nil, 0, err
	}
	for _, utxo := range utxos {
		amount, err := btcutil.NewAmount(utxo.FromAmount)
		if err != nil {
			return nil, 0, err
		}
		leftAmount += amount
	}
	for _, amount := range t.toAmounts {
		leftAmount -= amount
	}

	if leftAmount < 0 {
		return nil, 0, fmt.Errorf("left amount is under zero")
	}

	return msgTx, leftAmount, nil
}

func (t *RawTx) fund(msgTx *wire.MsgTx, utxos []*utxo, leftAmount btcutil.Amount) (msgTxFunded *wire.MsgTx, err error) {
	// set amount left without fee = sum(vin) - sum(vout) - fee
	var leftAmountWithoutFee int64
	{
		// sign tx ( to calculate fee )
		msgTxSigned, err := t.sign(msgTx, utxos)
		if err != nil {
			return nil, err
		}

		// get transaction size
		_, vsize := getRawTxSize(msgTxSigned)

		// 단위 수수료 변경 ( btc per kb -> satoshi per byte )
		feePerByteSatoshi := int64(t.fee.ToUnit(btcutil.AmountSatoshi + btcutil.AmountKiloBTC))
		if feePerByteSatoshi < 0 {
			return nil, fmt.Errorf("invalid fee per byte ( satoshi ) | %v", feePerByteSatoshi)
		}

		// p2pkh size = 34
		feeRawTx := int64(vsize+34) * feePerByteSatoshi

		// left amount 를 satoshi 단위로 변경
		leftAmountSatoshi := int64(leftAmount.ToUnit(btcutil.AmountSatoshi))
		if leftAmountSatoshi < 0 {
			return nil, fmt.Errorf("invalid left amount ( satoshi ) | %v", leftAmount)
		}

		// amount left 에서 전체 수수료로 뺀 값이 balance address 에 들어갈 금액
		leftAmountWithoutFee = leftAmountSatoshi - feeRawTx
		if leftAmountWithoutFee < 0 {
			return nil, fmt.Errorf("not enough amount left without fee | left : %v | fee - %v", leftAmount, feeRawTx)
		}
	}

	// add balance address to vout
	toAddr, err := txscript.PayToAddrScript(t.balanceAddr)
	if err != nil {
		return nil, err
	}
	msgTx.AddTxOut(wire.NewTxOut(leftAmountWithoutFee, toAddr))

	return msgTx, nil
}

func getRawTxSize(msgTx *wire.MsgTx) (size, vsize int) {
	size = msgTx.SerializeSize()
	sizeNoWitness := msgTx.SerializeSizeStripped()

	weight := (sizeNoWitness)*3 + size
	// round up
	vsize = int(math.Round(float64(weight) / 4))
	return size, vsize
}

func (t *RawTx) sign(msgTxFunded *wire.MsgTx, utxos []*utxo) (msgTxSigned *wire.MsgTx, err error) {
	rawTxInput := make([]RawTxInput, 0, len(utxos))
	for _, utxo := range utxos {
		rawTxInput = append(rawTxInput, RawTxInput{
			Txid:         utxo.Txid,
			Vout:         utxo.Vout,
			ScriptPubKey: utxo.ScriptPubKey,
			RedeemScript: utxo.RedeemScript,
			Amount:       utxo.FromAmount,
		})
	}

	msgTxSigned, err = t.client.SignRawTransactionWithKey(msgTxFunded, rawTxInput, t.fromPrivKeys)
	if err != nil {
		return nil, err
	}
	return msgTxSigned, nil
}

func (t *RawTx) send(msgTxSigned *wire.MsgTx) (txid string, err error) {
	hash, err := t.client.rpc.SendRawTransaction(msgTxSigned, false)
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}
