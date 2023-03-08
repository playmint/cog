package model

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	// SubscriptionBuffer is how many events will be buffered PER-CLIENT when
	// the client is being a slow consumer before they start getting dropped.
	// Too big a number on a busy sub will end up with a client processing a
	// big backlog slowly Too small a number and we may be dropping events
	// leaving the client out of sync.
	//
	// I picked "10" pretty much out of nowhere, so if you have arrived at this
	// code investigating dropped events, change this value
	SubscriptionBuffer = 10

	// NotificationBuffer is the total number of notifications that can be buffered
	// globally before they start getting dropped when ALL client consumers are being
	// too slow or if processing of the queue is taking too long.
	//
	// I picked "1000" pretty much out of nowhere, so if you have arrived at this
	// code investigating dropped events, change this value
	NotificationBuffer = 1000
)

type Subscriptions struct {
	State          map[string]map[uuid.UUID]chan *State
	TxByOwner      map[string]map[string]map[uuid.UUID]chan *ActionTransaction
	SessionByOwner map[string]map[string]map[uuid.UUID]chan *Session
	notifications  chan interface{}
	// Games map[string]map[uuid.UUID]chan *Game
	// Nodes map[string]map[uuid.UUID]chan *Node
	sync.RWMutex
}

func NewSubscriptions() (*Subscriptions, chan interface{}) {
	notifications := make(chan interface{}, NotificationBuffer)
	return &Subscriptions{
		State:     map[string]map[uuid.UUID]chan *State{},
		TxByOwner: map[string]map[string]map[uuid.UUID]chan *ActionTransaction{},
		// Games: map[string]map[uuid.UUID]chan *Game{},
		// Nodes: map[string]map[uuid.UUID]chan *Node{},
		notifications: notifications,
	}, notifications
}

func (subs *Subscriptions) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case notification := <-subs.notifications:
			switch obj := notification.(type) {
			case *State:
				for stateID, subs := range subs.State {
					if stateID != obj.ID {
						continue
					}
					for _, subscriber := range subs {
						select {
						case subscriber <- obj:
						default:
						}
					}
				}
			case *ActionTransaction:
				for routerID, subsByOwner := range subs.TxByOwner {
					if routerID != obj.RouterAddress {
						continue
					}
					for txOwner, subs := range subsByOwner {
						if txOwner != "" && txOwner != obj.Owner {
							continue
						}
						for _, subscriber := range subs {
							select {
							case subscriber <- obj:
							default:
							}
						}
					}
				}
			case *Session:
				for routerID, subsByOwner := range subs.SessionByOwner {
					if routerID != obj.RouterAddress {
						continue
					}
					for sessionOwner, subs := range subsByOwner {
						if sessionOwner != "" && sessionOwner != obj.Owner {
							continue
						}
						for _, subscriber := range subs {
							select {
							case subscriber <- obj:
							default:
							}
						}
					}
				}
			default:
			}
		}
	}
}

func (subs *Subscriptions) SubscribeState(ctx context.Context, stateID string) chan *State {
	id := uuid.New()

	go func() {
		<-ctx.Done()
		subs.Lock()
		defer subs.Unlock()
		chans, ok := subs.State[stateID]
		if !ok {
			return
		}
		delete(chans, id)
	}()

	subs.Lock()
	chans, ok := subs.State[stateID]
	if !ok {
		chans = map[uuid.UUID]chan *State{}
	}
	debouncedEvents := make(chan *State, SubscriptionBuffer)
	events := make(chan *State, SubscriptionBuffer)
	// subscribing to a state query is probably an anti-pattern
	// and may be removed in a future version.
	// it is too easy to subscribe to a complex query and end
	// up DOSing yourself as the subscription attempts to process
	// a flurry of state updates.
	//
	// Unlike most subscriptions (recv an event), state subscriptions
	// don't care about each change, they really only care about the
	// "most up to date change". So here we debounce the updates so
	// that if 100 updates come in within 1second it only results in
	// 1 update to the subscription.
	//
	// the fact we have to do this at all is a clue that this is
	// probably an antipattern.
	go func() {
		next := make(chan *State, 1)
		minWait := 500 * time.Millisecond
		maxWait := 3 * time.Second
		minTimeout := time.NewTicker(minWait)
		maxTimeout := time.NewTicker(maxWait)
		done := ctx.Done()
		for {
			select {
			case <-done:
				return
			case evt := <-events:
				select {
				case <-next:
				default:
				}
				next <- evt
				minTimeout.Reset(minWait)
			case <-maxTimeout.C:
				select {
				case evt := <-next:
					debouncedEvents <- evt
				default:
				}
			case <-minTimeout.C:
				select {
				case evt := <-next:
					debouncedEvents <- evt
				default:
				}
			}
		}
	}()
	chans[id] = events
	subs.State[stateID] = chans
	subs.Unlock()

	return debouncedEvents
}

func (subs *Subscriptions) SubscribeTransaction(ctx context.Context, routerID string, ownerFilter *string) chan *ActionTransaction {
	id := uuid.New()
	owner := ""
	if ownerFilter != nil {
		owner = *ownerFilter
	}

	go func() {
		<-ctx.Done()
		subs.Lock()
		defer subs.Unlock()
		chansByOwner, ok := subs.TxByOwner[routerID]
		if !ok {
			return
		}
		chans, ok := chansByOwner[owner]
		if !ok {
			return
		}
		delete(chans, id)
	}()

	subs.Lock()
	chansByOwner, ok := subs.TxByOwner[routerID]
	if !ok {
		chansByOwner = map[string]map[uuid.UUID]chan *ActionTransaction{}
	}
	chans, ok := chansByOwner[owner]
	if !ok {
		chans = map[uuid.UUID]chan *ActionTransaction{}
	}
	events := make(chan *ActionTransaction, SubscriptionBuffer)
	chans[id] = events
	chansByOwner[owner] = chans
	subs.TxByOwner[routerID] = chansByOwner
	subs.Unlock()

	return events
}

func (subs *Subscriptions) SubscribeSession(ctx context.Context, routerID string, ownerFilter *string) chan *Session {
	id := uuid.New()
	owner := ""
	if ownerFilter != nil {
		owner = *ownerFilter
	}

	go func() {
		<-ctx.Done()
		subs.Lock()
		defer subs.Unlock()
		chansByOwner, ok := subs.SessionByOwner[routerID]
		if !ok {
			return
		}
		chans, ok := chansByOwner[owner]
		if !ok {
			return
		}
		delete(chans, id)
	}()

	subs.Lock()
	chansByOwner, ok := subs.SessionByOwner[routerID]
	if !ok {
		chansByOwner = map[string]map[uuid.UUID]chan *Session{}
	}
	chans, ok := chansByOwner[owner]
	if !ok {
		chans = map[uuid.UUID]chan *Session{}
	}
	events := make(chan *Session, SubscriptionBuffer)
	chans[id] = events
	chansByOwner[owner] = chans
	subs.SessionByOwner[routerID] = chansByOwner
	subs.Unlock()

	return events
}
