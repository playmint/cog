package model

import (
	"context"
	"sync"

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

type StateEvent struct {
	StateID string
	Event   Event
}

type Subscriptions struct {
	Events         map[string]map[uuid.UUID]chan Event
	TxByOwner      map[string]map[string]map[uuid.UUID]chan *ActionTransaction
	SessionByOwner map[string]map[string]map[uuid.UUID]chan *Session
	notifications  chan interface{}
	sync.RWMutex
}

func NewSubscriptions() (*Subscriptions, chan interface{}) {
	notifications := make(chan interface{}, NotificationBuffer)
	return &Subscriptions{
		Events:        map[string]map[uuid.UUID]chan Event{},
		TxByOwner:     map[string]map[string]map[uuid.UUID]chan *ActionTransaction{},
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
			case *StateEvent:
				for stateID, subs := range subs.Events {
					if stateID != obj.StateID {
						continue
					}
					if obj.Event == nil {
						continue
					}
					for _, subscriber := range subs {
						select {
						case subscriber <- obj.Event:
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

func (subs *Subscriptions) SubscribeStateEvent(ctx context.Context, stateID string) chan Event {
	id := uuid.New()

	go func() {
		<-ctx.Done()
		subs.Lock()
		defer subs.Unlock()
		chans, ok := subs.Events[stateID]
		if !ok {
			return
		}
		delete(chans, id)
	}()

	subs.Lock()
	chans, ok := subs.Events[stateID]
	if !ok {
		chans = map[uuid.UUID]chan Event{}
	}
	events := make(chan Event, SubscriptionBuffer)
	chans[id] = events
	subs.Events[stateID] = chans
	subs.Unlock()

	return events
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
