package sequencer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
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
	"github.com/playmint/ds-node/pkg/config"
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
	PrivateKey *ecdsa.PrivateKey
	// pending/processing are actions, grouped by Router that are waiting to join a batch
	pending    map[string]*model.ActionBatch
	processing map[string]*model.ActionBatch
	// batches are lists of actions, grouped by Router has processed
	success       map[string][]*model.ActionBatch
	failure       map[string][]*model.ActionBatch
	httpClient    *alchemy.Client
	simHttpClient *alchemy.Client
	simWsClient   *alchemy.Client
	simBlock      uint64
	notifications chan interface{}
	idxr          indexer.Indexer
	log           zerolog.Logger
	sync.RWMutex
}

func NewMemorySequencer(ctx context.Context, key *ecdsa.PrivateKey, notifications chan interface{}, httpProviderURL string, simProviderURL string, idxr indexer.Indexer) (*MemorySequencer, error) {
	var err error
	seqr := &MemorySequencer{
		PrivateKey:    key,
		notifications: notifications,
		log:           log.With().Str("service", "sequencer").Logger(),
		idxr:          idxr,
	}
	// setup an RPC client
	seqr.httpClient, err = alchemy.Dial(
		httpProviderURL,
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
	timer := time.NewTimer(time.Duration(1 * time.Second))
	shutdown := ctx.Done()
	go func() {
		for {
			select {
			case <-timer.C:
				seqr.commit(ctx)
				timer.Reset(time.Duration(10 * time.Second))
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
	stateAddr common.Address,
	actionData []string,
	actionSig string,
) (*model.ActionTransaction, error) {
	if routerAddr.Hex() == "" || len(actionData) == 0 || actionSig == "" {
		return nil, fmt.Errorf("invalid action data")
	}
	seqr.Lock()
	defer seqr.Unlock()

	actualPending, ok := seqr.pending[routerAddr.Hex()]
	if !ok {
		actualPending = &model.ActionBatch{
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
		Batch:         actualPending,
	}

	if seqr.simHttpClient != nil && seqr.simBlock > 0 {
		_, err := seqr.commitTxWithClient(ctx, routerAddr, tx, seqr.simHttpClient)
		if err != nil {
			seqr.log.Error().
				Err(err).
				Uint64("block", seqr.simBlock).
				Str("batch", "sim").
				Msg("sim-tx-fail")
			return nil, err
		}
		seqr.log.Info().
			Uint64("block", seqr.simBlock).
			Str("batch", "sim").
			Msg("sim-tx-ok")
	} else {
		seqr.log.Info().
			Uint64("block", seqr.simBlock).
			Str("batch", "sim").
			Msg("sim-not-ready")
	}

	actualPending.Transactions = append(actualPending.Transactions, tx)
	seqr.pending[routerAddr.Hex()] = actualPending
	seqr.emitTx(tx)

	seqr.log.Info().
		Str("batch", actualPending.ID).
		Str("tx", tx.ID).
		Msg("queued")

	return tx, nil
}

func (seqr *MemorySequencer) commit(ctx context.Context) {
	seqr.Lock()
	seqr.processing = seqr.pending
	seqr.pending = map[string]*model.ActionBatch{}
	seqr.Unlock()

	latestBlock, _ := seqr.httpClient.BlockNumber(ctx)

	todo := 0
	for _, batch := range seqr.processing {
		if len(batch.Transactions) == 0 {
			continue
		}
		todo++
	}
	if todo == 0 {
		if _, err := seqr.httpClient.MineEmptyBlock(ctx); err != nil {
			seqr.log.Error().
				Err(err).
				Msg("mine-empty-block")
		} else {
			seqr.log.Info().
				Msg("mine-empty-block")
			latestBlock, err = seqr.httpClient.BlockNumber(ctx)
			if err != nil {
				seqr.log.Error().
					Err(err).
					Msg("mine-empty-block-get-block")
				return
			}
		}
	}

	for routerAddr, batch := range seqr.processing {
		if len(batch.Transactions) == 0 {
			continue
		}
		for _, action := range batch.Transactions {
			rcpt, err := seqr.commitTxWithClient(ctx, common.HexToAddress(routerAddr), action, seqr.httpClient)
			if rcpt != nil {
				block := int(uint32(rcpt.BlockNumber.Uint64()))
				batch.Block = &block
				txid := rcpt.TxHash.Hex()
				batch.Tx = &txid
				if rcpt.BlockNumber.Uint64() > latestBlock {
					latestBlock = rcpt.BlockNumber.Uint64()
				}
			}
			if err != nil {
				seqr.log.Error().
					Err(err).
					Str("batch", batch.ID).
					Msg("commit")
			}
			seqr.log.Info().
				Str("batch", batch.ID).
				Msg("commit")
		}
		// move the batch into the success pile
		batch.Status = model.ActionTransactionStatusSuccess
		successes, ok := seqr.success[batch.RouterAddress]
		if !ok {
			successes = []*model.ActionBatch{}
		}
		successes = append(successes, batch)
		seqr.success[batch.RouterAddress] = successes
		for _, tx := range batch.Transactions {
			seqr.emitTx(tx)
			seqr.log.Info().
				Str("batch", batch.ID).
				Str("tx", tx.ID).
				Str("status", string(batch.Status)).
				Msg("processed")
		}
	}
	seqr.Lock()
	seqr.processing = map[string]*model.ActionBatch{}
	if latestBlock > seqr.simBlock {
		seqr.simBlock = latestBlock
		// refork the simulation from last block
		seqr.log.Info().
			Uint64("block", seqr.simBlock).
			Msg("reset-sim")
		if err := seqr.resetSim(ctx, seqr.simBlock); err != nil {
			seqr.log.Error().
				Err(err).
				Uint64("block", seqr.simBlock).
				Msg("reset-sim-fail")
		} else {
			// fork the simulated index
			simxr, err := seqr.idxr.NewSim(ctx, seqr.simBlock, seqr.simHttpClient, seqr.simWsClient)
			if err != nil {
				seqr.log.Error().
					Err(err).
					Uint64("block", seqr.simBlock).
					Msg("reset-idxsim-fail")
			} else {
				// fast forward the pending queue
				for routerAddr, batch := range seqr.pending {
					for _, action := range batch.Transactions {
						_, err := seqr.commitTxWithClient(ctx, common.HexToAddress(routerAddr), action, seqr.simHttpClient)
						if err != nil {
							seqr.log.Error().
								Err(err).
								Str("batch", batch.ID).
								Msg("fast-fwd-sim-fail")
						}
					}
				}
				// hotswap the indexer
				seqr.idxr.SetSim(simxr)
				seqr.log.Info().
					Uint64("block", seqr.simBlock).
					Msg("reset-sim-ok")
			}
		}
	}

	seqr.Unlock()
}

func (seqr *MemorySequencer) resetSim(ctx context.Context, latest uint64) (err error) {
	seqr.log.Info().Msg("dial-client")
	// new simulation client
	seqr.simHttpClient, err = alchemy.Dial(
		config.SimulationProviderHTTP,
		config.SequencerMaxConcurrency,
		config.SequencerPrivateKey,
	)
	if err != nil {
		return err
	}
	seqr.simWsClient, err = alchemy.Dial(
		config.SimulationProviderWS,
		config.SequencerMaxConcurrency,
		config.SequencerPrivateKey,
	)
	if err != nil {
		return err
	}
	// fork sim chain
	seqr.log.Info().Msg("start-reset")
	x, err := seqr.simHttpClient.Reset(ctx, config.SequencerProviderHTTP, latest)
	if err != nil {
		return err
	}
	seqr.log.Info().Msg("enable-automine")
	_, err = seqr.simHttpClient.EnableAutoMine(ctx)
	if err != nil {
		return err
	}
	seqr.log.Info().Msgf("end-reset: url:%v, block:%v x:%v", config.SequencerProviderHTTP, latest, x)
	return nil
}

func (seqr *MemorySequencer) commitTxWithClient(
	ctx context.Context,
	routerAddr common.Address,
	action *model.ActionTransaction,
	client *alchemy.Client,
) (*types.Receipt, error) {
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
	txOpts.Value = big.NewInt(0)         // in wei
	txOpts.GasLimit = uint64(3000000000) // in units

	tx, err := sessionRouter.Dispatch(txOpts,
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
		reason := errorReason(ctx, client, client.Address(), tx, rcpt.BlockNumber)
		return rcpt, fmt.Errorf("%s", reason)
	}
	return rcpt, nil
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

	if rcpt.BlockNumber.Uint64() > seqr.simBlock {
		seqr.Lock()
		seqr.simBlock = rcpt.BlockNumber.Uint64()
		seqr.log.Info().
			Uint64("block", seqr.simBlock).
			Msg("reset-sim")
		if err := seqr.resetSim(ctx, seqr.simBlock); err != nil {
			seqr.log.Error().
				Err(err).
				Uint64("block", seqr.simBlock).
				Msg("reset-sim-fail")
		} else {
			// fork the simulated index
			simxr, err := seqr.idxr.NewSim(ctx, seqr.simBlock, seqr.simHttpClient, seqr.simWsClient)
			if err != nil {
				seqr.log.Error().
					Err(err).
					Uint64("block", seqr.simBlock).
					Msg("reset-idxsim-fail")
			} else {
				// fast forward the pending queue
				for routerAddr, batch := range seqr.pending {
					for _, action := range batch.Transactions {
						_, err := seqr.commitTxWithClient(ctx, common.HexToAddress(routerAddr), action, seqr.simHttpClient)
						if err != nil {
							seqr.log.Error().
								Err(err).
								Str("batch", batch.ID).
								Msg("fast-fwd-sim-fail")
						}
					}
				}
				// hotswap the indexer
				seqr.idxr.SetSim(simxr)
				seqr.log.Info().
					Uint64("block", seqr.simBlock).
					Msg("reset-sim-ok")
			}
		}
		seqr.Unlock()
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
