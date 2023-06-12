package alchemy

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/playmint/ds-node/pkg/client/types"
)

// Client defines typed wrappers for the Ethereum RPC API.
type Client struct {
	privateKey *ecdsa.PrivateKey
	publicKey  ecdsa.PublicKey
	nonces     map[common.Address]uint64
	rpc        *rpc.Client
	*ethclient.Client
	concurrencyLock chan struct{}
	sync.Mutex
}

// Dial connects a client to the given URL.
func Dial(rawurl string, concurrency int, key *ecdsa.PrivateKey) (*Client, error) {
	return DialContext(context.Background(), rawurl, concurrency, key)
}

func DialContext(ctx context.Context, rawurl string, concurrency int, key *ecdsa.PrivateKey) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c, concurrency, key)
}

// NewClient creates a client that uses the given RPC client.
func NewClient(rpc *rpc.Client, concurrency int, key *ecdsa.PrivateKey) (*Client, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("concurrency must be at least 1")
	}
	c := &Client{
		privateKey:      key,
		nonces:          map[common.Address]uint64{},
		rpc:             rpc,
		concurrencyLock: make(chan struct{}, concurrency),
		Client:          ethclient.NewClient(rpc),
	}
	if key != nil {
		c.publicKey = *key.Public().(*ecdsa.PublicKey)
	}
	return c, nil
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (c *Client) SubscribeFilterPending(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Transaction) (ethereum.Subscription, error) {
	args := map[string]interface{}{}
	if len(q.Addresses) > 0 {
		args["address"] = q.Addresses[0]
	}
	return c.rpc.EthSubscribe(ctx, ch, "alchemy_filteredNewFullPendingTransactions", args)
}

func (c *Client) PendingNonceAt(ctx context.Context, addr common.Address) (uint64, error) {
	_, ok := c.nonces[addr]
	if !ok {
		latestNonce, err := c.Client.PendingNonceAt(ctx, addr)
		if err != nil {
			return 0, err
		}
		c.nonces[addr] = latestNonce
	}
	return c.nonces[addr], nil
}

func (c *Client) IncrementRelayNonce(ctx context.Context) {
	c.nonces[c.Address()]++
}

func (c *Client) ConcurrencyLock() {
	c.concurrencyLock <- struct{}{}
}

func (c *Client) ConcurrencyUnlock() {
	<-c.concurrencyLock
}

func (c *Client) NewRelayTransactor(ctx context.Context) (*bind.TransactOpts, error) {

	nonce, err := c.PendingNonceAt(ctx, c.Address())
	if err != nil {
		return nil, err
	}

	// gasPrice, err := c.SuggestGasPrice(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	chainID, err := c.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(
		c.privateKey,
		chainID,
	)
	if err != nil {
		return nil, err
	}
	txOpts.Nonce = big.NewInt(int64(nonce))
	// txOpts.GasPrice = gasPrice

	return txOpts, nil
}

// [{"forking": {"jsonRpcUrl": "httKBrdf", "blockNumber": 14000000}}
func (c *Client) Reset(ctx context.Context, remoteURL string, remoteBlockNumber uint64) (uint64, error) {
	forking := map[string]interface{}{
		"jsonRpcUrl":  remoteURL,
		"blockNumber": remoteBlockNumber,
	}
	params := map[string]interface{}{
		"forking": forking,
	}
	var hex hexutil.Uint64
	err := c.rpc.CallContext(ctx, &hex, "anvil_reset", params)
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

func (c *Client) EstimateContractGas(ctx context.Context, opts *bind.TransactOpts, contract *common.Address, input []byte) error {

	// head header
	head, err := c.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}
	// Normalize value
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	// Estimate GasPrice
	gasPrice, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	// Estimate FeeCap
	gasTipCap, err := c.SuggestGasTipCap(opts.Context)
	if err != nil {
		return nil
	}
	// Estimate FeeCap
	gasFeeCap := opts.GasFeeCap
	if gasFeeCap == nil {
		gasFeeCap = new(big.Int).Add(
			gasTipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
	}
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
	}
	msg := ethereum.CallMsg{
		From:      opts.From,
		To:        contract,
		GasPrice:  gasPrice,
		GasTipCap: gasTipCap,
		GasFeeCap: gasFeeCap,
		Value:     value,
		Data:      input,
	}
	_, err = c.EstimateGas(opts.Context, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Address() common.Address {
	return crypto.PubkeyToAddress(c.publicKey)
}
