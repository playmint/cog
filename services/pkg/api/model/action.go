package model

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ActionTransaction struct {
	ID            string   `json:"id"`
	Payload       []string `json:"payload"`
	Sig           string   `json:"sig"`
	RouterAddress string
	Owner         string `json:"owner"`
	Batch         *ActionBatch
}

func (a *ActionTransaction) ActionBytes() [][]byte {
	bs := [][]byte{}
	for _, payload := range a.Payload {
		b, _ := hexutil.Decode(payload)
		bs = append(bs, b)
	}
	return bs
}

func (a *ActionTransaction) ActionSig() []byte {
	permit, _ := hexutil.Decode(a.Sig)
	return permit
}

func (a *ActionTransaction) Router() *Router {
	return &Router{
		ID: a.RouterAddress,
	}
}

func (a *ActionTransaction) Status() ActionTransactionStatus {
	if a.Batch == nil {
		return ActionTransactionStatusUnknown
	}
	return a.Batch.Status
}

type ActionBatch struct {
	ID            string                  `json:"id"`
	Tx            *string                 `json:"tx"`
	Status        ActionTransactionStatus `json:"status"`
	Transactions  []*ActionTransaction    `json:"transactions"`
	RouterAddress string
	Block         *int `json:"block"`
}
