package eth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitSize(t *testing.T) {
	// wei to eth
	WeiToEth, err := Conv_WeiToEth("1")
	require.NoError(t, err)
	assert.Equalf(t, WeiToEth, "0.000000000000000001", "invalid wei_to_eth | expect : 0.000000000000000001 | result : %s", WeiToEth)

	// wei to gwei
	WeiToGwei, err := Conv_WeiToGwei("1")
	require.NoError(t, err)
	assert.Equalf(t, WeiToGwei, "0.000000001", "invalid wei_to_gwei | expect : 0.000000001 | result : %s", WeiToGwei)

	// eth to wei
	EthToWei, err := Conv_EthToWei("1")
	require.NoError(t, err)
	assert.Equalf(t, EthToWei, "1000000000000000000", "invalid eth_to_wei | expect : 1000000000000000000 | result : %s", EthToWei)
}

// London HardFork(EIP-1559) 이후 변경된 수수료 체계(dynamic fee) 관련 테스트
func TestDynamicFee(t *testing.T) {
	for _, test := range []struct {
		inputGasUsed   uint64
		inputBaseFee   string // Base Fee Per Gas
		inputTipCap    string // Max Priorty Fee Per Gas
		inputFeeCap    string // Max Fee Per Gas
		expectFeeBurnt string
		expectFeeTip   string
		expectFeeSave  string
		errMsg         string
	}{
		// fee cap - base fee > tip cap 일 경우
		{1, "100", "10", "200", "100", "10", "90", ""},
		{100, "100", "10", "200", "10000", "1000", "9000", ""},
		// 0x69cf78dbdf0a77dac3d87314d09c03c3eb7e282d5953b73b28c69c04ead60b84
		{109245, "0.000000102752233976", "0.000000116551928074", "0.000000116551928074", "0.01122516780070812", "0.00150754758173601", "0", ""},
		// 0xae40090438da2b40d6589e7653d8372c730b685c52ed8eff8832595123c3f433
		{968546, "0.000000076679330156", "0.00000008", "0.00000008", "0.074267458505273176", "0.003216221494726824", "0", ""},
		// 0xe36a6a0861955ed894126a4ffb31ff758eda14dbb3e4443a26fa35b05e663059
		{968546, "0.000000076679330156", "0.000000086", "0.000000086", "0.074267458505273176", "0.009027497494726824", "0", ""},

		// fee cap - base fee < tip cap 일 경우
		{1, "100", "30", "120", "100", "20", "0", ""},
		{100, "100", "30", "120", "10000", "2000", "0", ""},
		// 0x6d19344a5f889dd5320edaa91f9bd08fe83af76960ffae982879bad45c84aab4
		{21000, "0.000000102992906016", "0.000000002", "0.000000219", "0.002162851026336", "0.000042", "0.002394148973664", ""},
		// 0xac461bc7e9afce4956b0ea9e904d0d11623e434add09ec5ec4d7fff1ac2b0e3f
		{21000, "0.00000010519268421", "0.0000000015", "0.000000160565775374", "0.00220904636841", "0.0000315", "0.001131334914444", ""},
		// 0x916540514841555e342cb4521df49ff73de4830cca14b9a86035c76ba8524754
		{21000, "0.00000010519268421", "0.000000002", "0.000000211", "0.00220904636841", "0.000042", "0.00217995363159", ""},
		// 0x916540514841555e342cb4521df49ff73de4830cca14b9a86035c76ba8524754
		{21000, "0.00000010519268421", "0.000000002", "0.000000211", "0.00220904636841", "0.000042", "0.00217995363159", ""},
		// 0x322183d429472221719db82462a832cacb288cd39033777b015c4b07effe6c67
		{21000, "0.000000083933307745", "0.000000001520462116", "0.000000141684830556", "0.001762599462645", "0.000031929704436", "0.001180852274595", ""},
		// 0xa9d2f086cf42c7f91521707d3533fdc53dafede07e415b720564dd7975b27742
		{105547, "0.000000083933307745", "0.000000002330847154", "0.000000091933832321", "0.008858908832561515", "0.000246013924563238", "0.000598417442859834", ""},
	} {
		feeBurnt, feeTip, feeSave, err := CalcFeeCost_DynamicFee(test.inputGasUsed, test.inputBaseFee, test.inputTipCap, test.inputFeeCap)
		if err != nil && err.Error() != test.errMsg {
			t.Errorf("\nerr is not same | output - %v | expect - %v", err, test.errMsg)
		}
		if feeBurnt != test.expectFeeBurnt {
			t.Errorf("\nfee burnt is not same | output - %v | expect - %v", feeBurnt, test.expectFeeBurnt)
		}
		if feeTip != test.expectFeeTip {
			t.Errorf("\nfee tip is not same | output - %v | expect - %v", feeTip, test.expectFeeTip)
		}
		if feeSave != test.expectFeeSave {
			t.Errorf("\nfee save is not same | output - %v | expect - %v", feeSave, test.expectFeeSave)
		}
	}
}

/*
func Test_txid__encode_decode(t *testing.T) {
	s_txid := "0xc9ec67d71b6a59eac2908ce8676c95c2df2e036f04bf0c30e7beeefc907e4d2b"
	bt_txid, err := Encode__txid(s_txid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(bt_txid)

	td_s_txid_copy, err := Decode__txid(bt_txid)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s_txid)
	fmt.Println(td_s_txid_copy)
}
*/
