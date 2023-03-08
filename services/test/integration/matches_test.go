package integration_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func EqualBig(expected interface{}) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(actual *big.Int) (bool, error) {
		var expectedHex string
		switch v := expected.(type) {
		case *big.Int:
			expectedHex = hexutil.EncodeBig(v)
		case int:
			expectedHex = hexutil.EncodeBig(big.NewInt(int64(v)))
		case int64:
			expectedHex = hexutil.EncodeBig(big.NewInt(int64(v)))
		case uint64:
			expectedHex = hexutil.EncodeBig(big.NewInt(int64(v)))
		}
		actualHex := hexutil.EncodeBig(actual)

		return expectedHex == actualHex, nil
	}).WithTemplate("Expected:\n{{.FormattedActual}}\n{{.To}} equal \n{{format .Data 1}}").WithTemplateData(expected)
}
