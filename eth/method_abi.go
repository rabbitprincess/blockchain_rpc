package eth

import (
	"context"
	"math/big"

	token "blockchain_rpc/eth/smart_contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

type ApprovalLog struct {
	TokenOwner common.Address
	Spender    common.Address
	Tokens     *big.Int
}

func (t *Client) FilterQueryTransfer(contractAddresses []string, blockHash string) (logs []types.Log, err error) {
	logTransferSigHash := crypto.Keccak256Hash([]byte("transfer(address,uint256)"))
	ethBlockHash := common.HexToHash(blockHash)

	var ethContractAddresses []common.Address
	if len(contractAddresses) > 0 {
		ethContractAddresses = make([]common.Address, 0, len(contractAddresses))
		for _, address := range contractAddresses {
			ethContractAddresses = append(ethContractAddresses, common.HexToAddress(address))
		}
	}

	query := ethereum.FilterQuery{
		BlockHash: &ethBlockHash,
		Addresses: ethContractAddresses,
		Topics:    [][]common.Hash{{logTransferSigHash}},
	}
	return t.rpc_client.FilterLogs(context.Background(), query)
}

func (t *Client) MakeErc20TransferBytecode(addressTo string, amountTo *big.Int) (bytecodeTransfer []byte, err error) {
	ethAddrTo := common.HexToAddress(addressTo)
	hash := sha3.NewLegacyKeccak256()
	_, err = hash.Write([]byte("Transfer(address,address,uint256)"))
	if err != nil {
		return nil, err
	}
	methodId := hash.Sum(nil)[:4]
	addressPadded := common.LeftPadBytes(ethAddrTo.Bytes(), 32)
	amountPadded := common.LeftPadBytes(amountTo.Bytes(), 32)

	bytecodeTransfer = make([]byte, 0, len(methodId)+len(addressPadded)+len(amountPadded))
	bytecodeTransfer = append(bytecodeTransfer, methodId...)
	bytecodeTransfer = append(bytecodeTransfer, addressPadded...)
	bytecodeTransfer = append(bytecodeTransfer, amountPadded...)
	return bytecodeTransfer, nil
}

type Erc20Info struct {
	IsFunded    bool
	Name        string
	Symbol      string
	Decimals    uint8
	TotalSupply string
}

func (t *Client) GetErc20Info(contractAddr string) (info *Erc20Info, err error) {
	ethContractAddr := common.HexToAddress(contractAddr)
	pt_token, err := token.NewToken(ethContractAddr, t.rpc_client)
	if err != nil {
		return nil, err
	}

	info = &Erc20Info{}
	if pt_token == nil {
		return info, nil
	}

	opts := &bind.CallOpts{}
	info.Name, err = pt_token.Name(opts)
	if err != nil {
		return nil, err
	}

	info.Symbol, err = pt_token.Symbol(opts)
	if err != nil {
		return nil, err
	}
	info.Decimals, err = pt_token.Decimals(opts)
	if err != nil {
		return nil, err
	}
	totalSupply, err := pt_token.TotalSupply(opts)
	if err != nil {
		return nil, err
	}
	info.TotalSupply = totalSupply.String()

	return info, nil
}

func (t *Client) GetErc20BalanceOf(address string, contractAddr string) {

}
