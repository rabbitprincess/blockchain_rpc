package xrp

import "encoding/json"

//---------------------------------------------------------------------------//
// type

type CmdType string

type RPCError struct {
	ErrName string  `json:"error,omitempty"`
	ErrCode ErrCode `json:"error_code"`
	ErrMsg  string  `json:"error_message"`
}

func (t RPCError) Error() string {
	bt_json, _ := json.MarshalIndent(&t, "", "\t")
	return string(bt_json)
}

type ErrCode int64

const (
	Error ErrCode = iota
	Ok
)

//---------------------------------------------------------------------------//
// cmd

type CmdReq interface {
	Marshal() ([]byte, error)
	Type() CmdType
}

type CmdRes interface {
	Unmarshal([]byte) error
	Error() *RPCError
}
