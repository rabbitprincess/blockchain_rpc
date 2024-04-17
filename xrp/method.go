package xrp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rabbitprincess/blockchain_rpc/xrp/types"
)

func (t *Client) sendCmd(req types.CmdReq, res types.CmdRes) (err error) {
	// marshal
	btReq, err := req.Marshal()
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(btReq)

	// httpReq
	httpReq, err := http.NewRequest("POST", t.url, buf)
	if err != nil {
		return err
	}
	httpRes, err := t.rpc.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	btRes, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	// unmarshal
	err = res.Unmarshal(btRes)
	if err != nil {
		return err
	}
	// 후처리 - error
	resErr := res.Error()
	if resErr.ErrCode != types.Ok {
		return resErr
	}

	return nil
}

//-------------------------------------------------------------------------------------------//
// server

func (t *Client) GetServerInfo() (res *types.Res_serverInfo_result, err error) {
	cmdReq := types.Req_serverInfo{Method: types.Cmd_ServerInfo}
	cmdRes := types.Res_serverInfo{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result
	return res, nil
}

func (t *Client) GetBaseFee() (baseFee string, err error) {
	req := types.Req_fee{Method: types.Cmd_Fee}
	res := types.Res_fee{}
	err = t.sendCmd(&req, &res)
	if err != nil {
		return "", err
	}
	// 에러 처리
	rpcError := res.Error()
	if rpcError.ErrCode != types.Ok {
		return "", rpcError
	}

	baseFee = res.Result.Drops.BaseFee
	return baseFee, nil
}

//-------------------------------------------------------------------------------------------//
// account

func (t *Client) GetAccountInfo(account string) (res *types.Res_accountInfo_result, err error) {
	cmdReq := types.Req_accountInfo{
		Method: types.Cmd_AccountInfo,
		Params: []types.Req_accountInfo_params{{Account: account, Strict: true, Queue: true}},
	}
	cmdRes := types.Res_accountInfo{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}

	res = &cmdRes.Result
	return res, nil
}

func (t *Client) GetAccountTx(account string, ledgerMin, ledgerMax int) (res *types.Res_accountTx_result, err error) {
	cmdReq := types.Req_accountTx{
		Method: types.Cmd_AccountTx,
		Params: []types.Req_accountTx_params{
			{
				Account:        account,
				LedgerHash:     "",
				LedgerIndexMin: ledgerMin,
				LedgerIndexMax: ledgerMax,
			},
		},
	}
	cmdRes := types.Res_accountTx{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result
	return res, nil
}

func (t *Client) WalletPropose(privKey string) (res *types.Res_walletPropose_Result, err error) {
	cmdReq := types.Req_walletPropose{Method: types.Cmd_WalletPropose}
	if privKey != "" {
		cmdReq.Params = []types.Req_walletPropose_params{{Seed: privKey, KeyType: types.Secp256k1}}
	}
	cmdRes := types.Res_walletPropose{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result
	return res, nil
}

//-------------------------------------------------------------------------------------------//
// tx

func (t *Client) MakeTransaction(txHash string) (res *types.TransactionRes, err error) {
	cmdReq := types.Req_tx{Method: types.Cmd_Tx}
	cmdRes := types.Res_tx{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result
	return res, nil
}

func (t *Client) SignTransaction(tx *types.TransactionRes, privKey string) (txid string, err error) {
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}

	cmdReq := types.Req_sign{
		Method: types.Cmd_Sign,
		Params: []types.Req_sign_params{
			{
				Secret: privKey,
				TxJson: txJSON,
			},
		},
	}
	cmdRes := types.Res_sign{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return "", err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return "", rpcError
	}
	txid = string(cmdRes.Result.TxBlob)
	return txid, nil
}

func (t *Client) SendTransaction(txid string) (res *types.Res_submit_result, err error) {
	cmdReq := types.Req_submit{Method: types.Cmd_Submit}
	cmdRes := types.Res_submit{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result
	return res, nil
}

//-------------------------------------------------------------------------------------------//
// ledger

func (t *Client) GetLedgerByNumber(ledgerNumber int, includeTx bool) (res *types.LedgerRes, err error) {
	cmdReq := types.Req_ledger{
		Method: types.Cmd_Ledger,
		Params: []types.Req_ledger_params{
			{
				LedgerIndex:  strconv.FormatInt(int64(ledgerNumber), 10),
				Accounts:     false,
				Full:         false,
				Transactions: includeTx,
				Expand:       includeTx,
				OwnerFunds:   false,
			},
		},
	}
	cmdRes := types.Res_ledger{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result.Ledger
	return res, nil
}

func (t *Client) GetLedgerByHash(ledgerHash string, includeTx bool) (res *types.LedgerRes, err error) {
	cmdReq := types.Req_ledger{
		Method: types.Cmd_Ledger,
		Params: []types.Req_ledger_params{
			{
				LedgerHash:   ledgerHash,
				Accounts:     false,
				Full:         false,
				Transactions: includeTx,
				Expand:       includeTx,
				OwnerFunds:   false,
			},
		},
	}
	cmdRes := types.Res_ledger{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result.Ledger
	return res, nil
}

func (t *Client) GetLedgerLast(includeTx bool) (res *types.LedgerRes, err error) {
	cmdReq := types.Req_ledger{
		Method: types.Cmd_Ledger,
		Params: []types.Req_ledger_params{
			{
				LedgerIndex:  "validated",
				Accounts:     false,
				Full:         false,
				Transactions: includeTx,
				Expand:       includeTx,
				OwnerFunds:   false,
			},
		},
	}
	cmdRes := types.Res_ledger{}
	err = t.sendCmd(&cmdReq, &cmdRes)
	if err != nil {
		return nil, err
	}

	// 에러 처리
	rpcError := cmdRes.Error()
	if rpcError.ErrCode != types.Ok {
		return nil, rpcError
	}
	res = &cmdRes.Result.Ledger
	return res, nil
}
