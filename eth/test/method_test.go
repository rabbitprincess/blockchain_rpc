package eth_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/gokch/blockchain_rpc/eth"
)

const (
	DEF_url_local_goeril   = "http://10.1.1.61:8546"
	DEF_url_remote_goeril  = "https://goerli.infura.io/v3/c88b23ec7c6d4ec581e45f1f1a9afc9e"
	DEF_url_local_mainnet  = "https://eth-node.cccv.to:8545"
	DEF_url_remote_mainnet = "https://mainnet.infura.io/v3/c88b23ec7c6d4ec581e45f1f1a9afc9e"

	url string = DEF_url_remote_mainnet

	DEF_private_key  = "aa41425d6df6460c9dc413275830a152d5a9851713661f8c48f2494461bee885"
	DEF_address      = "0xe5bDa4eEd3FD91793632604B79cFc97372617eB4"
	DEF_token__BNB   = "0x64BBF67A8251F7482330C33E65b08B835125e018"
	DEF_token__KCH   = "0x48c8fb83907FcD67cA5F703658f2416630E3bA2a"
	DEF_token__AERGO = "0x91Af0fBB28ABA7E31403Cb457106Ce79397FD4E6"
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
	info, err := client.GetErc20Info(DEF_token__AERGO)
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

	// get block info
	blockInfo, err := client.GetBlockInfo(15934533)
	if err != nil {
		t.Fatal(err)
	}
	blockHash := blockInfo.Hash().String()

	// get logs in block ( aergo token only )
	logs, err := client.FilterLogs([]string{DEF_token__AERGO}, blockHash)
	if err != nil {
		t.Fatal(err)
	}

	// decode logs ( transfer only )
	transfers, err := eth.DecodeTransfers(logs)
	if err != nil {
		t.Fatal(err)
	}

	for _, transfer := range transfers {
		fmt.Println("token  :", transfer.ContractAddr)
		fmt.Println("from   :", transfer.From)
		fmt.Println("to     :", transfer.To)
		fmt.Println("amount :", transfer.Amount)
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

	blockNumber, err := client.GetBlockMostRecent()
	if err != nil {
		t.Fatal(err)
	}

	blockInfo, err := client.GetBlockInfo(blockNumber)
	if err != nil {
		t.Fatal(err)
	}

	blockRLP, err := eth.EncodeBlockRLP(blockInfo)
	if err != nil {
		t.Fatal(err)
	}

	blockInfoNew, err := eth.DecodeBlockRLP(blockRLP)
	if err != nil {
		t.Fatal(err)
	}
	blockHashNew, err := eth.EncodeBlockRLP(blockInfoNew)
	if err != nil {
		t.Fatal(err)
	}
	// 검증
	if bytes.Equal(blockInfo.Extra(), blockInfoNew.Extra()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", blockInfo.Extra(), blockInfoNew.Extra())
	}
	if bytes.Equal(blockRLP, blockHashNew) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", blockRLP, blockHashNew)
	}
}
