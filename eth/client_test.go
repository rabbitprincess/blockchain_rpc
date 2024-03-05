package eth

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

const (
	DEF_url_local_goeril   = "http://10.1.1.61:8546"
	DEF_url_remote_goeril  = "https://goerli.infura.io/v3/c88b23ec7c6d4ec581e45f1f1a9afc9e"
	DEF_url_local_mainnet  = "https://eth-node.cccv.to:8545"
	DEF_url_remote_mainnet = "https://mainnet.infura.io/v3/c88b23ec7c6d4ec581e45f1f1a9afc9e"

	url string = DEF_url_remote_mainnet

	DEF_PrivKey    = "aa41425d6df6460c9dc413275830a152d5a9851713661f8c48f2494461bee885"
	DEF_Address    = "0xe5bDa4eEd3FD91793632604B79cFc97372617eB4"
	DEF_TokenBNB   = "0x64BBF67A8251F7482330C33E65b08B835125e018"
	DEF_TokenKCH   = "0x48c8fb83907FcD67cA5F703658f2416630E3bA2a"
	DEF_TokenAERGO = "0x91Af0fBB28ABA7E31403Cb457106Ce79397FD4E6"
)

func TestGetNewAddress(t *testing.T) {
	client, err := NewClient(url)
	require.NoError(t, err)

	defer client.Close()
	privKey, addr, err := client.GetNewAddress()
	require.NoError(t, err)

	fmt.Printf("private key : %v\n", privKey)
	fmt.Printf("address : %v\n", addr)
}

func TestWithdraw(t *testing.T) {
	// 출금 입력값
	var (
		privKey  = DEF_PrivKey
		from     = DEF_Address
		to       = "0xD76C201f700E5bAE854BD0722a8B29F87F9a9cCB"
		contract = DEF_TokenKCH
		amount   = "0.01"
	)

	client, err := NewClient(url)
	require.NoError(t, err)
	defer client.Close()

	// 출금 과정
	{
		gasPrice, gasTip, err := client.SuggestGasInfo()
		require.NoError(t, err)
		nonce, err := client.GetAddressNonce(from)
		require.NoError(t, err)
		tokenInfo, err := client.GetErc20Info(contract)
		require.NoError(t, err)

		rawTx := &RawTx{}
		rawTx.Init(client, params.TestChainConfig, gasPrice, gasTip, 100000, nonce+1, uint64(tokenInfo.Decimals), privKey, from, contract, to, amount)
		txid, err := rawTx.SendTx()
		require.NoError(t, err)
		fmt.Println("tx send success, txid : ", txid)
	}
}
