package eth

import (
	"fmt"

	"github.com/gokch/snum"
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
//  basefee = Base Fee Per Gas
//  tipcap  = Max Priority Fee Per Gas
//  feecap  = Max Fee Per Gas
// ret
//  fee_burnt = Base Fee amount ( burnt )
//  fee_tip   = Tip Fee amount ( to block miner )
// 	fee_save  = Saving Fee amount ( refund )
func CalcFeeCost__DynamicFee(
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

		if feeSave, err = snFeeCap.GetStr(); err != nil {
			return "", "", "", err
		}
		if feeTip, err = snTipCap.GetStr(); err != nil {
			return "", "", "", err
		}
	} else { // tipcap 이 더 크면
		feeSave = "0"
		if feeTip, err = snFeeCap.GetStr(); err != nil {
			return "", "", "", err
		}
	}

	// get fee burnt
	if feeBurnt, err = snBaseFee.GetStr(); err != nil {
		return "", "", "", err
	}

	return feeBurnt, feeTip, feeSave, nil
}
