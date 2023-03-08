package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	BlockHash    *common.Hash   `json:"blockHash"`
	BlockNumber  *uint64        `json:"blockNumber"`
	To           common.Address `json:"to" gencodec:"required"`
	From         common.Address `json:"from"`
	Gas          string         `json:"gas"`
	GasPrice     string         `json:"gasPrice"`
	Input        string         `json:"input"`
	Nonce        string         `json:"nonce"`
	TxHash       common.Hash    `json:"hash" gencodec:"required"`
	TxIndex      uint64         `json:"transactionIndex"`
	Value        string         `json:"value"`
	V            string         `json:"v"`
	R            string         `json:"r"`
	S            string         `json:"s"`
	Subscription string         `json:"subscription"`
}
