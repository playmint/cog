package config

import (
	"crypto/ecdsa"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/crypto"
)

func getRequiredEnvString(name string) string {
	v := os.Getenv(name)
	if v == "" {
		return ""
	}
	return v
}

func getOptionalEnvInt(name string, defvalue int) int {
	vs := os.Getenv(name)
	if vs == "" {
		return defvalue
	}
	v, err := strconv.Atoi(vs)
	if err != nil {
		return 0
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
		return nil
	}
	publicKey := privateKey.Public()
	_, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil
	}
	// relayAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey
}
