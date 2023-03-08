package model

import "github.com/ethereum/go-ethereum/common"

type Game struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`

	DispatcherAddress common.Address
	StateAddress      common.Address
	RouterAddress     common.Address

	*Graph
}

func (game *Game) Dispatcher() *Dispatcher {
	return &Dispatcher{
		ID: game.DispatcherAddress.Hex(),
	}
}

func (game *Game) State() *State {
	return &State{
		ID: game.StateAddress.Hex(),
	}
}

func (game *Game) Router() *Router {
	return &Router{
		ID: game.RouterAddress.Hex(),
	}
}
