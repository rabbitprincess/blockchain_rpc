package types

//--------------------------------------------------------------------------//
// ledger

type Ledger struct {
	Accepted            bool             `json:"accepted"`
	AccountHash         string           `json:"account_hash"`
	CloseFlags          int              `json:"close_flags"`
	CloseTime           int              `json:"close_time"`
	CloseTimeHuman      string           `json:"close_time_human"`
	CloseTimeResolution int              `json:"close_time_resolution"`
	Closed              bool             `json:"closed"`
	Hash                string           `json:"hash"`
	LedgerHash          string           `json:"ledger_hash"`
	LedgerIndex         string           `json:"ledger_index"`
	ParentCloseTime     int              `json:"parent_close_time"`
	ParentHash          string           `json:"parent_hash"`
	SeqNum              string           `json:"seqNum"`
	TotalCoins          string           `json:"total_coins"`
	TransactionHash     string           `json:"transaction_hash"`
	Transactions        []TransactionRes `json:"transactions"`
}

//--------------------------------------------------------------------------//
// tx

type TransactionRes struct {
	Account            string    `json:"Account"`
	Amount             Amount    `json:"Amount"`
	Destination        string    `json:"Destination"`
	DestinationTag     uint32    `json:"DestinationTag"`
	Fee                string    `json:"Fee"`
	Flags              TxFlag    `json:"Flags"`
	LastLedgerSequence int64     `json:"LastLedgerSequence"`
	OfferSequence      int64     `json:"OfferSequence"`
	Sequence           int64     `json:"Sequence"`
	SigningPubKey      string    `json:"SigningPubKey"`
	TakerGets          Amount    `json:"TakerGets"`
	TakerPays          Amount    `json:"TakerPays"`
	TransactionType    TxType    `json:"TransactionType"`
	TxnSignature       string    `json:"TxnSignature"`
	Date               int64     `json:"date"`
	Hash               string    `json:"hash"`
	InLedger           int64     `json:"inLedger"`
	LedgerIndex        int64     `json:"ledger_index"`
	Metadata           Metadata  `json:"meta"`
	Status             CmdStatus `json:"status"`
	Validated          bool      `json:"validated"`
	RPCError
}

type TransactionPayment struct {
	TransactionType TxType `json:"transaction_type"`
	Account         string `json:"account"`
	Fee             string `json:"fee"`
	Sequence        uint32 `json:"sequence"`
	Destination     string `json:"destination"`
	Amount          string `json:"amount"`
	DestinationTag  uint32 `json:"destination_tag,omitempty"`
	InvoiceID       []byte `json:"invoice_id,omitempty"`
}

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
