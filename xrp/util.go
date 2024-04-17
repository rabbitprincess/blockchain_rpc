package xrp

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rabbitprincess/snum"
)

const (
	DEF_unitSize_drop = 1
	DEF_unitSize_xrp  = 6 // 1 xrp = 1000000 drop
)

func Conv_DropToXrp(drop string) (xrp string, err error) {
	snXrp := &snum.Snum{}
	err = snXrp.SetStr(drop)
	if err != nil {
		return "", err
	}

	sn10 := &snum.Snum{}
	sn10.SetUint64(10)
	sn10.Pow(DEF_unitSize_xrp)
	snXrp.Div(sn10)

	return snXrp.String(), nil
}

func Conv_XrpToDrop(xrp string) (drop string, err error) {
	snDrop := &snum.Snum{}
	err = snDrop.SetStr(xrp)
	if err != nil {
		return "", err
	}

	sn10 := &snum.Snum{}
	sn10.SetUint64(10)
	sn10.Pow(DEF_unitSize_xrp)
	snDrop.Mul(sn10)

	return snDrop.String(), nil
}

func EncodeTxid(s string) (bt []byte, err error) {
	txHash := common.HexToHash(s)
	return txHash.Bytes(), nil
}

func DecodeTxid(bt []byte) (s string, err error) {
	txHash := common.BytesToHash(bt)
	s = txHash.Hex()

	// 후처리 - xrp 스펙에서 0x 를 붙이면 조회 불가능해 0x 를 제거한다
	if s[0:2] == "0x" {
		s = s[2:]
	}
	return s, nil
}
