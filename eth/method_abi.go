package eth

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	token "github.com/gokch/blockchain_rpc/eth/smart_contract"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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

func (t *Client) FilterQueryTransfer(contractAddresses []string, blockHash string) (logs []*types.Log, err error) {
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
	typesLogs, err := t.rpc_client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, err
	}
	logs = make([]*types.Log, 0, len(typesLogs))
	for _, typesLog := range typesLogs {
		logs = append(logs, &typesLog)
	}
	return logs, nil
}

func (t *Client) MakeErc20TransferBytecode(addressTo string, amountTo *big.Int) (bytecodeTransfer []byte) {
	ethAddrTo := common.HexToAddress(addressTo)
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte("Transfer(address,address,uint256)"))

	methodId := hash.Sum(nil)[:4]
	addressPadded := common.LeftPadBytes(ethAddrTo.Bytes(), 32)
	amountPadded := common.LeftPadBytes(amountTo.Bytes(), 32)

	bytecodeTransfer = make([]byte, 0, len(methodId)+len(addressPadded)+len(amountPadded))
	bytecodeTransfer = append(bytecodeTransfer, methodId...)
	bytecodeTransfer = append(bytecodeTransfer, addressPadded...)
	bytecodeTransfer = append(bytecodeTransfer, amountPadded...)
	return bytecodeTransfer
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
	token, err := token.NewToken(ethContractAddr, t.rpc_client)
	if err != nil {
		return nil, err
	}

	info = &Erc20Info{}
	if token == nil {
		return info, nil
	}

	opts := &bind.CallOpts{}
	info.Name, err = token.Name(opts)
	if err != nil {
		return nil, err
	}

	info.Symbol, err = token.Symbol(opts)
	if err != nil {
		return nil, err
	}
	info.Decimals, err = token.Decimals(opts)
	if err != nil {
		return nil, err
	}
	totalSupply, err := token.TotalSupply(opts)
	if err != nil {
		return nil, err
	}
	info.TotalSupply = totalSupply.String()

	return info, nil
}

func (t *Client) GetErc20BalanceOf(addr string, contractAddr string) (balance string, err error) {
	// decimal 추출
	info, err := t.GetErc20Info(contractAddr)
	if err != nil {
		return "", err
	}

	token, err := token.NewToken(common.HexToAddress(contractAddr), t.rpc_client)
	if err != nil {
		return "", err
	}
	bigBalance, err := token.BalanceOf(&bind.CallOpts{}, common.HexToAddress(addr))
	if err != nil {
		return "", err
	}

	// decimal 에 따라 자릿수 변경
	balance, err = Conv_WeiToUnit(bigBalance.String(), info.Decimals)
	if err != nil {
		return "", err
	}
	return balance, nil
}

type TransferErc20 struct {
	From   string
	To     string
	Amount string

	Removed      bool
	BlockHash    string
	TxHash       string
	ContractAddr string
	Data         []byte
}

func DecodeTransfers(logs []*types.Log) (transfers []*TransferErc20, err error) {
	abi, err := abi.JSON(strings.NewReader(token.TokenABI))
	if err != nil {
		return nil, err
	}
	fnSigHash := crypto.Keccak256Hash([]byte("transfer(address,uint256)"))
	logSigHash := crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
	transfers = make([]*TransferErc20, 0, len(logs))
	for _, log := range logs {
		switch log.Topics[0] {
		case fnSigHash, logSigHash:
			if len(log.Data) == 0 {
				// erc 721 거래 등의 이유로 data 가 없을 수 있음
				continue
			}

			var transfer TransferErc20
			err = abi.UnpackIntoInterface(&transfer, "Transfer", log.Data)
			if err != nil {
				return nil, err
			}
			if len(log.Topics) < 3 {
				return nil, fmt.Errorf("invalid contract form (%v)", "Transfer")
			}
			transfer.From = log.Topics[1].Hex()
			transfer.To = log.Topics[2].Hex()
			transfer.BlockHash = log.BlockHash.Hex()
			transfer.TxHash = log.TxHash.Hex()
			transfer.ContractAddr = log.Address.Hex()
			transfer.Removed = log.Removed
			transfer.Data = log.Data
			transfers = append(transfers, &transfer)
		default:
			continue // 다른 method 일 경우
		}
	}
	return transfers, nil
}
