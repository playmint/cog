package sequencer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/contracts/router"
	"github.com/playmint/ds-node/pkg/indexer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

// -- end dummy action

type Sequencer interface {
	Ready() chan struct{}
	Enqueue(
		ctx context.Context,
		routerAddress common.Address,
		owner common.Address,
		stateAddress common.Address,
		actionPayload []string,
		actionSig string,
	) (*model.ActionTransaction, error)
	Signin(
		ctx context.Context,
		routerAddr common.Address,
		dispatcherAddr common.Address,
		sessionKey common.Address,
		ttl uint32,
		scopes uint32,
		permit string,
	) error
	Signout(ctx context.Context, routerAddr common.Address, sessionKey common.Address, permit string) error

	GetTransactions(routerAddr common.Address, owner *string, status []model.ActionTransactionStatus) ([]*model.ActionTransaction, error)
	GetTransaction(routerAddr common.Address, id string) (*model.ActionTransaction, error)
}

var _ Sequencer = &MemorySequencer{}

type MemorySequencer struct {
	PrivateKey        *ecdsa.PrivateKey
	chainProviderHTTP string
	chainHttpClient   *alchemy.Client
	notifications     chan interface{}
	idxr              indexer.Indexer
	log               zerolog.Logger
}

func NewMemorySequencer(
	ctx context.Context,
	key *ecdsa.PrivateKey,
	notifications chan interface{},
	chainProviderHTTP string,
	idxr indexer.Indexer,
) (*MemorySequencer, error) {

	var err error
	seqr := &MemorySequencer{
		PrivateKey:        key,
		notifications:     notifications,
		log:               log.With().Str("service", "sequencer").Logger(),
		idxr:              idxr,
		chainProviderHTTP: chainProviderHTTP,
	}
	// setup an RPC client
	seqr.chainHttpClient, err = alchemy.Dial(
		seqr.chainProviderHTTP,
		1,
		seqr.PrivateKey,
	)
	if err != nil {
		return nil, err
	}

	return seqr, nil
}

func (seqr *MemorySequencer) Ready() chan struct{} {
	ch := make(chan struct{})
	go func() {
		close(ch)
	}()
	return ch
}

func (seqr *MemorySequencer) emitTx(tx *model.ActionTransaction) {
	seqr.notifications <- tx
}

// Enqueue dispatches and waits for action commit
func (seqr *MemorySequencer) Enqueue(ctx context.Context,
	routerAddr common.Address,
	ownerAddr common.Address,
	stateAddr common.Address,
	actionData []string,
	actionSig string,
) (*model.ActionTransaction, error) {
	if routerAddr.Hex() == "" || len(actionData) == 0 || actionSig == "" {
		return nil, fmt.Errorf("invalid action data")
	}

	actionTx := &model.ActionTransaction{
		ID:            uuid.NewV4().String(),
		Owner:         ownerAddr.Hex(),
		Payload:       actionData,
		Sig:           actionSig,
		RouterAddress: routerAddr.Hex(),
		Batch:         &model.ActionBatch{},
	}

	tx, err := seqr.dispatch(ctx, routerAddr, actionTx)
	if err != nil {
		seqr.log.Error().
			Err(err).
			Msg("action-rejected")
		return actionTx, err
	}

	_, err = WaitMined(ctx, seqr.chainHttpClient, tx)
	if err != nil {
		seqr.log.Error().
			Err(err).
			Str("hash", tx.Hash().Hex()).
			Msg("action-fail")
		return actionTx, err
	}

	seqr.log.Info().
		Str("hash", tx.Hash().Hex()).
		Msg("action-commited")

	return actionTx, nil
}

func (seqr *MemorySequencer) dispatch(
	ctx context.Context,
	routerAddr common.Address,
	action *model.ActionTransaction,
) (*types.Transaction, error) {
	// lock client
	client := seqr.chainHttpClient
	client.Lock()
	defer client.Unlock()
	// prep action data
	actions := [][][]byte{}
	sigs := [][]byte{}
	actions = append(actions, action.ActionBytes())
	sigs = append(sigs, action.ActionSig())

	sessionRouter, err := router.NewSessionRouter(routerAddr, client)
	if err != nil {
		return nil, err
	}

	txOpts, err := client.NewRelayTransactor(ctx)
	if err != nil {
		return nil, err
	}

	txOpts.Context = ctx
	txOpts.Value = big.NewInt(0)
	txOpts.GasLimit = uint64(15000000)
	txOpts.GasPrice = big.NewInt(1000000500)

	tx, err := sessionRouter.Dispatch(txOpts,
		actions,
		sigs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed commit batch tx: %v", err)
	}
	defer client.IncrementRelayNonce(ctx)
	return tx, nil
}

func WaitMined(ctx context.Context, client *alchemy.Client, tx *types.Transaction) (*types.Receipt, error) {
	// wait til batch success
	maxWait, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	time.Sleep(50 * time.Millisecond)
	rcpt, err := bind.WaitMined(maxWait, client, tx)
	if err != nil {
		return nil, err
	}
	switch rcpt.Status {
	case 1:
	default:
		reason := errorReason(ctx, client, client.Address(), tx, rcpt.BlockNumber)
		return rcpt, fmt.Errorf("%s", reason)
	}
	return rcpt, nil
}

func (seqr *MemorySequencer) Signout(ctx context.Context, routerAddr common.Address, sessionKey common.Address, permit string) error {
	client := seqr.chainHttpClient
	client.Lock()
	defer client.Unlock()

	// lookup the account contract
	sessionRouter, err := router.NewSessionRouter(routerAddr, client)
	if err != nil {
		return err
	}

	// setup tx
	txOpts, err := client.NewRelayTransactor(ctx)
	if err != nil {
		return err
	}
	txOpts.Value = big.NewInt(0)      // in wei
	txOpts.GasLimit = uint64(3000000) // in units

	// decode the permit into sig parts
	sig, err := hexutil.Decode(permit)
	if err != nil {
		return err
	}

	_, err = sessionRouter.RevokeAddr(txOpts, sessionKey, sig)
	if err != nil {
		return fmt.Errorf("failed perform signin tx for session=%v: %v", sessionKey, err)
	}
	defer client.IncrementRelayNonce(ctx)

	return nil
}

func (seqr *MemorySequencer) Signin(ctx context.Context, routerAddr common.Address, dispatcherAddr common.Address, sessionKey common.Address, ttl uint32, scopes uint32, permit string) error {
	client := seqr.chainHttpClient
	client.Lock()
	defer client.Unlock()

	// lookup the account contract
	sessionRouter, err := router.NewSessionRouter(routerAddr, client)
	if err != nil {
		return err
	}

	// setup tx
	txOpts, err := client.NewRelayTransactor(ctx)
	if err != nil {
		return err
	}
	txOpts.Value = big.NewInt(0)      // in wei
	txOpts.GasLimit = uint64(3000000) // in units

	// decode the permit into sig parts
	sig, err := hexutil.Decode(permit)
	if err != nil {
		return err
	}

	tx, err := sessionRouter.AuthorizeAddr0(txOpts, dispatcherAddr, ttl, scopes, sessionKey, sig)
	if err != nil {
		return fmt.Errorf("failed perform signin tx for session=%v: %v", sessionKey, err)
	}
	seqr.log.Info().
		Str("session", sessionKey.Hex()).
		Uint32("ttl", ttl).
		Uint32("scopes", scopes).
		Str("dispatcher", dispatcherAddr.Hex()).
		Str("router", routerAddr.Hex()).
		Msg("signin")
	defer client.IncrementRelayNonce(ctx)

	// wait mined
	maxWait, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()
	rcpt, err := bind.WaitMined(maxWait, client, tx)
	if err != nil {
		return err
	}
	switch rcpt.Status {
	case 1:
	default:
		return fmt.Errorf("signin failed: %v", rcpt)
	}

	return nil
}

func (seqr *MemorySequencer) GetBatch(routerAddr common.Address, id string) (*model.ActionBatch, error) {
	return nil, nil
}

func (seqr *MemorySequencer) GetTransactions(routerAddr common.Address, owner *string, status []model.ActionTransactionStatus) ([]*model.ActionTransaction, error) {
	txs := []*model.ActionTransaction{}
	return txs, nil
}

func (seqr *MemorySequencer) GetTransaction(routerAddr common.Address, id string) (*model.ActionTransaction, error) {
	txs, err := seqr.GetTransactions(routerAddr, nil, nil)
	if err != nil {
		return nil, err
	}
	for _, tx := range txs {
		if tx.ID == id {
			return tx, nil
		}
	}
	return nil, nil
}

func inStatusList(haystack []model.ActionTransactionStatus, needle model.ActionTransactionStatus) bool {
	for _, hay := range haystack {
		if needle == hay {
			return true
		}
	}
	return false
}

// - Sig and Payload are non empty
// - A batch is set
// - A router destination is set
// - the signature correctly verifies
func ValidateActions(payloads []string, sig string) (*common.Address, error) {
	byteArray, _ := abi.NewType("bytes[]", "bytes[]", nil)
	args := abi.Arguments{
		abi.Argument{Type: byteArray},
	}
	actions := [][]byte{}
	for _, action := range payloads {
		b, _ := hexutil.Decode(action)
		actions = append(actions, b)
	}
	payload, err := args.Pack(actions)
	if err != nil {
		return nil, fmt.Errorf("invalid action: unable to pack actions byte array: %v", err)
	}
	if len(payload) < 4 {
		return nil, fmt.Errorf("invalid action: expected payload len>=4 got len=%d", len(payload))
	}
	permit, _ := hexutil.Decode(sig)
	if len(permit) == 0 {
		return nil, fmt.Errorf("invalid action: expected sig len>0 got len=%d", len(permit))
	}
	if permit[len(permit)-1] > 26 {
		permit[len(permit)-1] -= 27
	}
	digest := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n32"),
		crypto.Keccak256Hash(payload).Bytes(),
	)
	pubKey, err := crypto.SigToPub(digest.Bytes(), permit)
	if err != nil || pubKey == nil {
		return nil, fmt.Errorf("invalid action: unable to recover signer: %v", err)
	}
	signer := crypto.PubkeyToAddress(*pubKey)
	if signer.Hex() == common.BigToAddress(big.NewInt(0)).Hex() {
		return nil, fmt.Errorf("invalid action: bad signature")
	}
	return &signer, nil
}

func errorReason(ctx context.Context, b ethereum.ContractCaller, from common.Address, tx *types.Transaction, blockNum *big.Int) (reason string) {
	msg := ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	_, err := b.CallContract(ctx, msg, blockNum)
	if err != nil {
		reason = err.Error()
		reason = strings.TrimPrefix(reason, "CallContract: ")
		reason = strings.TrimPrefix(reason, "execution reverted: ")
		return reason
	}
	return "failed"
}
