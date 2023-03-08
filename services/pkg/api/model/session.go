package model

type Session struct {
	ID            string        `json:"id"`
	Owner         string        `json:"owner"`
	Scope         *SessionScope `json:"scope"`
	Expires       int           `json:"expires"`
	RouterAddress string
}
