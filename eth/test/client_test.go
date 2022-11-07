package eth_test

import (
	"fmt"
	"testing"

	"github.com/gokch/blockchain_rpc/eth"

	"github.com/ethereum/go-ethereum/params"
)

func Test_GetNewAddress(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	s_private_key, s_address, err := client.GetNewAddress()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("private key : %v\n", s_private_key)
	fmt.Printf("address : %v\n", s_address)
}

func Test_Withdraw(t *testing.T) {
	// 출금 입력값
	var (
		privKey  = DEF_private_key
		from     = DEF_address
		to       = "0xD76C201f700E5bAE854BD0722a8B29F87F9a9cCB"
		contract = DEF_token__KCH
		amount   = "0.01"
	)

	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// 출금 과정
	{
		gasPrice, gasTip, err := client.SuggestGasInfo()
		if err != nil {
			t.Fatal(err)
		}
		nonce, err := client.GetAddressNonce(from)
		if err != nil {
			t.Fatal(err)
		}
		tokenInfo, err := client.GetErc20Info(contract)
		if err != nil {
			t.Fatal(err)
		}

		rawTx := &eth.RawTx{}
		rawTx.Init(client, params.TestChainConfig, gasPrice, gasTip, 100000, nonce+1, uint64(tokenInfo.Decimals), privKey, from, contract, to, amount)
		txid, err := rawTx.SendTx()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("tx send success, txid : ", txid)
	}
}
