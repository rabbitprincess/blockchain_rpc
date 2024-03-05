package eth

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestServerInfo(t *testing.T) {
	client, err := NewClient(url)
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

func TestGetAddressBalance(t *testing.T) {
	client, err := NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	balance, err := client.GetAddressBalance(DEF_Address)
	if err != nil {
		t.Fatal("balance at fail |", err)
	}
	fmt.Printf("balance : %s\n", balance)
}

//-----------------------------------------------------------------------------//
// fee

func TestFeeSuggestGas(t *testing.T) {
	client, err := NewClient(url)
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

	gasPriceGwei, err := Conv_WeiToGwei(snPrice)
	if err != nil {
		t.Fatal(err)
	}

	gasTipCapGwei, err := Conv_WeiToGwei(snTipCap)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("suggest gas price(gwei) : %v\nsuggest gas tip(gwei) : %v\n", gasPriceGwei, gasTipCapGwei)
}

//-----------------------------------------------------------------------------//
// tx

func TestTxGetRawData(t *testing.T) {
	client := &Client{}
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

	hexTxid, err := EncodeTxRLP(txRaw)
	if err != nil {
		t.Fatal(err)
	}
	buf := bytes.NewBuffer(nil)
	txRaw.EncodeRLP(buf)
	pt_tx_raw__new, err := DecodeTxRLP(hexTxid)
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

func TestGetErc20Info(t *testing.T) {
	client, err := NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	// BNB
	info, err := client.GetErc20Info(DEF_TokenAERGO)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)

	// KCH
	info, err = client.GetErc20Info(DEF_TokenKCH)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)
}

func TestGetErc20TokenBalance(t *testing.T) {
	client, err := NewClient(url)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	balance, err := client.GetErc20BalanceOf(DEF_Address, DEF_TokenKCH)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(balance)
}

func TestGetErc20QueryAndReceipt(t *testing.T) {
	client, err := NewClient(url)
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
	logs, err := client.FilterLogs([]string{DEF_TokenAERGO}, blockHash)
	if err != nil {
		t.Fatal(err)
	}

	// decode logs ( transfer only )
	transfers, err := DecodeTransfers(logs)
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

func TestGetAddressFrom(t *testing.T) {
	client, err := NewClient(url)
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
		from, err := AddressGetSender(tx, params.RopstenChainConfig)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("tx / from : %v | %v\n", tx.Hash().String(), from)
	}
}

//-----------------------------------------------------------------------------//
// block

func TestBlockGetInfo(t *testing.T) {
	client, err := NewClient(url)
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

func TestBlockEncode(t *testing.T) {
	client, err := NewClient(url)
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

	blockRLP, err := EncodeBlockRLP(blockInfo)
	if err != nil {
		t.Fatal(err)
	}

	blockInfoNew, err := DecodeBlockRLP(blockRLP)
	if err != nil {
		t.Fatal(err)
	}
	blockHashNew, err := EncodeBlockRLP(blockInfoNew)
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
