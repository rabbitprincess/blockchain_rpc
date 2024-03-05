package eth

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	defer client.Close()

	balance, err := client.GetAddressBalance(DEF_Address)
	require.NoError(t, err, "balance at failed")
	fmt.Printf("balance : %s\n", balance)
}

//-----------------------------------------------------------------------------//
// fee

func TestFeeSuggestGas(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	price, tipcap, err := client.SuggestGasInfo()
	require.NoError(t, err)

	gasPriceGwei, err := Conv_WeiToGwei(price.String())
	require.NoError(t, err)

	gasTipCapGwei, err := Conv_WeiToGwei(tipcap.String())
	require.NoError(t, err)

	fmt.Printf("suggest gas price(gwei) : %v\nsuggest gas tip(gwei) : %v\n", gasPriceGwei, gasTipCapGwei)
}

//-----------------------------------------------------------------------------//
// tx

func TestTxGetRawData(t *testing.T) {
	client := &Client{}
	err := client.Open(url)
	require.NoError(t, err)
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	require.NoError(t, err)

	blockInfo, err := client.GetBlockInfo(blockNumber)
	require.NoError(t, err)

	txInfos := blockInfo.Transactions()
	txLast := txInfos[blockInfo.Transactions().Len()-1]
	txRaw, _, err := client.GetTxInfo(txLast.Hash().Hex())
	require.NoError(t, err, "get tx failed")

	hexTxid, err := EncodeTxRLP(txRaw)
	require.NoError(t, err)

	buf := bytes.NewBuffer(nil)
	txRaw.EncodeRLP(buf)
	newTx, err := DecodeTxRLP(hexTxid)
	require.NoError(t, err)

	if bytes.Equal(txRaw.Value().Bytes(), newTx.Value().Bytes()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", txRaw.Value().Bytes(), newTx.Value().Bytes())
	}
	if bytes.Equal(txRaw.Data(), newTx.Data()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", txRaw.Data(), newTx.Data())
	}
}

func TestGetErc20Info(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	// BNB
	info, err := client.GetErc20Info(DEF_TokenAERGO)
	require.NoError(t, err)
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)

	// KCH
	info, err = client.GetErc20Info(DEF_TokenKCH)
	require.NoError(t, err)
	fmt.Println(info.IsFunded, info.Name, info.Symbol, info.TotalSupply)
}

func TestGetErc20TokenBalance(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	balance, err := client.GetErc20BalanceOf(DEF_Address, DEF_TokenKCH)
	require.NoError(t, err)

	fmt.Println(balance)
}

func TestGetErc20QueryAndReceipt(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	// get block info
	blockInfo, err := client.GetBlockInfo(15934533)
	require.NoError(t, err)
	blockHash := blockInfo.Hash().String()

	// get logs in block ( aergo token only )
	logs, err := client.FilterLogs([]string{DEF_TokenAERGO}, blockHash)
	require.NoError(t, err)

	// decode logs ( transfer only )
	transfers, err := DecodeTransfers(logs)
	require.NoError(t, err)

	for _, transfer := range transfers {
		fmt.Println("token  :", transfer.ContractAddr)
		fmt.Println("from   :", transfer.From)
		fmt.Println("to     :", transfer.To)
		fmt.Println("amount :", transfer.Amount)
	}
}

func TestGetAddressFrom(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	require.NoError(t, err)

	blockInfo, err := client.GetBlockInfo(blockNumber)
	require.NoError(t, err)

	for _, tx := range blockInfo.Transactions() {
		from, err := AddressGetSender(tx, params.RopstenChainConfig)
		require.NoError(t, err)
		fmt.Printf("tx / from : %v | %v\n", tx.Hash().String(), from)
	}
}

//-----------------------------------------------------------------------------//
// block

func TestBlockGetInfo(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	require.NoError(t, err)

	fmt.Printf("%v\n", blockNumber)
	blockInfo, err := client.GetBlockInfo(blockNumber)
	require.NoError(t, err)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("block hash - %v\n", blockInfo.TxHash())
	fmt.Printf("block number - %v\n", blockInfo.NumberU64())
	fmt.Printf("tx len - %v\n", len(blockInfo.Transactions()))

	for i, tx := range blockInfo.Transactions() {
		fmt.Printf("\t%v\n", i)
		fmt.Printf("\ttxhash - %v\n", tx.Hash())
		fmt.Printf("\tnonce - %v\n", tx.Nonce())
		fmt.Printf("\tto - %v\n", tx.To())
		fmt.Printf("\tamount - %v\n", tx.Value().String())
	}
}

func TestBlockEncode(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	blockNumber, err := client.GetBlockMostRecent()
	require.NoError(t, err)

	blockInfo, err := client.GetBlockInfo(blockNumber)
	require.NoError(t, err)

	blockRLP, err := EncodeBlockRLP(blockInfo)
	require.NoError(t, err)

	blockInfoNew, err := DecodeBlockRLP(blockRLP)
	require.NoError(t, err)

	blockHashNew, err := EncodeBlockRLP(blockInfoNew)
	require.NoError(t, err)

	// 검증
	if bytes.Equal(blockInfo.Extra(), blockInfoNew.Extra()) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", blockInfo.Extra(), blockInfoNew.Extra())
	}
	if bytes.Equal(blockRLP, blockHashNew) != true {
		t.Fatalf("encode decode error\n\tori : %v\n\tnew : %v\n", blockRLP, blockHashNew)
	}
}
