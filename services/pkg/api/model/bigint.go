package model

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type BigInt *big.Int

func MarshalBigInt(bignum *big.Int) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		if bignum.Sign() == 0 {
			_, _ = w.Write([]byte(`"0x0"`))
		} else {
			_, _ = w.Write([]byte(strconv.Quote(hexutil.Encode(bignum.Bytes()))))
		}
	})
}

func UnmarshalBigInt(v interface{}) (*big.Int, error) {
	switch v := v.(type) {
	case string:
		b, err := hexutil.Decode(v)
		if err != nil {
			return nil, fmt.Errorf("%v failed to decode as BigInt", v)
		}
		return big.NewInt(0).SetBytes(b), nil
	case []byte:
		return big.NewInt(0).SetBytes(v), nil
	case int:
		return big.NewInt(int64(v)), nil
	case bool:
		n := 0
		if v {
			n = 1
		}
		return big.NewInt(int64(n)), nil
	default:
		return nil, fmt.Errorf("%T is not a decodable as BigInt", v)
	}
}

// genqlient marshallers

func ClientMarshalBigInt(bignum *BigInt) ([]byte, error) {
	n := (*big.Int)(*bignum)
	return []byte(strconv.Quote(hexutil.Encode(n.Bytes()))), nil
}

func ClientUnmarshalBigInt(b []byte, v *BigInt) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "0x0" {
		*v = BigInt(big.NewInt(0))
		return nil
	}

	nBytes, err := hexutil.Decode(s)
	if err != nil {
		return err
	}
	bn := BigInt(big.NewInt(0).SetBytes(nBytes))
	*v = bn
	return nil
}
