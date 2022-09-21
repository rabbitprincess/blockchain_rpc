package types

import "encoding/json"

type Hash128 [16]byte
type Hash160 [20]byte
type Hash256 [32]byte
type PrivateKey []byte
type PublicKey [33]byte
type RegularKey [20]byte
type Seed []byte

//--------------------------------------------------------------------------//
// key

type TD_s_key_type string

const (
	TD_s_key_type__secp256k1 TD_s_key_type = "secp256k1"
	TD_s_key_type__ed25519   TD_s_key_type = "ed25519"
)

//--------------------------------------------------------------------------//
// currency

type TD_u1_currecncy_type uint8

const (
	TD_u1__XRP       TD_u1_currecncy_type = 0
	TD_u1__STANDARD  TD_u1_currecncy_type = 1
	TD_u1__DEMURRAGE TD_u1_currecncy_type = 2
	TD_u1__HEX       TD_u1_currecncy_type = 3
	TD_u1__UNKNOWN   TD_u1_currecncy_type = 4
)

//--------------------------------------------------------------------------//
// drop

type Drops struct {
	BaseFee       string `json:"base_fee"`
	MedianFee     string `json:"median_fee"`
	MinimumFee    string `json:"minimum_fee"`
	OpenLedgerFee string `json:"open_ledger_fee"`
}

//--------------------------------------------------------------------------//
// level

type Levels struct {
	MedianLevel     string `json:"median_level"`
	MinimumLevel    string `json:"minimum_level"`
	OpenLedgerLevel string `json:"open_ledger_level"`
	ReferenceLevel  string `json:"reference_level"`
}

//--------------------------------------------------------------------------//
// amount

type Amount struct {
	Value    string
	Currency string
	Issuer   string
}

func (t *Amount) UnmarshalJSON(bt []byte) (err error) {
	if len(bt) == 0 {
		return nil
	}
	if bt[0] != '{' {
		err = json.Unmarshal(bt, &t.Value)
		if err != nil {
			return err
		}
		return nil
	}
	dummy := struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
		Issuer   string `json:"issuer"`
	}{}
	err = json.Unmarshal(bt, &dummy)
	if err != nil {
		return err
	}
	t.Value = dummy.Value
	t.Currency = dummy.Currency
	t.Issuer = dummy.Issuer
	return nil
}

type PathElem struct {
	Account  *string
	Currency *string
	Issuer   *string
}
type Memo struct {
	MemoType   []byte
	MemoData   []byte
	MemoFormat []byte
}

//--------------------------------------------------------------------------//
// wallet

type Wallet struct {
	S_private_key    string `json:"privateKey"`
	S_public_key     string `json:"publicKey"`
	S_wallet_address string `json:"classicAddress"`
	S_seed           string `json:"seed"`
}

//--------------------------------------------------------------------------//
// tx type
type TxType string

const (
	Tx_Payment              TxType = "Payment"
	Tx_AccountSet           TxType = "AccountSet"
	Tx_AccountDelete        TxType = "AccountDelete"
	Tx_SetRegularKey        TxType = "SetRegularKey"
	Tx_OfferCreate          TxType = "OfferCreate"
	Tx_OfferCancel          TxType = "OfferCancel"
	Tx_TrustSet             TxType = "TrustSet"
	Tx_EnableAmendment      TxType = "EnableAmendment"
	Tx_SetFee               TxType = "SetFee"
	Tx_UNLModify            TxType = "UNLModify"
	Tx_TicketCreate         TxType = "TicketCreate"
	Tx_EscrowCreate         TxType = "EscrowCreate"
	Tx_EscrowFinish         TxType = "EscrowFinish"
	Tx_EscrowCancel         TxType = "EscrowCancel"
	Tx_SignerListSet        TxType = "SignerListSet"
	Tx_PaymentChannelCreate TxType = "PaymentChannelCreate"
	Tx_PaymentChannelFund   TxType = "PaymentChannelFund"
	Tx_PaymentChannelClaim  TxType = "PaymentChannelClaim"
	Tx_CheckCreate          TxType = "CheckCreate"
	Tx_CheckCash            TxType = "CheckCash"
	Tx_CheckCancel          TxType = "CheckCancel"
)

//--------------------------------------------------------------------------//
// tx flag
type TxFlag uint32

const (
	// Universal flags
	Tx_CanonicalSignature TxFlag = 0x80000000

	// Payment flags
	Tx_NoDirectRipple TxFlag = 0x00010000
	Tx_PartialPayment TxFlag = 0x00020000
	Tx_LimitQuality   TxFlag = 0x00040000
	Tx_Circle         TxFlag = 0x00080000 // Not implemented

	// AccountSet flags
	Tx_SetRequireDest   TxFlag = 0x00000001
	Tx_SetRequireAuth   TxFlag = 0x00000002
	Tx_SetDisallowXRP   TxFlag = 0x00000003
	Tx_SetDisableMaster TxFlag = 0x00000004
	Tx_SetAccountTxnID  TxFlag = 0x00000005
	Tx_NoFreeze         TxFlag = 0x00000006
	Tx_GlobalFreeze     TxFlag = 0x00000007
	Tx_DefaultRipple    TxFlag = 0x00000008
	Tx_RequireDestTag   TxFlag = 0x00010000
	Tx_OptionalDestTag  TxFlag = 0x00020000
	Tx_RequireAuth      TxFlag = 0x00040000
	Tx_OptionalAuth     TxFlag = 0x00080000
	Tx_DisallowXRP      TxFlag = 0x00100000
	Tx_AllowXRP         TxFlag = 0x00200000

	// OfferCreate flags
	Tx_Passive           TxFlag = 0x00010000
	Tx_ImmediateOrCancel TxFlag = 0x00020000
	Tx_FillOrKill        TxFlag = 0x00040000
	Tx_Sell              TxFlag = 0x00080000

	// TrustSet flags
	Tx_SetAuth       TxFlag = 0x00010000
	Tx_SetNoRipple   TxFlag = 0x00020000
	Tx_ClearNoRipple TxFlag = 0x00040000
	Tx_SetFreeze     TxFlag = 0x00100000
	Tx_ClearFreeze   TxFlag = 0x00200000

	// EnableAmendments flags
	Tx_GotMajority  TxFlag = 0x00010000
	Tx_LostMajority TxFlag = 0x00020000

	// PaymentChannelClaim flags
	Tx_Renew TxFlag = 0x00010000
	Tx_Close TxFlag = 0x00020000
)

//---------------------------------------------------------------------------//
// tx metadata

type Metadata struct {
	AffectedNodes     []AffectedNodes `json:"AffectedNodes"`
	TransactionIndex  int             `json:"TransactionIndex"`
	TransactionResult string          `json:"TransactionResult"`
	DeliveredAmount   Amount          `json:"delivered_amount"`
}

type AffectedNodes struct {
	ModifiedNode ModifiedNode `json:"ModifiedNode"`
}

type ModifiedNode struct {
	FinalFields       FinalFields    `json:"FinalFields"`
	LedgerEntryType   string         `json:"LedgerEntryType"`
	LedgerIndex       string         `json:"LedgerIndex"`
	PreviousFields    PreviousFields `json:"PreviousFields"`
	PreviousTxnID     string         `json:"PreviousTxnID"`
	PreviousTxnLgrSeq int            `json:"PreviousTxnLgrSeq"`
}

type FinalFields struct {
	Account    string      `json:"Account"`
	Balance    interface{} `json:"Balance"`
	Flags      TxFlag      `json:"Flags"`
	OwnerCount int         `json:"OwnerCount"`
	Sequence   int         `json:"Sequence"`
}

type PreviousFields struct {
	Balance  interface{} `json:"Balance"`
	Sequence int         `json:"Sequence"`
}

//---------------------------------------------------------------------------//
// cmd

type CmdType string

const (
	Cmd_Ping            CmdType = "ping"
	Cmd_ServerInfo      CmdType = "server_info"
	Cmd_RipplePathFind  CmdType = "ripple_path_find"
	Cmd_WalletPropose   CmdType = "wallet_propose"
	Cmd_Submit          CmdType = "submit"
	Cmd_Sign            CmdType = "sign"
	Cmd_Fee             CmdType = "fee"
	Cmd_AccountInfo     CmdType = "account_info"
	Cmd_AccountTx       CmdType = "account_tx"
	Cmd_AccountChannels CmdType = "account_channels"
	Cmd_Tx              CmdType = "tx"
	Cmd_Ledger          CmdType = "ledger"
	Cmd_LedgerClosed    CmdType = "ledger_closed"
	Cmd_LedgerData      CmdType = "ledger_data"
	Cmd_LedgerHeader    CmdType = "ledger_header"
	Cmd_Subscribe       CmdType = "subscribe"
)

type CmdStatus string

const (
	Cmd_success CmdStatus = "success"
	Cmd_fail    CmdStatus = "fail"
)
