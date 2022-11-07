package eth_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/gokch/blockchain_rpc/eth"
)

const (
	DEF_url_local_goeril  = "http://10.1.1.61:8546"
	DEF_url_remote_goeril = "https://goerli.infura.io/v3/c88b23ec7c6d4ec581e45f1f1a9afc9e"

	url string = DEF_url_remote_goeril

	DEF_private_key = "aa41425d6df6460c9dc413275830a152d5a9851713661f8c48f2494461bee885"
	DEF_address     = "0xe5bDa4eEd3FD91793632604B79cFc97372617eB4"
	DEF_token__BNB  = "0x64BBF67A8251F7482330C33E65b08B835125e018"
	DEF_token__KCH  = "0x48c8fb83907FcD67cA5F703658f2416630E3bA2a"
)

func Test_ServerInfo(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	info, err := client.GetServerInfo()
	if err != nil {
		t.Fatal("get server info fail |", err)
	}
	if info == nil {
		t.Fatal("server info not syncing")
	}
	fmt.Printf("current block : %v\n", info.CurrentBlock)
	fmt.Printf("highest block : %v\n", info.HighestBlock)
	fmt.Printf("highest block : %v\n", info.StartingBlock)
}

//-----------------------------------------------------------------------------//
// address

func Test_GetAddressBalance(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	balance, err := client.GetAddressBalance(DEF_address)
	if err != nil {
		t.Fatal("balance at fail |", err)
	}
	fmt.Printf("balance : %s\n", balance)
}

//-----------------------------------------------------------------------------//
// fee

func Test_FeeSuggestGas(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	price, tipcap, err := client.SuggestGasInfo()
	if err != nil {
		t.Fatal(err)
	}

	snPrice := price.String()
	snTipCap := tipcap.String()

	gasPriceGwei, err := eth.Conv_WeiToGwei(snPrice)
	if err != nil {
		t.Fatal(err)
	}

	gasTipCapGwei, err := eth.Conv_WeiToGwei(snTipCap)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("suggest gas price(gwei) : %v\nsuggest gas tip(gwei) : %v\n", gasPriceGwei, gasTipCapGwei)
}

//-----------------------------------------------------------------------------//
// tx

func Test_TxGetRawData(t *testing.T) {
	client := &eth.Client{}
	err := client.Open(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}
	blockInfo, err := client.GetBlockInfo(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	txInfos := blockInfo.Transactions()
	txLast := txInfos[blockInfo.Transactions().Len()-1]
	txRaw, _, err := client.GetTxInfo(txLast.Hash().Hex())
	if err != nil {
		t.Fatal("get tx fail |", err)
	}

	hexTxid, err := eth.EncodeTxRLP(txRaw)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	txRaw.EncodeRLP(buf)
	pt_tx_raw__new, err := eth.DecodeTxRLP(hexTxid)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(txRaw.Value().Bytes(), pt_tx_raw__new.Value().Bytes()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", txRaw.Value().Bytes(), pt_tx_raw__new.Value().Bytes())
	}
	if bytes.Equal(txRaw.Data(), pt_tx_raw__new.Data()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", txRaw.Data(), pt_tx_raw__new.Data())
	}
}

func Test_GetErc20Info(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// BNB
	info, err := client.GetErc20Info(DEF_token__BNB)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)

	// KCH
	info, err = client.GetErc20Info(DEF_token__KCH)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)
}

func Test_GetErc20TokenBalance(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	balance, err := client.GetErc20BalanceOf(DEF_address, DEF_token__KCH)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance)
}

func Test_GetErc20QueryAndReceipt(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	n_block_num, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}
	_ = n_block_num
	pt_block_info, err := client.GetBlockInfo(12895502)
	if err != nil {
		t.Fatal(err)
	}
	txs := pt_block_info.Transactions()
	blockHash := pt_block_info.Hash().String()
	fmt.Println(pt_block_info.Number())

	logs, err := client.FilterQueryTransfer(nil, blockHash)
	transfers, err := eth.DecodeTransfers(logs)
	for _, pt_transfer := range transfers {
		fmt.Println(pt_transfer.ContractAddr)
	}
	for _, pt_tx := range txs {
		isContract, err := client.ValidAddress(pt_tx.To().String())
		if err != nil {
			t.Fatal(err)
		}
		if isContract == true {
			receipt, err := client.GetTxReceipt(pt_tx.Hash().String())
			if err != nil {
				t.Fatal(err)
			} else if receipt == nil {
				t.Fatal("not exist receipt")
			}

			transfers, err := eth.DecodeTransfers(receipt.Logs)
			for _, pt_transfer := range transfers {
				fmt.Println(pt_transfer.ContractAddr)
			}
		} else {
			fmt.Println("it is not smart contract")
		}
	}
}

func Test_GetAddressFrom(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}
	blockInfo, err := client.GetBlockInfo(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	for _, tx := range blockInfo.Transactions() {
		from, err := eth.AddressGetSender(tx, params.RopstenChainConfig)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("tx / from : %v | %v\n", tx.Hash().String(), from)
	}
}

//-----------------------------------------------------------------------------//
// block

func Test_BlockGetInfo(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	n8_block_number, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", n8_block_number)
	pt_block_info, err := client.GetBlockInfo(n8_block_number)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(pt_block_info)
	pt_block_info.Bloom()

	fmt.Printf("block hash - %v\n", pt_block_info.TxHash())
	fmt.Printf("block number - %v\n", pt_block_info.NumberU64())
	fmt.Printf("tx len - %v\n", len(pt_block_info.Transactions()))

	for i, tx := range pt_block_info.Transactions() {
		fmt.Printf("\t%v\n", i)
		fmt.Printf("\ttxhash - %v\n", tx.Hash())
		fmt.Printf("\tnonce - %v\n", tx.Nonce())
		fmt.Printf("\tto - %v\n", tx.To())
		fmt.Printf("\tamount - %v\n", tx.Value().String())
	}
}

func Test_BlockEncode(t *testing.T) {
	client, err := eth.NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	n8_block_number, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}

	pt_block_info, err := client.GetBlockInfo(n8_block_number)
	if err != nil {
		t.Fatal(err)
	}

	bt_hex__block_hash, err := eth.EncodeBlockRLP(pt_block_info)
	if err != nil {
		t.Fatal(err)
	}

	pt_block_info__new, err := eth.DecodeBlockRLP(bt_hex__block_hash)
	if err != nil {
		t.Fatal(err)
	}
	bt_hex__block_hash__new, err := eth.EncodeBlockRLP(pt_block_info__new)
	if err != nil {
		t.Fatal(err)
	}
	// 검증
	if bytes.Equal(pt_block_info.Extra(), pt_block_info__new.Extra()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", pt_block_info.Extra(), pt_block_info__new.Extra())
	}
	if bytes.Equal(bt_hex__block_hash, bt_hex__block_hash__new) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", bt_hex__block_hash, bt_hex__block_hash__new)
	}
}
