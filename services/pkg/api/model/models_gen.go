// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Event interface {
	IsEvent()
}

type Account struct {
	ID string `json:"id"`
}

// annotations are off-chain data attached to nodes that are guarenteed
// to have been made available to all clients, but are not usable within logic.
// for example; a "name" might be an annotation because there is no logic on-chain
// that expects to verify the name OR a node might be annotated with some JSON
// meatadata containing static details for a client to display. Since values
// are stored only in calldata (not in state storage), values can be larger than
// usually cost-effective for an equivilent value stored in state.
type Annotation struct {
	ID    string `json:"id"`
	Ref   string `json:"ref"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type BlockEvent struct {
	ID        string `json:"id"`
	Block     int    `json:"block"`
	Simulated bool   `json:"simulated"`
}

func (BlockEvent) IsEvent() {}

type ContractConfig struct {
	Name    string `json:"name"`
	ChainID int    `json:"chainId"`
	Address string `json:"address"`
}

type Dispatcher struct {
	ID string `json:"id"`
}

type ERC721Attribute struct {
	DisplayType string `json:"display_type"`
	TraitType   string `json:"trait_type"`
	Value       string `json:"value"`
}

type ERC721Metadata struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Image           string             `json:"image"`
	ImageData       string             `json:"image_data"`
	BackgroundColor string             `json:"background_color"`
	AnimationURL    string             `json:"animation_url"`
	YoutubeURL      string             `json:"youtube_url"`
	ExternalURL     string             `json:"external_url"`
	Attributes      []*ERC721Attribute `json:"attributes"`
}

// match condition for traversing/filtering the graph.
type Match struct {
	// ids only match if node is any of these ids, if empty match any id
	Ids []string `json:"ids"`
	// via only follow edges of these rel types, if empty follow all edges
	Via []*RelMatch `json:"via"`
	// kinds only matches if node kind is any of these kinds, if empty match any kind
	Kinds []string `json:"kinds"`
	// has only matches nodes that directly have the Rel, similar to via but subtle difference.
	// given the graph...
	//
	// 	A --HAS_RED--> B --HAS_BLUE--> C
	// 	A --HAS_BLUE--> Y --HAS_RED--> Z
	//
	// match(via: ["HAS_RED", "HAS_BLUE"]) would return B,C,Y,Z
	// match(via: ["HAS_RED", "HAS_BLUE"], has: ["HAS_RED"]) would return B,Z
	Has []*RelMatch `json:"has"`
	// `limit` stops matches after that many edges have been collected
	Limit *int `json:"limit"`
	// how many connections of connections allow to follow when searching
	// for a match. default=0 (meaning only direct connections)
	MaxDepth *int `json:"maxDepth"`
}

// RelMatch configures the types of edges that can be matched.
//
// rel is the human friendly name of the relationship.
//
// dir is either IN/OUT/BOTH and ditactes if we consider the edge pointing in an
// outbound or inbound direction from this node.
type RelMatch struct {
	Rel string             `json:"rel"`
	Dir *RelMatchDirection `json:"dir"`
	Key *int               `json:"key"`
}

type Router struct {
	ID           string               `json:"id"`
	Sessions     []*Session           `json:"sessions"`
	Session      *Session             `json:"session"`
	Transactions []*ActionTransaction `json:"transactions"`
	Transaction  *ActionTransaction   `json:"transaction"`
}

type SessionScope struct {
	FullAccess bool `json:"FullAccess"`
}

type State struct {
	ID        string `json:"id"`
	Block     int    `json:"block"`
	Simulated bool   `json:"simulated"`
	// nodes returns any nodes that match the Match filter.
	Nodes []*Node `json:"nodes"`
	// node returns the first node that mates the Match filter.
	Node *Node `json:"node"`
}

type ActionTransactionStatus string

const (
	ActionTransactionStatusUnknown ActionTransactionStatus = "UNKNOWN"
	ActionTransactionStatusPending ActionTransactionStatus = "PENDING"
	ActionTransactionStatusSuccess ActionTransactionStatus = "SUCCESS"
	ActionTransactionStatusFailed  ActionTransactionStatus = "FAILED"
)

var AllActionTransactionStatus = []ActionTransactionStatus{
	ActionTransactionStatusUnknown,
	ActionTransactionStatusPending,
	ActionTransactionStatusSuccess,
	ActionTransactionStatusFailed,
}

func (e ActionTransactionStatus) IsValid() bool {
	switch e {
	case ActionTransactionStatusUnknown, ActionTransactionStatusPending, ActionTransactionStatusSuccess, ActionTransactionStatusFailed:
		return true
	}
	return false
}

func (e ActionTransactionStatus) String() string {
	return string(e)
}

func (e *ActionTransactionStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ActionTransactionStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ActionTransactionStatus", str)
	}
	return nil
}

func (e ActionTransactionStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type AttributeKind string

const (
	AttributeKindBool         AttributeKind = "BOOL"
	AttributeKindInt8         AttributeKind = "INT8"
	AttributeKindInt16        AttributeKind = "INT16"
	AttributeKindInt32        AttributeKind = "INT32"
	AttributeKindInt64        AttributeKind = "INT64"
	AttributeKindInt128       AttributeKind = "INT128"
	AttributeKindInt256       AttributeKind = "INT256"
	AttributeKindInt          AttributeKind = "INT"
	AttributeKindUINt8        AttributeKind = "UINT8"
	AttributeKindUINt16       AttributeKind = "UINT16"
	AttributeKindUINt32       AttributeKind = "UINT32"
	AttributeKindUINt64       AttributeKind = "UINT64"
	AttributeKindUINt128      AttributeKind = "UINT128"
	AttributeKindUINt256      AttributeKind = "UINT256"
	AttributeKindBytes        AttributeKind = "BYTES"
	AttributeKindString       AttributeKind = "STRING"
	AttributeKindAddress      AttributeKind = "ADDRESS"
	AttributeKindBytes4       AttributeKind = "BYTES4"
	AttributeKindBoolArray    AttributeKind = "BOOL_ARRAY"
	AttributeKindInt8Array    AttributeKind = "INT8_ARRAY"
	AttributeKindInt16Array   AttributeKind = "INT16_ARRAY"
	AttributeKindInt32Array   AttributeKind = "INT32_ARRAY"
	AttributeKindInt64Array   AttributeKind = "INT64_ARRAY"
	AttributeKindInt128Array  AttributeKind = "INT128_ARRAY"
	AttributeKindInt256Array  AttributeKind = "INT256_ARRAY"
	AttributeKindIntArray     AttributeKind = "INT_ARRAY"
	AttributeKindUINt8Array   AttributeKind = "UINT8_ARRAY"
	AttributeKindUINt16Array  AttributeKind = "UINT16_ARRAY"
	AttributeKindUINt32Array  AttributeKind = "UINT32_ARRAY"
	AttributeKindUINt64Array  AttributeKind = "UINT64_ARRAY"
	AttributeKindUINt128Array AttributeKind = "UINT128_ARRAY"
	AttributeKindUINt256Array AttributeKind = "UINT256_ARRAY"
	AttributeKindBytesArray   AttributeKind = "BYTES_ARRAY"
	AttributeKindStringArray  AttributeKind = "STRING_ARRAY"
)

var AllAttributeKind = []AttributeKind{
	AttributeKindBool,
	AttributeKindInt8,
	AttributeKindInt16,
	AttributeKindInt32,
	AttributeKindInt64,
	AttributeKindInt128,
	AttributeKindInt256,
	AttributeKindInt,
	AttributeKindUINt8,
	AttributeKindUINt16,
	AttributeKindUINt32,
	AttributeKindUINt64,
	AttributeKindUINt128,
	AttributeKindUINt256,
	AttributeKindBytes,
	AttributeKindString,
	AttributeKindAddress,
	AttributeKindBytes4,
	AttributeKindBoolArray,
	AttributeKindInt8Array,
	AttributeKindInt16Array,
	AttributeKindInt32Array,
	AttributeKindInt64Array,
	AttributeKindInt128Array,
	AttributeKindInt256Array,
	AttributeKindIntArray,
	AttributeKindUINt8Array,
	AttributeKindUINt16Array,
	AttributeKindUINt32Array,
	AttributeKindUINt64Array,
	AttributeKindUINt128Array,
	AttributeKindUINt256Array,
	AttributeKindBytesArray,
	AttributeKindStringArray,
}

func (e AttributeKind) IsValid() bool {
	switch e {
	case AttributeKindBool, AttributeKindInt8, AttributeKindInt16, AttributeKindInt32, AttributeKindInt64, AttributeKindInt128, AttributeKindInt256, AttributeKindInt, AttributeKindUINt8, AttributeKindUINt16, AttributeKindUINt32, AttributeKindUINt64, AttributeKindUINt128, AttributeKindUINt256, AttributeKindBytes, AttributeKindString, AttributeKindAddress, AttributeKindBytes4, AttributeKindBoolArray, AttributeKindInt8Array, AttributeKindInt16Array, AttributeKindInt32Array, AttributeKindInt64Array, AttributeKindInt128Array, AttributeKindInt256Array, AttributeKindIntArray, AttributeKindUINt8Array, AttributeKindUINt16Array, AttributeKindUINt32Array, AttributeKindUINt64Array, AttributeKindUINt128Array, AttributeKindUINt256Array, AttributeKindBytesArray, AttributeKindStringArray:
		return true
	}
	return false
}

func (e AttributeKind) String() string {
	return string(e)
}

func (e *AttributeKind) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = AttributeKind(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid AttributeKind", str)
	}
	return nil
}

func (e AttributeKind) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

// RelMatchDirection indicates a direction of the relationship to match.  Edges
// are directional (they have a src node on one end and a dst node on the other)
// Sometimes we want to traverse the graph following this direction, sometimes we
// want to traverse in the oppersite direction, and sometimes it is purely the
// fact that two nodes are connected that we care about.
type RelMatchDirection string

const (
	RelMatchDirectionIn   RelMatchDirection = "IN"
	RelMatchDirectionOut  RelMatchDirection = "OUT"
	RelMatchDirectionBoth RelMatchDirection = "BOTH"
)

var AllRelMatchDirection = []RelMatchDirection{
	RelMatchDirectionIn,
	RelMatchDirectionOut,
	RelMatchDirectionBoth,
}

func (e RelMatchDirection) IsValid() bool {
	switch e {
	case RelMatchDirectionIn, RelMatchDirectionOut, RelMatchDirectionBoth:
		return true
	}
	return false
}

func (e RelMatchDirection) String() string {
	return string(e)
}

func (e *RelMatchDirection) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RelMatchDirection(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RelMatchDirection", str)
	}
	return nil
}

func (e RelMatchDirection) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
