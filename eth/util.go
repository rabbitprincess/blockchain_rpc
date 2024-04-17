package eth

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rabbitprincess/snum_sort/snum"
)

const (
	UnitSize_wei  = 1
	UnitSize_gwei = 9
	UnitSize_eth  = 18

	DEF_s_zero_address         = "0x0000000000000000000000000000000000000000"
	DEF_s_regex__valid_address = "^0x[0-9a-fA-F]{40}$"
)

//----------------------------------------------------------------------------------------//
// conv unit size

func Conv_WeiToEth(wei string) (eth string, err error) {
	return Conv_WeiToUnit(wei, UnitSize_eth)
}

func Conv_EthToWei(eth string) (wei string, err error) {
	return Conv_UnitToWei(eth, UnitSize_eth)
}

func Conv_WeiToGwei(wei string) (gwei string, err error) {
	return Conv_WeiToUnit(wei, UnitSize_gwei)
}

func Conv_GweiToWei(gwei string) (wei string, err error) {
	return Conv_UnitToWei(gwei, UnitSize_gwei)
}

func Conv_WeiToUnit(wei string, decimal uint8) (unit string, err error) {
	sn := &snum.Snum{}
	err = sn.SetStr(wei)
	if err != nil {
		return "", err
	}

	// decimal -> decimal size ( 10 ** _u1_decimals )
	snDecimal := &snum.Snum{}
	snDecimal.SetUint64(10)
	snDecimal.Pow(int64(decimal))
	sn.Div(snDecimal)

	return sn.String(), nil
}

func Conv_UnitToWei(unit string, decimal uint8) (wei string, err error) {
	sn := &snum.Snum{}
	err = sn.SetStr(unit)
	if err != nil {
		return "", err
	}

	// get decimal size ( 10 ** _u1_decimals )
	snDecimal := &snum.Snum{}
	snDecimal.SetUint64(10)
	snDecimal.Pow(int64(decimal))
	sn.Mul(snDecimal)

	return sn.String(), nil
}

// get Dynamic Fee by type 2 Transaction ( After London HardFork, EIP 1559 )
// arg
//
//	basefee = Base Fee Per Gas
//	tipcap  = Max Priority Fee Per Gas
//	feecap  = Max Fee Per Gas
//
// ret
//
//	feeBurnt = Base Fee amount ( burnt )
//	feeTip   = Tip Fee amount ( to block miner )
//	feeSave  = Saving Fee amount ( refund )
func CalcFeeCost_DynamicFee(
	gasUsed uint64,
	baseFee string,
	gasTipCap string,
	gasFeeCap string,
) (
	feeBurnt string,
	feeTip string,
	feeSave string,
	err error,
) {
	// snum 형식으로 변환
	snBaseFee := &snum.Snum{}
	snTipCap := &snum.Snum{}
	snFeeCap := &snum.Snum{}
	{
		if err = snBaseFee.SetStr(baseFee); err != nil {
			return "", "", "", err
		}
		if err = snTipCap.SetStr(gasTipCap); err != nil {
			return "", "", "", err
		}
		if err = snFeeCap.SetStr(gasFeeCap); err != nil {
			return "", "", "", err
		}
	}

	// gas_used 곱하기
	snGasUsed := &snum.Snum{}
	snGasUsed.SetUint64(gasUsed)
	snBaseFee.Mul(snGasUsed)
	snTipCap.Mul(snGasUsed)
	snFeeCap.Mul(snGasUsed)

	// basefee 와 feecap 비교
	if snFeeCap.Cmp(snBaseFee) < 0 { // feecap 이 더 작으면 에러 처리
		return "", "", "", fmt.Errorf("not enough fee cap | basefee - %v | feecap - %v", baseFee, gasFeeCap)
	}

	// feecap 에서 base fee 차감
	snFeeCap.Sub(snBaseFee)

	// (feecap - basefee) 와 tipcap 을 비교
	if snFeeCap.Cmp(snTipCap) > 0 { // feecap - basefee 가 더 크면
		snFeeCap.Sub(snTipCap)

		feeSave = snFeeCap.String()
		feeTip = snTipCap.String()
	} else { // tipcap 이 더 크면
		feeSave = "0"
		feeTip = snFeeCap.String()
	}

	// get fee burnt
	if feeBurnt = snBaseFee.GetStr(); feeBurnt == "" {
		return "", "", "", err
	}

	return feeBurnt, feeTip, feeSave, nil
}

//----------------------------------------------------------------------------------------//
// rlp encode, decode

func EncodeBlockRLP(block *types.Block) ([]byte, error) {
	bt_bytes_buffer := bytes.NewBuffer(nil)
	err := block.EncodeRLP(bt_bytes_buffer)
	if err != nil {
		return nil, err
	}
	return bt_bytes_buffer.Bytes(), nil
}

func DecodeBlockRLP(blockRLP []byte) (*types.Block, error) {
	bt_bytes_buffer := bytes.NewBuffer(blockRLP)
	block := &types.Block{}
	err := block.DecodeRLP(rlp.NewStream(bt_bytes_buffer, 0))
	if err != nil {
		return nil, err
	}
	return block, nil
}

func EncodeTxRLP(tx *types.Transaction) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := tx.EncodeRLP(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeTxRLP(txRLP []byte) (*types.Transaction, error) {
	buf := bytes.NewBuffer(txRLP)
	tx := &types.Transaction{}
	err := tx.DecodeRLP(rlp.NewStream(buf, 0))
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func EncodeReceiptRLP(receipt *types.Receipt) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := receipt.EncodeRLP(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeReceiptRLP(receiptRLP []byte) (*types.Receipt, error) {
	buf := bytes.NewBuffer(receiptRLP)
	receipt := &types.Receipt{}
	err := receipt.DecodeRLP(rlp.NewStream(buf, 0))
	if err != nil {
		return nil, err
	}
	return receipt, nil
}

func EncodeTxid(hexTxid string) (rawTxid []byte, err error) {
	pt_tx_hash := common.HexToHash(hexTxid)
	return pt_tx_hash.Bytes(), nil
}

func DecodeTxid(rawTxid []byte) (hexTxid string, err error) {
	pt_tx_hash := common.BytesToHash(rawTxid)
	return pt_tx_hash.Hex(), nil
}

//----------------------------------------------------------------------------------------//
// address

func AddressValid(addr interface{}) bool {
	re := regexp.MustCompile(DEF_s_regex__valid_address)
	switch v := addr.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

func AddressIsZero(addr interface{}) bool {
	var address common.Address
	switch v := addr.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	bt_address_zero := common.FromHex(DEF_s_zero_address)
	bt_address := address.Bytes()
	return reflect.DeepEqual(bt_address, bt_address_zero)
}

func AddressChecksum(address string) (addressChecksum string) {
	pt_wallet_addr := common.HexToAddress(address)
	return pt_wallet_addr.Hex()
}

func AddressGetSender(tx *types.Transaction, config *params.ChainConfig) (from string, err error) {
	addr, err := types.Sender(types.LatestSigner(config), tx)
	if err != nil {
		return "", err
	}
	return addr.Hex(), nil
}

//----------------------------------------------------------------------------------------//
// genesis block

func GetGenesis(config *params.ChainConfig) (genesis *core.Genesis, err error) {
	switch config.ChainID {
	case params.MainnetChainConfig.ChainID:
		genesis = core.DefaultGenesisBlock()
	case params.RopstenChainConfig.ChainID:
		genesis = core.DefaultRopstenGenesisBlock()
	case params.GoerliChainConfig.ChainID:
		genesis = core.DefaultGoerliGenesisBlock()
	case params.RinkebyChainConfig.ChainID:
		genesis = core.DefaultRinkebyGenesisBlock()
	default:
		return nil, fmt.Errorf("invalid chain config | chainID : %v", config.ChainID)
	}

	return genesis, nil
}
