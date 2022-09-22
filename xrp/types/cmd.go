package types

import "encoding/json"

type CmdReq interface {
	Marshal() ([]byte, error)
	Type() CmdType
}

type CmdRes interface {
	Unmarshal([]byte) error
	Error() *RPCError
}

//---------------------------------------------------------------------------//
// cmd

type Req_ping struct {
	Command CmdType `json:"command"`
	ID      int     `json:"id"`
}

// ping
type Res_ping struct {
	ID     int             `json:"id"`
	Result Res_ping_result `json:"result"`
	Status string          `json:"status"`
	Type   string          `json:"type"`
}

type Res_ping_result struct {
	RPCError
}

func (t Req_ping) Type() CmdType {
	return t.Command
}

func (t *Req_ping) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_ping) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_ping) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

// server info
type Req_serverInfo struct {
	Method CmdType `json:"method"`
}

type Res_serverInfo struct {
	Result Res_serverInfo_result `json:"result"`
}

type Res_serverInfo_result struct {
	Info   Res_serverInfo_result_info `json:"info"`
	Status CmdStatus                  `json:"status"`
	RPCError
}

type Res_serverInfo_result_info struct {
	BuildVersion             string          `json:"build_version"`
	CompleteLedgers          string          `json:"complete_ledgers"`
	Hostid                   string          `json:"hostid"`
	IoLatencyMs              int             `json:"io_latency_ms"`
	JqTransOverflow          string          `json:"jq_trans_overflow"`
	LastClose                LastClose       `json:"last_close"`
	LoadFactor               float64         `json:"load_factor"`
	NodeSize                 string          `json:"node_size"`
	PeerDisconnects          string          `json:"peer_disconnects"`
	PeerDisconnectsResources string          `json:"peer_disconnects_resources"`
	Peers                    int             `json:"peers"`
	PubkeyNode               string          `json:"pubkey_node"`
	PubkeyValidator          string          `json:"pubkey_validator"`
	ServerState              string          `json:"server_state"`
	ServerStateDurationUs    string          `json:"server_state_duration_us"`
	StateAccounting          StateAccounting `json:"state_accounting"`
	Time                     string          `json:"time"`
	Uptime                   int             `json:"uptime"`
	ValidatedLedger          ValidatedLedger `json:"validated_ledger"`
	ValidationQuorum         int             `json:"validation_quorum"`
}

type LastClose struct {
	ConvergeTimeS float64 `json:"converge_time_s"`
	Proposers     int     `json:"proposers"`
}

type StateAccounting struct {
	Connected    ConnectState `json:"connected"`
	Disconnected ConnectState `json:"disconnected"`
	Full         ConnectState `json:"full"`
	Syncing      ConnectState `json:"syncing"`
	Tracking     ConnectState `json:"tracking"`
}

type ConnectState struct {
	DurationUs  string      `json:"duration_us"`
	Transitions interface{} `json:"transitions"`
}

type ValidatedLedger struct {
	Age            int     `json:"age"`
	BaseFeeXrp     float64 `json:"base_fee_xrp"`
	Hash           string  `json:"hash"`
	ReserveBaseXrp int     `json:"reserve_base_xrp"`
	ReserveIncXrp  int     `json:"reserve_inc_xrp"`
	Seq            int     `json:"seq"`
}

func (t Req_serverInfo) Type() CmdType {
	return t.Method
}

func (t *Req_serverInfo) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_serverInfo) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_serverInfo) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// fee

type Req_fee struct {
	Method CmdType `json:"method"`
}

type Res_fee struct {
	Result Res_fee_Result `json:"result"`
}

type Res_fee_Result struct {
	CurrentLedgerSize  string `json:"current_ledger_size"`
	CurrentQueueSize   string `json:"current_queue_size"`
	Drops              Drops  `json:"drops"`
	ExpectedLedgerSize string `json:"expected_ledger_size"`
	LedgerCurrentIndex int    `json:"ledger_current_index"`
	Levels             Levels `json:"levels"`
	MaxQueueSize       string `json:"max_queue_size"`
	Status             string `json:"status"`

	RPCError
}

func (t Req_fee) Type() CmdType {
	return t.Method
}

func (t *Req_fee) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_fee) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_fee) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// tx ( tx info )

type Req_tx struct {
	Method CmdType         `json:"method"`
	Params []Req_tx_params `json:"params"`
}

type Req_tx_params struct {
	Transaction string `json:"transaction"`
	Binary      bool   `json:"binary"`
	MinLedger   int64  `json:"min_ledger"`
	MaxLedger   int64  `json:"max_ledger"`
}

type Res_tx struct {
	Result TransactionRes `json:"result"`
}

func (t Req_tx) Type() CmdType {
	return t.Method
}

func (t *Req_tx) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_tx) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_tx) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// account info

type Req_accountInfo struct {
	Method CmdType                  `json:"method"`
	Params []Req_accountInfo_params `json:"params"`
}

type Req_accountInfo_params struct {
	Account     string `json:"account"`
	Strict      bool   `json:"strict"`
	LedgerIndex string `json:"ledger_index"`
	Queue       bool   `json:"queue"`
}

type Res_accountInfo struct {
	Result Res_accountInfo_result `json:"result"`
}

type Res_accountInfo_result struct {
	AccountData Res_accountInfo_result_data `json:"account_data"`
	LedgerIndex int                         `json:"ledger_index"`
	Status      CmdStatus                   `json:"status"`
	RPCError
	Request Req_accountInfo_params `json:"request,omitempty"`
}

type Res_accountInfo_result_data struct {
	Account           string `json:"Account"`
	Balance           string `json:"Balance"`
	Flags             TxFlag `json:"Flags"`
	LedgerEntryType   string `json:"LedgerEntryType"`
	OwnerCount        int    `json:"OwnerCount"`
	PreviousTxnID     string `json:"PreviousTxnID"`
	PreviousTxnLgrSeq int    `json:"PreviousTxnLgrSeq"`
	Sequence          int    `json:"Sequence"`
	Index             string `json:"index"`
}

func (t Req_accountInfo) Type() CmdType {
	return t.Method
}

func (t *Req_accountInfo) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_accountInfo) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_accountInfo) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// account tx

type Req_accountTx struct {
	Method CmdType                `json:"method"`
	Params []Req_accountTx_params `json:"params"`
}

type Req_accountTx_params struct {
	Account        string `json:"account"`
	Binary         bool   `json:"binary"`
	Forward        bool   `json:"forward"`
	LedgerHash     string `json:"ledger_hash"`
	LedgerIndex    int    `json:"ledger_index"`
	LedgerIndexMin int    `json:"ledger_index_min"`
	LedgerIndexMax int    `json:"ledger_index_max"`
	Limit          int    `json:"limit"`
}

type Res_accountTx struct {
	Result Res_accountTx_result `json:"result"`
}

type Res_accountTx_result struct {
	Account        string                      `json:"account"`
	Marker         Res_accountTx_result_marker `json:"marker"`
	Status         CmdStatus                   `json:"status"`
	Transactions   []Transactions              `json:"transactions"`
	Validated      bool                        `json:"validated"`
	Limit          int                         `json:"limit"`
	LedgerIndexMax int                         `json:"ledger_index_max"`
	LedgerIndexMin int                         `json:"ledger_index_min"`
	RPCError
}

type Res_accountTx_result_marker struct {
	Ledger int `json:"ledger"`
	Seq    int `json:"seq"`
}

type Transactions struct {
	Meta      Metadata       `json:"meta"`
	Tx        TransactionRes `json:"tx"`
	Validated bool           `json:"validated"`
}

func (t Req_accountTx) Type() CmdType {
	return t.Method
}

func (t *Req_accountTx) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_accountTx) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_accountTx) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// ledger ( get ledger )

type Req_ledger struct {
	Method CmdType             `json:"method"`
	Params []Req_ledger_params `json:"params"`
}

type Req_ledger_params struct {
	LedgerHash   string `json:"ledger_hash"`
	LedgerIndex  string `json:"ledger_index"`
	Accounts     bool   `json:"accounts"`
	Full         bool   `json:"full"`
	Transactions bool   `json:"transactions"`
	Expand       bool   `json:"expand"`
	OwnerFunds   bool   `json:"owner_funds"`
}

type Res_ledger struct {
	Result Res_ledger_Result `json:"result"`
}

type Res_ledger_Result struct {
	Ledger      LedgerRes `json:"ledger"`
	LedgerHash  string    `json:"ledger_hash"`
	LedgerIndex int       `json:"ledger_index"`
	Status      CmdStatus `json:"status"`
	Validated   bool      `json:"validated"`
	RPCError
}

func (t Req_ledger) Type() CmdType {
	return t.Method
}

func (t *Req_ledger) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_ledger) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_ledger) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// wallet propose

type Req_walletPropose_by_passphrase struct {
	Method     CmdType `json:"method"`
	Passphrase string  `json:"passphrase"`
}

type Req_walletPropose struct {
	Method CmdType                    `json:"method"`
	Params []Req_walletPropose_params `json:"params"`
}

type Req_walletPropose_params struct {
	Seed    string        `json:"seed"`
	KeyType TD_s_key_type `json:"key_type"`
}

type Res_walletPropose struct {
	Result Res_walletPropose_Result `json:"result"`
}

type Res_walletPropose_Result struct {
	AccountID     string    `json:"account_id"`
	KeyType       string    `json:"key_type"`
	MasterKey     string    `json:"master_key"`
	MasterSeed    string    `json:"master_seed"`
	MasterSeedHex string    `json:"master_seed_hex"`
	PublicKey     string    `json:"public_key"`
	PublicKeyHex  string    `json:"public_key_hex"`
	Status        CmdStatus `json:"status"`
	RPCError
}

func (t Req_walletPropose_by_passphrase) Type() CmdType {
	return t.Method
}

func (t *Req_walletPropose_by_passphrase) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Req_walletPropose) Type() CmdType {
	return t.Method
}

func (t *Req_walletPropose) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_walletPropose) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_walletPropose) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// sign ( sign tx )

type Req_sign struct {
	Method CmdType           `json:"method"`
	Params []Req_sign_params `json:"params"`
}

type Req_sign_params struct {
	Secret string `json:"secret"`
	TxJson []byte `json:"tx_json"`
}

type Res_sign struct {
	Result Res_sign_result `json:"result"`
}

type Res_sign_result struct {
	Status CmdStatus               `json:"status"`
	TxBlob string                  `json:"tx_blob"`
	TxJSON *Res_sign_result_TxJSON `json:"tx_json"`
	RPCError
}

type Res_sign_result_TxJSON struct {
	Account         string `json:"Account"`
	Amount          Amount `json:"Amount"`
	Destination     string `json:"Destination"`
	Fee             string `json:"Fee"`
	Flags           TxFlag `json:"Flags"`
	Sequence        int    `json:"Sequence"`
	SigningPubKey   string `json:"SigningPubKey"`
	TransactionType string `json:"TransactionType"`
	TxnSignature    string `json:"TxnSignature"`
	Hash            string `json:"hash"`
}

func (t Req_sign) Type() CmdType {
	return t.Method
}

func (t *Req_sign) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_sign) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_sign) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// sign_by_passphrase ( sign tx ) - res 는 sign 과 공유

type Req_sign_by_passphrase struct {
	Method CmdType                         `json:"method"`
	Params []Req_sign_by_passphrase_params `json:"params"`
}

type Req_sign_by_passphrase_params struct {
	KeyType    string              `json:"key_type"`
	Passphrase string              `json:"passphrase"`
	TxJson     *TransactionPayment `json:"tx_json"`
}

type Res_sign_by_passphrase struct {
	Result Res_sign_result `json:"result"`
}

func (t Req_sign_by_passphrase) Type() CmdType {
	return t.Method
}

func (t *Req_sign_by_passphrase) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

//------------------------------------------------------------------------//
// submit ( submit tx )

type Req_submit struct {
	Method CmdType             `json:"method"`
	Params []Req_submit_params `json:"params"`
}

type Req_submit_params struct {
	TxBlob string `json:"tx_blob"`
}

type Res_submit struct {
	Result               Res_submit_result `json:"result"`
	ValidatedLedgerIndex int               `json:"validated_ledger_index"`
}

type Res_submit_result struct {
	Accepted                 bool                      `json:"accepted"`
	AccountSequenceAvailable int                       `json:"account_sequence_available"`
	AccountSequenceNext      int                       `json:"account_sequence_next"`
	Applied                  bool                      `json:"applied"`
	Broadcast                bool                      `json:"broadcast"`
	EngineResult             string                    `json:"engine_result"`
	EngineResultCode         TxResult                  `json:"engine_result_code"`
	EngineResultMessage      string                    `json:"engine_result_message"`
	Status                   CmdStatus                 `json:"status"`
	Kept                     bool                      `json:"kept"`
	OpenLedgerCost           string                    `json:"open_ledger_cost"`
	Queued                   bool                      `json:"queued"`
	TxBlob                   string                    `json:"tx_blob"`
	TxJSON                   Res_submit_result_tx_json `json:"tx_json"`
	RPCError
}

type Res_submit_result_tx_json struct {
	Account         string `json:"Account"`
	Amount          Amount `json:"Amount"`
	Destination     string `json:"Destination"`
	Fee             string `json:"Fee"`
	Flags           TxFlag `json:"Flags"`
	Sequence        int    `json:"Sequence"`
	SigningPubKey   string `json:"SigningPubKey"`
	TransactionType string `json:"TransactionType"`
	TxnSignature    string `json:"TxnSignature"`
	Hash            string `json:"hash"`
}

func (t Req_submit) Type() CmdType {
	return t.Method
}

func (t *Req_submit) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_submit) Error() *RPCError {
	return &t.Result.RPCError
}

func (t *Res_submit) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// walletPropose ( nodejs )

type Req_offline_walletPropose struct {
	Method       CmdType `json:"method"`
	PrivateKey   string  `json:"private_key"`
	N_wallet_cnt int     `json:"wallet_cnt"`
}

type Res_offline_walletPropose struct {
	Arrpt_wallet []*Wallet `json:"wallet"`
	RPCError
}

func (t Req_offline_walletPropose) Type() CmdType {
	return t.Method
}

func (t *Req_offline_walletPropose) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_offline_walletPropose) Error() *RPCError {
	return &t.RPCError
}

func (t *Res_offline_walletPropose) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}

//------------------------------------------------------------------------//
// sign ( nodejs )

type Req_offline_sign struct {
	Method     CmdType `json:"method"`
	PrivateKey string  `json:"private_key"`
	TxJSON     []byte  `json:"tx"`
}

type Res_offline_sign struct {
	TxBlob string `json:"tx_blob"`
	Hash   string `json:"hash"`
	RPCError
}

func (t Req_offline_sign) Type() CmdType {
	return t.Method
}

func (t *Req_offline_sign) Marshal() ([]byte, error) {
	return json.Marshal(&t)
}

func (t Res_offline_sign) Error() *RPCError {
	return &t.RPCError
}

func (t *Res_offline_sign) Unmarshal(_bt []byte) error {
	return json.Unmarshal(_bt, &t)
}
