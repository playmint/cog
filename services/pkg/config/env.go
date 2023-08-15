package config

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func getRequiredEnvString(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(fmt.Sprintf("required environment variable %s missing", name))
	}
	return v
}

func getOptionalEnvAddress(name string, defvalue common.Address) common.Address {
	v := os.Getenv(name)
	if v == "" {
		return defvalue
	}
	return common.HexToAddress(v)
}

func getOptionalEnvInt(name string, defvalue int) int {
	vs := os.Getenv(name)
	if vs == "" {
		return defvalue
	}
	v, err := strconv.Atoi(vs)
	if err != nil {
		panic(fmt.Sprintf("environment variable %s contains invlaid value %s", name, vs))
	}
	return v
}

func getOptionalEnvBool(name string, defvalue string) bool {
	vs := os.Getenv(name)
	if vs == "" {
		vs = defvalue
	}
	v := vs != "false"
	return v
}

func getRequiredEnvKey(name string) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(os.Getenv(name))
	if err != nil {
		panic(fmt.Errorf("unable to decode private key: %v", err))
	}
	publicKey := privateKey.Public()
	_, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic(fmt.Errorf("unable to extract public key: %v", err))
	}
	// relayAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey
}
