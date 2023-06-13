package cog

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/benbjohnson/immutable"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/playmint/ds-node/pkg/api/model"
	"github.com/playmint/ds-node/pkg/client"
	"github.com/playmint/ds-node/pkg/contracts/router"
	"github.com/playmint/ds-node/pkg/indexer/eventwatcher"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var FULL_ACCESS uint32 = 0xffffffff

type SessionStore struct {
	sessions      *immutable.Map[string, *immutable.Map[string, *model.Session]]
	abi           *abi.ABI
	events        *eventwatcher.Watcher
	notifications chan interface{}
	log           zerolog.Logger
	sync.RWMutex
}

func NewSessionStore(ctx context.Context, watcher *eventwatcher.Watcher, notifications chan interface{}) (*SessionStore, error) {
	cabi, err := abi.JSON(strings.NewReader(router.SessionRouterABI))
	if err != nil {
		return nil, err
	}
	store := &SessionStore{
		abi:           &cabi,
		events:        watcher,
		notifications: notifications,
		sessions:      immutable.NewMap[string, *immutable.Map[string, *model.Session]](nil),
		log:           log.With().Str("service", "indexer").Str("component", "sessionstore").Logger(),
	}

	// watch all events from all contracts that match the SessionCreate topic
	query := [][]interface{}{{cabi.Events["SessionCreate"].ID}}
	topics, err := abi.MakeTopics(query...)
	if err != nil {
		return nil, err
	}
	queue := watcher.SubscribeTopic(topics[0])

	go store.watch(ctx, queue)
	return store, nil
}

func (rs *SessionStore) emitSession(s *model.Session) {
	rs.notifications <- s
}

func (rs *SessionStore) watch(ctx context.Context, events chan types.Log) {
	for {
		select {
		case <-ctx.Done():
			return
		case rawEvent := <-events:
			eventABI, err := rs.abi.EventByID(rawEvent.Topics[0])
			if err != nil {
				rs.log.Warn().Err(err).Msg("unhandleable event topic")
				continue
			}
			rs.log.Debug().Msgf("recv %v", eventABI.RawName)
			switch eventABI.RawName {
			case "SessionCreate":
				var evt router.SessionRouterSessionCreate
				if err := unpackLog(rs.abi, &evt, eventABI.RawName, rawEvent); err != nil {
					rs.log.Warn().Err(err).Msgf("undecodable %T event", evt)
					continue
				}
				evt.Raw = rawEvent
				err := rs.setSession(&evt)
				if err != nil && client.IsRetryable(err) {
					time.Sleep(1 * time.Second)
					events <- rawEvent
				} else if err != nil {
					rs.log.Error().Err(err).Msgf("failed process %T event", evt)
				}
			default:
				rs.log.Warn().Msgf("ignoring unhandled event type %v", eventABI)
			}
		}
	}
}

// An action was registered, update the mapping for the game
func (rs *SessionStore) setSession(evt *router.SessionRouterSessionCreate) error {
	rs.Lock()
	defer rs.Unlock()

	if evt.Raw.Removed {
		// hmmm blockchain reorg occured so what do we do?
		// just ignore for now, but this probably needs more thought
		return nil
	}

	// create new session object
	session := &model.Session{
		ID:    evt.Session.Hex(),
		Owner: evt.Owner.Hex(),
		Scope: &model.SessionScope{
			FullAccess: (evt.Scopes & FULL_ACCESS) == FULL_ACCESS,
		},
		Expires:       int(evt.Exp),
		RouterAddress: evt.Raw.Address.Hex(),
	}

	// fetch existing sessions for this router
	sessions, sessionsExist := rs.sessions.Get(evt.Raw.Address.Hex())
	if !sessionsExist {
		sessions = immutable.NewMap[string, *model.Session](nil)
	}

	// add to the set
	sessions = sessions.Set(evt.Session.Hex(), session)
	rs.sessions = rs.sessions.Set(evt.Raw.Address.Hex(), sessions)

	rs.emitSession(session)

	return nil
}

func (rs *SessionStore) GetSession(routerAddr common.Address, sessionID string) *model.Session {
	rs.RLock()
	defer rs.RUnlock()

	sessions, ok := rs.sessions.Get(routerAddr.Hex())
	if !ok {
		return nil
	}
	session, ok := sessions.Get(sessionID)
	if !ok {
		return nil
	}
	return session
}

func (rs *SessionStore) GetSessions(routerAddr common.Address, owner *string) []*model.Session {
	rs.RLock()
	defer rs.RUnlock()

	sessions := []*model.Session{}

	routerSessions, ok := rs.sessions.Get(routerAddr.Hex())
	if !ok {
		return nil
	}
	itr := routerSessions.Iterator()
	for !itr.Done() {
		_, session, ok := itr.Next()
		if !ok {
			continue
		}
		if owner != nil && *owner != session.Owner {
			continue
		}
		sessions = append(sessions, session)
	}
	return sessions
}
