package sequencer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client/alchemy"
	"github.com/playmint/ds-node/pkg/config"
	"github.com/playmint/ds-node/pkg/contracts/router"
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
	PrivateKey *ecdsa.PrivateKey
	// pending/processing are actions, grouped by Router that are waiting to join a batch
	pending    map[string]*model.ActionBatch
	processing map[string]*model.ActionBatch
	// batches are lists of actions, grouped by Router has processed
	success       map[string][]*model.ActionBatch
	failure       map[string][]*model.ActionBatch
	httpClient    *alchemy.Client
	notifications chan interface{}
	log           zerolog.Logger
	sync.RWMutex
}

func NewMemorySequencer(ctx context.Context, key *ecdsa.PrivateKey, notifications chan interface{}) (*MemorySequencer, error) {
	var err error
	seqr := &MemorySequencer{
		PrivateKey:    key,
		notifications: notifications,
		log:           log.With().Str("service", "sequencer").Logger(),
	}
	// setup an RPC client
	seqr.httpClient, err = alchemy.Dial(
		config.SequencerProviderHTTP,
		config.SequencerMaxConcurrency,
		config.SequencerPrivateKey,
	)
	if err != nil {
		return nil, err
	}

	// setup the queues and batches
	seqr.pending = map[string]*model.ActionBatch{}
	seqr.processing = map[string]*model.ActionBatch{}
	seqr.success = map[string][]*model.ActionBatch{}
	seqr.failure = map[string][]*model.ActionBatch{}

	// drain the queue every few seconds
	timer := time.NewTimer(time.Duration(config.SequencerMinBatchDelaySeconds) * time.Second)
	shutdown := ctx.Done()
	go func() {
		for {
			select {
			case <-timer.C:
				seqr.commit(ctx)
				timer.Reset(time.Duration(config.SequencerMinBatchDelaySeconds) * time.Second)
			case <-shutdown:
				timer.Stop()
				return
			}
		}
	}()

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

// Enqueue pushes the tx into the pending batch for the relevent Router
func (seqr *MemorySequencer) Enqueue(ctx context.Context,
	routerAddr common.Address,
	ownerAddr common.Address,
	actionData []string,
	actionSig string,
) (*model.ActionTransaction, error) {
	if routerAddr.Hex() == "" || len(actionData) == 0 || actionSig == "" {
		return nil, fmt.Errorf("invalid action data")
	}
	seqr.Lock()
	defer seqr.Unlock()

	pendingBatch, ok := seqr.pending[routerAddr.Hex()]
	if !ok {
		pendingBatch = &model.ActionBatch{
			ID:            uuid.NewV4().String(),
			Status:        model.ActionTransactionStatusPending,
			RouterAddress: routerAddr.Hex(),
		}
	}

	tx := &model.ActionTransaction{
		ID:            uuid.NewV4().String(),
		Owner:         ownerAddr.Hex(),
		Payload:       actionData,
		Sig:           actionSig,
		RouterAddress: routerAddr.Hex(),
		Batch:         pendingBatch,
	}
	pendingBatch.Transactions = append(pendingBatch.Transactions, tx)

	// do a dry run of the pending batch to see if it explodes
	// TODO: implement local chain simulation instead of this
	// TOOD: the commitBatch logic should be exluding actions with expired session
	dryRun := true
	_, err := seqr.commitBatch(ctx, routerAddr, pendingBatch, dryRun)
	if err != nil {
		return nil, fmt.Errorf("action will cause batch to fail: %v", err)
	}

	seqr.pending[routerAddr.Hex()] = pendingBatch
	seqr.emitTx(tx)

	seqr.log.Info().
		Str("batch", pendingBatch.ID).
		Str("tx", tx.ID).
		Msg("queued")

	return tx, nil
}

func (seqr *MemorySequencer) commit(ctx context.Context) {
	// TODO: do this without locking the enqueue process so much
	seqr.Lock()
	defer seqr.Unlock()

	seqr.processing = seqr.pending
	seqr.pending = map[string]*model.ActionBatch{}

	for routerAddr, batch := range seqr.processing {
		if len(batch.Transactions) == 0 {
			continue
		}
		rcpt, err := seqr.commitBatch(ctx, common.HexToAddress(routerAddr), batch, false)
		if rcpt != nil {
			block := int(uint32(rcpt.BlockNumber.Uint64()))
			batch.Block = &block
			txid := rcpt.TxHash.Hex()
			batch.Tx = &txid
		}
		if err != nil {
			seqr.log.Error().
				Err(err).
				Str("batch", batch.ID).
				Msg("commit")
			// move the batch into the failed pile
			// TODO: better handling of batch failure
			// TODO: better logging/metrics of batch failure
			// TODO: retries/removing bad actions from the batch
			batch.Status = model.ActionTransactionStatusFailed
			failures, ok := seqr.failure[batch.RouterAddress]
			if !ok {
				failures = []*model.ActionBatch{}
			}
			failures = append(failures, batch)
			seqr.failure[batch.RouterAddress] = failures
			seqr.processing = map[string]*model.ActionBatch{}
		} else {
			seqr.log.Info().
				Str("batch", batch.ID).
				Msg("commit")
			// move the batch into the success pile
			batch.Status = model.ActionTransactionStatusSuccess
			successes, ok := seqr.success[batch.RouterAddress]
			if !ok {
				successes = []*model.ActionBatch{}
			}
			successes = append(successes, batch)
			seqr.success[batch.RouterAddress] = successes
			seqr.processing = map[string]*model.ActionBatch{}
		}
		for _, tx := range batch.Transactions {
			seqr.emitTx(tx)
			seqr.log.Info().
				Str("batch", batch.ID).
				Str("tx", tx.ID).
				Str("status", string(batch.Status)).
				Msg("processed")
		}
	}
}

func (seqr *MemorySequencer) commitBatch(
	ctx context.Context,
	routerAddr common.Address,
	batch *model.ActionBatch,
	dryRun bool,
) (*types.Receipt, error) {

	actions := [][][]byte{}
	sigs := [][]byte{}
	for _, action := range batch.Transactions {
		actions = append(actions, action.ActionBytes())
		sigs = append(sigs, action.ActionSig())
	}

	client := seqr.httpClient
	client.Lock()
	defer client.Unlock()

	sessionRouter, err := router.NewSessionRouter(routerAddr, client)
	if err != nil {
		return nil, err
	}

	txOpts, err := client.NewRelayTransactor(ctx)
	if err != nil {
		return nil, err
	}

	txOpts.Context = ctx
	txOpts.Value = big.NewInt(0) // in wei
	// txOpts.GasLimit = uint64(3000000000)                                  // in units
	// txOpts.GasPrice = txOpts.GasPrice.Mul(txOpts.GasPrice, big.NewInt(2)) // up to double

	var tx *types.Transaction
	if dryRun {
		cabi, err := abi.JSON(strings.NewReader(router.SessionRouterABI))
		if err != nil {
			return nil, err
		}
		input, err := cabi.Pack("dispatch", actions, sigs)
		if err != nil {
			return nil, err
		}
		err = client.EstimateContractGas(ctx, txOpts, &routerAddr, input)
		if err != nil {
			return nil, fmt.Errorf("failed dry-run of batch tx: %v", err)
		}
		return nil, err
	} else {
		tx, err = sessionRouter.Dispatch(txOpts,
			actions,
			sigs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed commit batch tx: %v", err)
		}
		defer client.IncrementRelayNonce(ctx)
		// wait til batch success
		maxWait, cancel := context.WithTimeout(ctx, 1*time.Minute)
		defer cancel()
		rcpt, err := bind.WaitMined(maxWait, client, tx)
		if err != nil {
			return nil, err
		}
		switch rcpt.Status {
		case 1:
		default:
			return nil, fmt.Errorf("tx failed: %v", rcpt)
		}
		return rcpt, nil
	}

	// batch.Inc() // TODO: metrics
	// cost, _ := weiToGwei(tx.Cost()).Float64()
	// txCost.Observe(cost)

}

func (seqr *MemorySequencer) Signout(ctx context.Context, routerAddr common.Address, sessionKey common.Address, permit string) error {
	client := seqr.httpClient
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
	txOpts.Value = big.NewInt(0) // in wei
	// txOpts.GasLimit = uint64(3000000)                                     // in units
	// txOpts.GasPrice = txOpts.GasPrice.Mul(txOpts.GasPrice, big.NewInt(2)) // up to double

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

	// batch.Inc() // TODO: metrics
	// cost, _ := weiToGwei(tx.Cost()).Float64()
	// txCost.Observe(cost)

	return nil
}

func (seqr *MemorySequencer) Signin(ctx context.Context, routerAddr common.Address, dispatcherAddr common.Address, sessionKey common.Address, ttl uint32, scopes uint32, permit string) error {
	client := seqr.httpClient
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
	txOpts.Value = big.NewInt(0) // in wei
	// txOpts.GasLimit = uint64(3000000)                                     // in units
	// txOpts.GasPrice = txOpts.GasPrice.Mul(txOpts.GasPrice, big.NewInt(2)) // up to double

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

	// batch.Inc() // TODO: metrics
	// cost, _ := weiToGwei(tx.Cost()).Float64()
	// txCost.Observe(cost)

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
		return fmt.Errorf("tx failed: %v", rcpt)
	}

	return nil
}

func (seqr *MemorySequencer) GetBatch(routerAddr common.Address, id string) (*model.ActionBatch, error) {
	r := routerAddr.Hex()

	if batch, ok := seqr.pending[r]; ok {
		if batch.ID == id {
			return batch, nil
		}
	}
	if batch, ok := seqr.processing[r]; ok {
		if batch.ID == id {
			return batch, nil
		}
	}
	if batches, ok := seqr.success[r]; ok {
		for _, batch := range batches {
			if batch.ID == id {
				return batch, nil
			}
		}
	}
	if batches, ok := seqr.failure[r]; ok {
		for _, batch := range batches {
			if batch.ID == id {
				return batch, nil
			}
		}
	}
	return nil, nil
}

func (seqr *MemorySequencer) GetTransactions(routerAddr common.Address, owner *string, status []model.ActionTransactionStatus) ([]*model.ActionTransaction, error) {
	seqr.RLock()
	defer seqr.RUnlock()

	txs := []*model.ActionTransaction{}

	if len(status) == 0 || inStatusList(status, model.ActionTransactionStatusPending) {
		for _, batch := range seqr.pending {
			for _, tx := range batch.Transactions {
				if owner == nil || *owner == tx.Owner {
					txs = append(txs, tx)
				}
			}
		}
		for _, batch := range seqr.processing {
			for _, tx := range batch.Transactions {
				if owner == nil || *owner == tx.Owner {
					txs = append(txs, tx)
				}
			}
		}
	}
	if len(status) == 0 || inStatusList(status, model.ActionTransactionStatusSuccess) {
		for _, batches := range seqr.success {
			for _, batch := range batches {
				for _, tx := range batch.Transactions {
					if owner == nil || *owner == tx.Owner {
						txs = append(txs, tx)
					}
				}
			}
		}
	}
	if len(status) == 0 || inStatusList(status, model.ActionTransactionStatusFailed) {
		for _, batches := range seqr.failure {
			for _, batch := range batches {
				for _, tx := range batch.Transactions {
					if owner == nil || *owner == tx.Owner {
						txs = append(txs, tx)
					}
				}
			}
		}
	}

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
		{Type: byteArray},
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
