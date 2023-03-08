package contracts

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
)

// IsMoreRecent returns true if event a is "newer" than
// event b
func IsMoreRecent(a, b *types.Log) bool {
	if a.BlockNumber > b.BlockNumber {
		return true
	}
	if a.BlockNumber < b.BlockNumber {
		return false
	}
	if a.TxIndex > b.TxIndex {
		return true
	}
	if a.TxIndex < b.TxIndex {
		return false
	}
	if a.Index == b.Index {
		// this should not happen if it does, then data
		// will not be consistent so just panic
		panic(fmt.Errorf("unorderable event: %v", a))
	}
	return a.Index > b.Index
}
