package eth_test

import (
	"testing"

	"github.com/gokch/blockchain_rpc/eth"
)

func TestUnitSize(t *testing.T) {
	// wei to eth
	WeiToEth, _ := eth.Conv_WeiToEth("1")
	if WeiToEth != "0.000000000000000001" {
		t.Errorf("invalid wei_to_eth | expect : 0.000000000000000001 | result : %s", WeiToEth)
	}

	// wei to gwei
	WeiToGwei, _ := eth.Conv_WeiToGwei("1")
	if WeiToGwei != "0.000000001" {
		t.Errorf("invalid wei_to_gwei | expect : 0.000000001 | result : %s", WeiToGwei)
	}

	// eth to wei
	EthToWei, _ := eth.Conv_EthToWei("1")
	if EthToWei != "1000000000000000000" {
		t.Errorf("invalid eth_to_wei | expect : 1000000000000000000 | result : %s", EthToWei)
	}
}

// London HardFork(EIP-1559) 이후 변경된 수수료 체계(dynamic fee) 관련 테스트
func TestDynamicFee(t *testing.T) {
	type TestInput struct {
		inputGasUsed   uint64
		inputBaseFee   string // Base Fee Per Gas
		inputTipCap    string // Max Priorty Fee Per Gas
		inputFeeCap    string // Max Fee Per Gas
		expectFeeBurnt string
		expectFeeTip   string
		expectFeeSave  string
		errMsg         string
	}

	fn := func(_t_input TestInput) {
		feeBurnt, feeTip, feeSave, err := eth.CalcFeeCost_DynamicFee(_t_input.inputGasUsed, _t_input.inputBaseFee, _t_input.inputTipCap, _t_input.inputFeeCap)
		if err != nil && err.Error() != _t_input.errMsg {
			t.Errorf("\nerr is not same | output - %v | expect - %v", err, _t_input.errMsg)
		}
		if feeBurnt != _t_input.expectFeeBurnt {
			t.Errorf("\nfee burnt is not same | output - %v | expect - %v", feeBurnt, _t_input.expectFeeBurnt)
		}
		if feeTip != _t_input.expectFeeTip {
			t.Errorf("\nfee tip is not same | output - %v | expect - %v", feeTip, _t_input.expectFeeTip)
		}
		if feeSave != _t_input.expectFeeSave {
			t.Errorf("\nfee save is not same | output - %v | expect - %v", feeSave, _t_input.expectFeeSave)
		}
	}

	// base test
	{
		// fee cap - base fee > tip cap 일 경우
		fn(TestInput{
			inputGasUsed:   1,
			inputBaseFee:   "100",
			inputTipCap:    "10",
			inputFeeCap:    "200",
			expectFeeBurnt: "100",
			expectFeeTip:   "10",
			expectFeeSave:  "90",
			errMsg:         ""},
		)
		fn(TestInput{
			inputGasUsed:   100,
			inputBaseFee:   "100",
			inputTipCap:    "10",
			inputFeeCap:    "200",
			expectFeeBurnt: "10000",
			expectFeeTip:   "1000",
			expectFeeSave:  "9000",
			errMsg:         ""},
		)

		// fee cap - base fee < tip cap 일 경우
		fn(TestInput{
			inputGasUsed:   1,
			inputBaseFee:   "100",
			inputTipCap:    "30",
			inputFeeCap:    "120",
			expectFeeBurnt: "100",
			expectFeeTip:   "20",
			expectFeeSave:  "0",
			errMsg:         ""},
		)
		fn(TestInput{
			inputGasUsed:   100,
			inputBaseFee:   "100",
			inputTipCap:    "30",
			inputFeeCap:    "120",
			expectFeeBurnt: "10000",
			expectFeeTip:   "2000",
			expectFeeSave:  "0",
			errMsg:         ""},
		)
	}

	// 실제 tx 예제
	{
		// fee cap - base fee > tip cap 일 경우
		{
			// 0x69cf78dbdf0a77dac3d87314d09c03c3eb7e282d5953b73b28c69c04ead60b84
			fn(TestInput{
				inputGasUsed:   109245,
				inputBaseFee:   "0.000000102752233976",
				inputTipCap:    "0.000000116551928074",
				inputFeeCap:    "0.000000116551928074",
				expectFeeBurnt: "0.01122516780070812",
				expectFeeTip:   "0.00150754758173601",
				expectFeeSave:  "0",
				errMsg:         ""},
			)

			// 0xae40090438da2b40d6589e7653d8372c730b685c52ed8eff8832595123c3f433
			fn(TestInput{
				inputGasUsed:   968546,
				inputBaseFee:   "0.000000076679330156",
				inputTipCap:    "0.00000008",
				inputFeeCap:    "0.00000008",
				expectFeeBurnt: "0.074267458505273176",
				expectFeeTip:   "0.003216221494726824",
				expectFeeSave:  "0",
				errMsg:         ""},
			)

			// 0xe36a6a0861955ed894126a4ffb31ff758eda14dbb3e4443a26fa35b05e663059
			fn(TestInput{
				inputGasUsed:   968546,
				inputBaseFee:   "0.000000076679330156",
				inputTipCap:    "0.000000086",
				inputFeeCap:    "0.000000086",
				expectFeeBurnt: "0.074267458505273176",
				expectFeeTip:   "0.009027497494726824",
				expectFeeSave:  "0",
				errMsg:         ""},
			)
		}

		// fee cap - base fee < tip cap 일 경우
		{
			// 0x6d19344a5f889dd5320edaa91f9bd08fe83af76960ffae982879bad45c84aab4
			fn(TestInput{
				inputGasUsed:   21000,
				inputBaseFee:   "0.000000102992906016",
				inputTipCap:    "0.000000002",
				inputFeeCap:    "0.000000219",
				expectFeeBurnt: "0.002162851026336",
				expectFeeTip:   "0.000042",
				expectFeeSave:  "0.002394148973664",
				errMsg:         ""},
			)

			// 0xac461bc7e9afce4956b0ea9e904d0d11623e434add09ec5ec4d7fff1ac2b0e3f
			fn(TestInput{
				inputGasUsed:   21000,
				inputBaseFee:   "0.00000010519268421",
				inputTipCap:    "0.0000000015",
				inputFeeCap:    "0.000000160565775374",
				expectFeeBurnt: "0.00220904636841",
				expectFeeTip:   "0.0000315",
				expectFeeSave:  "0.001131334914444",
				errMsg:         ""},
			)

			// 0x916540514841555e342cb4521df49ff73de4830cca14b9a86035c76ba8524754
			fn(TestInput{
				inputGasUsed:   21000,
				inputBaseFee:   "0.00000010519268421",
				inputTipCap:    "0.000000002",
				inputFeeCap:    "0.000000211",
				expectFeeBurnt: "0.00220904636841",
				expectFeeTip:   "0.000042",
				expectFeeSave:  "0.00217995363159",
				errMsg:         ""},
			)

			// 0x916540514841555e342cb4521df49ff73de4830cca14b9a86035c76ba8524754
			fn(TestInput{
				inputGasUsed:   21000,
				inputBaseFee:   "0.00000010519268421",
				inputTipCap:    "0.000000002",
				inputFeeCap:    "0.000000211",
				expectFeeBurnt: "0.00220904636841",
				expectFeeTip:   "0.000042",
				expectFeeSave:  "0.00217995363159",
				errMsg:         ""},
			)

			// 0x322183d429472221719db82462a832cacb288cd39033777b015c4b07effe6c67
			fn(TestInput{
				inputGasUsed:   21000,
				inputBaseFee:   "0.000000083933307745",
				inputTipCap:    "0.000000001520462116",
				inputFeeCap:    "0.000000141684830556",
				expectFeeBurnt: "0.001762599462645",
				expectFeeTip:   "0.000031929704436",
				expectFeeSave:  "0.001180852274595",
				errMsg:         ""},
			)

			// 0xa9d2f086cf42c7f91521707d3533fdc53dafede07e415b720564dd7975b27742
			fn(TestInput{
				inputGasUsed:   105547,
				inputBaseFee:   "0.000000083933307745",
				inputTipCap:    "0.000000002330847154",
				inputFeeCap:    "0.000000091933832321",
				expectFeeBurnt: "0.008858908832561515",
				expectFeeTip:   "0.000246013924563238",
				expectFeeSave:  "0.000598417442859834",
				errMsg:         ""},
			)
		}
	}
}

/*
func Test_txid__encode_decode(t *testing.T) {
	s_txid := "0xc9ec67d71b6a59eac2908ce8676c95c2df2e036f04bf0c30e7beeefc907e4d2b"
	bt_txid, err := eth.Encode__txid(s_txid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(bt_txid)

	td_s_txid_copy, err := eth.Decode__txid(bt_txid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s_txid)
	fmt.Println(td_s_txid_copy)
}
*/
