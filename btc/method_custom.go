package btc

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func (t *Client) sendAndRecv(req interface{}, res interface{}) (err error) {
	chanRes := t.rpc_client.SendCmd(req)
	bt, err := rpcclient.ReceiveFuture(chanRes)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bt, res)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) ScanTxOutSet(addresses ...btcutil.Address) (result *ScanTxOutSetResult, err error) {
	cmd := NewScanTxOutSetCmd("start", addresses)
	result = &ScanTxOutSetResult{}
	err = t.sendAndRecv(cmd, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (t *Client) SignRawTransactionWithKey(tx *wire.MsgTx, inputs []RawTxInput, privKeys []string) (txSigned *wire.MsgTx, err error) {
	var txid string
	if tx != nil {
		// Serialize the transaction and convert to hex string.
		buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
		if err := tx.Serialize(buf); err != nil {
			return nil, err
		}
		txid = hex.EncodeToString(buf.Bytes())
	}

	cmd := NewSignRawTransactionCmd(txid, &inputs, &privKeys, nil)
	result := &SignRawTransactionResult{}
	err = t.sendAndRecv(cmd, result)
	if err != nil {
		return nil, err
	}

	// Decode the serialized transaction hex to raw bytes.
	serializedTx, err := hex.DecodeString(result.Hex)
	if err != nil {
		return nil, err
	}

	// Deserialize the transaction and return it.
	var msgTx wire.MsgTx
	err = msgTx.Deserialize(bytes.NewReader(serializedTx))
	if err != nil {
		return nil, err
	}
	return &msgTx, nil
}

//---------------------------------------------------------------------------------//
// custom cmd struct

func init() {
	// 커스텀 cmd 등록
	btcjson.MustRegisterCmd("scantxoutset", (*ScanTxOutSetCmd)(nil), btcjson.UFWalletOnly)
	btcjson.MustRegisterCmd("signrawtransactionwithkey", (*SignRawTransactionCmd)(nil), btcjson.UFWalletOnly)
}

type ScanTxOutSetCmd struct {
	Action      *string
	ScanObjects *[]string
}

func NewScanTxOutSetCmd(Action string, addresses []btcutil.Address) *ScanTxOutSetCmd {
	cmd := &ScanTxOutSetCmd{}
	cmd.Action = &Action
	scans := make([]string, 0, len(addresses))
	for _, address := range addresses {
		s_scan := fmt.Sprintf("addr(%s)", address.String())
		scans = append(scans, s_scan)
	}
	cmd.ScanObjects = &scans
	return cmd
}

// ScanTxOutSetResult models a successful response from the scantxoutset request.
type ScanTxOutSetResult struct {
	Success     bool       `json:"success"`
	TxOuts      int64      `json:"txouts"`
	Height      int32      `json:"height"`
	BestBlock   string     `json:"bestblock"`
	Unspents    []UnSpents `json:"unspents"`
	TotalAmount float64    `json:"total_amount"`
}

type UnSpents struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	Desc         string  `json:"desc,omitempty"`
	Amount       float64 `json:"amount"`
	Height       int32   `json:"height"`
}

// RawTxInput models the data needed for raw transaction input that is used in
// the SignRawTransactionCmd struct.
type RawTxInput struct {
	Txid         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript"`
	Amount       float64 `json:"amount"`
}

// SignRawTransactionCmd defines the signrawtransaction JSON-RPC command.
type SignRawTransactionCmd struct {
	RawTx    string        `json:"hexstring"`
	PrivKeys *[]string     `json:"privkeys"`
	Inputs   *[]RawTxInput `json:"prevtxs"`
	Flags    *string       `jsonrpcdefault:"\"ALL\""`
}

// NewSignRawTransactionCmd returns a new instance which can be used to issue a
// signrawtransaction JSON-RPC command.
//
// The parameters which are pointers indicate they are optional.  Passing nil
// for optional parameters will use the default value.
func NewSignRawTransactionCmd(hexEncodedTx string, inputs *[]RawTxInput, privKeys *[]string, flags *string) *SignRawTransactionCmd {
	return &SignRawTransactionCmd{
		RawTx:    hexEncodedTx,
		Inputs:   inputs,
		PrivKeys: privKeys,
		Flags:    flags,
	}
}

// SignRawTransactionError models the data that contains script verification
// errors from the signrawtransaction request.
type SignRawTransactionError struct {
	TxID      string `json:"txid"`
	Vout      uint32 `json:"vout"`
	ScriptSig string `json:"scriptSig"`
	Sequence  uint32 `json:"sequence"`
	Error     string `json:"error"`
}

// SignRawTransactionResult models the data from the signrawtransaction
// command.
type SignRawTransactionResult struct {
	Hex      string                    `json:"hex"`
	Complete bool                      `json:"complete"`
	Errors   []SignRawTransactionError `json:"errors,omitempty"`
}
