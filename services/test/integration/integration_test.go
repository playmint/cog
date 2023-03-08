package integration_test

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API QuickTest", Ordered, func() {

	getQueryEndpoint := func() string {
		addr := os.Getenv("QUERY_ENDPOINT")
		if addr == "" {
			return "http://localhost:3080/query"
		}
		return addr
	}

	var (
		httpClient = &http.Client{
			Timeout: 15 * time.Second,
		}
		client      = graphql.NewClient(getQueryEndpoint(), httpClient)
		pollTimeout = 10 * time.Second
	)

	var (
		alicePrivateKey, _ = crypto.HexToECDSA("8451fe83595a88521cac265c1356568ab84c2b5d81654f36ea309be28ded5a0d")
		alicePublicKey     = "0x4208a6518500E980ED44Da94ea31b85c85ec4568"
		sessionPrivateKey  = newPrivateKey()
		sessionPublicKey   = publicAddress(sessionPrivateKey)
	)

	var (
		gameID = "latest"
	)

	type Arg struct {
		Kind  string
		Value interface{}
	}

	encodeAction := func(name string, args ...Arg) []byte {
		actionArgumentKinds := abi.Arguments{}
		actionArgumentValues := []interface{}{}
		for _, arg := range args {
			t, _ := abi.NewType(arg.Kind, arg.Kind, nil)
			actionArgumentKinds = append(actionArgumentKinds, abi.Argument{Type: t})
			actionArgumentValues = append(actionArgumentValues, arg.Value)
		}
		m := abi.NewMethod(name, name, abi.Function, "", false, false, actionArgumentKinds, nil)
		b, err := actionArgumentKinds.Pack(actionArgumentValues...)
		Expect(err).ToNot(HaveOccurred())
		return append(m.ID[:4], b...)
	}

	packActions := func(actions [][]byte) []byte {
		byteArray, _ := abi.NewType("bytes[]", "bytes[]", nil)
		args := abi.Arguments{
			{Type: byteArray},
		}
		b, err := args.Pack(actions)
		Expect(err).ToNot(HaveOccurred())
		return b
	}

	dispatchSigned := func(ctx context.Context, key *ecdsa.PrivateKey, actionName string, actionArgs ...Arg) (*dispatchResponse, error) {
		// abi encode the action into action bundle
		action := encodeAction(actionName, actionArgs...)
		actions := [][]byte{action}
		// sign the bundle with the session key
		authMessage := crypto.Keccak256Hash(
			[]byte("\x19Ethereum Signed Message:\n32"),
			crypto.Keccak256Hash(packActions(actions)).Bytes(),
		)
		sig, err := crypto.Sign(authMessage.Bytes(), key)
		Expect(err).ToNot(HaveOccurred())
		sig[len(sig)-1] += 27
		// send mutation
		return dispatch(
			ctx, client,
			gameID,
			[]string{hexutil.Encode(action)},
			hexutil.Encode(sig),
		)
	}

	sessionsCountByOwner := func() int {
		res, err := getSessionsByOwner(context.TODO(), client,
			gameID,
			alicePublicKey, // owner address
		)
		Expect(err).ToNot(HaveOccurred())
		return len(res.Game.Router.Sessions)
	}

	sessionOwner := func() string {
		res, err := getSessionByID(context.TODO(), client,
			gameID,
			sessionPublicKey, // session address
		)
		Expect(err).ToNot(HaveOccurred())
		return res.Game.Router.Session.Owner
	}

	transactionStatus := func(txid string) func() ActionTransactionStatus {
		return func() ActionTransactionStatus {
			res, err := getTransactionByID(context.TODO(), client,
				gameID,
				txid,
			)
			Expect(err).ToNot(HaveOccurred())
			return res.Game.Router.Transaction.Status
		}
	}

	// depracted
	// It("should allow alice to signup for a player account", func(ctx SpecContext) {
	// 	// sign a signup mutation with alice's private key
	// 	res, err := signup(
	// 		ctx, client,
	// 		gameID,
	// 		"ignored-as-this-is-a-noop-at-the-moment",
	// 	)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(res.Signup).To(BeTrue())
	// })

	It("should eventully fetch the deployed game", func(ctx SpecContext) {
		// we start with this check with a large timeout to help avoid the race
		// between starting up the supporting services/contracts and running the tests
		getGameID := func() string {
			ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()
			res, err := getGame(ctx, client, gameID)
			if err != nil {
				return ""
			}
			return res.Game.Id
		}
		Eventually(getGameID).
			Within(30 * time.Second).
			ProbeEvery(1 * time.Second).
			ShouldNot(BeEmpty())
	})

	It("should authorize a session key for alice's account", func(ctx SpecContext) {
		// build a session auth message
		sessionAddr := common.HexToAddress(sessionPublicKey)
		signinMessage := []byte("You are signing in with session: ")
		// sign it
		sig, err := crypto.Sign(crypto.Keccak256Hash(
			[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(signinMessage)+20)),
			signinMessage,
			sessionAddr.Bytes(),
		).Bytes(), alicePrivateKey)
		Expect(err).ToNot(HaveOccurred())
		sig[len(sig)-1] += 27
		// submit signin request
		res, err := signin(
			ctx, client,
			gameID, // gameID
			sessionPublicKey,
			hexutil.Encode(sig),
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Signin).To(BeTrue())
	})

	It("should eventully fetch the newly authorized session", func(ctx SpecContext) {
		Eventually(sessionOwner, pollTimeout).Should(Equal(alicePublicKey))
	})

	It("should not have at least one session by owner", func(ctx SpecContext) {
		Eventually(sessionsCountByOwner, pollTimeout).Should(BeNumerically(">", 0))
	})

	It("should send a session signed RESET_MAP action via dispatch", func(ctx SpecContext) {
		res, err := dispatchSigned(ctx, sessionPrivateKey, "RESET_MAP")
		time.Sleep(5 * time.Second)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Dispatch.Id).ToNot(BeEmpty())
		Expect(res.Dispatch.Status).To(Equal(ActionTransactionStatusPending))
		Eventually(transactionStatus(res.Dispatch.Id), pollTimeout).Should(Equal(ActionTransactionStatusSuccess))
	})

	prevTransactionBlock := 0

	It("should send a session signed SPAWN_SEEKER action via dispatch", func(ctx SpecContext) {
		res, err := dispatchSigned(ctx, sessionPrivateKey, "SPAWN_SEEKER",
			Arg{"uint32", uint32(1)},
			Arg{"uint8", uint8(0)},
			Arg{"uint8", uint8(0)},
			Arg{"uint8", uint8(33)},
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Dispatch.Id).ToNot(BeEmpty())
		Expect(res.Dispatch.Status).To(Equal(ActionTransactionStatusPending))
		Eventually(transactionStatus(res.Dispatch.Id), pollTimeout).Should(Equal(ActionTransactionStatusSuccess))

		// grab the block number the batch was mined at so we can use it later
		tx, err := getTransactionByID(ctx, client, gameID, res.Dispatch.Id)
		Expect(tx.Game.Router.Transaction.Batch.Tx).ToNot(BeEmpty())
		Expect(err).ToNot(HaveOccurred())
		prevTransactionBlock = tx.Game.Router.Transaction.Batch.Block
		Expect(prevTransactionBlock).ToNot(BeNil())
	})

	It("should send a session signed REVEAL_SEED action via dispatch", func(ctx SpecContext) {
		res, err := dispatchSigned(ctx, sessionPrivateKey, "REVEAL_SEED",
			Arg{"uint32", uint32(prevTransactionBlock)}, // seedid
			Arg{"uint32", uint32(42)},                   // entropy
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Dispatch.Id).ToNot(BeEmpty())
		Expect(res.Dispatch.Status).To(Equal(ActionTransactionStatusPending))
		Eventually(transactionStatus(res.Dispatch.Id), pollTimeout).Should(Equal(ActionTransactionStatusSuccess))
	})

	It("should send a session signed MOVE_SEEKER action via dispatch", func(ctx SpecContext) {
		res, err := dispatchSigned(ctx, sessionPrivateKey, "MOVE_SEEKER",
			Arg{"uint32", uint32(1)}, // seekerid
			Arg{"uint8", uint8(0)},   // NORTH enum
		)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Dispatch.Id).ToNot(BeEmpty())
		Expect(res.Dispatch.Status).To(Equal(ActionTransactionStatusPending))
		Eventually(transactionStatus(res.Dispatch.Id), pollTimeout).Should(Equal(ActionTransactionStatusSuccess))
	})

	It("should fetch seeker location", func(ctx SpecContext) {
		res, err := getSeekers(ctx, client, gameID)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.Game.State.Seekers).To(HaveLen(1))
		Expect(res.Game.State.Seekers[0].Kind).To(Equal("Seeker"))
		Expect(len(res.Game.State.Seekers[0].Position.Keys)).To(BeNumerically(">", 1))
		Expect(res.Game.State.Seekers[0].Position.Keys[0]).To(EqualBig(0))
		Expect(res.Game.State.Seekers[0].Position.Keys[1]).To(EqualBig(1))
	})

	It("should reject bad signatures during dispatch", func(ctx SpecContext) {
		// abi encode an action
		action := encodeAction("EXAMPLE_ACTION")
		authMessage := crypto.Keccak256Hash(
			[]byte("\x19Ethereum Signed Message:\n32"),
			crypto.Keccak256Hash(action).Bytes(),
		)
		// sign it with the session key
		sig, err := crypto.Sign(authMessage.Bytes(), sessionPrivateKey)
		Expect(err).ToNot(HaveOccurred())
		sig[len(sig)-1] += 27
		// break the signature by changing the payload
		action = encodeAction("HACKED_ACTION")
		// send mutation
		_, err = dispatch(
			ctx, client,
			gameID,
			[]string{hexutil.Encode(action)},
			hexutil.Encode(sig),
		)
		Expect(err).To(MatchError(ContainSubstring("no session for signer or invalid signature")))
	})

})

func newPrivateKey() *ecdsa.PrivateKey {
	key, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func publicAddress(privateKey *ecdsa.PrivateKey) string {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address.Hex()
}
