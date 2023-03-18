package model

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/benbjohnson/immutable"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/contracts/state"
)

type CompoundKeyKind uint8

const (
	CompoundKeyKindNone CompoundKeyKind = iota
	CompoundKeyKindUint160
	CompoundKeyKindUint8Array
	CompoundKeyKindInt8Array
	CompoundKeyKindUint16Array
	CompoundKeyKindInt16Array
	CompoundKeyKindUint32Array
	CompoundKeyKindInt32Array
	CompoundKeyKindUint64Array
	CompoundKeyKindInt64Array
	CompoundKeyKindAddress
	CompoundKeyKindBytes
	CompoundKeyKindString
)

type Graph struct {
	// cache of which nodes we have seen edges between
	nodes *immutable.Map[string, bool]
	// nodeID => relID => relKey => EdgeData
	edges *immutable.Map[string, *DirectedEdge]
	// mapping of relID => friendly name
	rels *immutable.Map[string, *state.StateEdgeTypeRegister]
	// mapping of kindID => friendly name
	kinds *immutable.Map[string, *state.StateNodeTypeRegister]
	// nodeID => label => annotationRef
	labels *immutable.Map[string, *immutable.Map[string, string]]
	// annotationRef => data
	ann *immutable.Map[string, string]
	// block is the last seen update to the graph
	block uint64
}

func NewGraph(block uint64) *Graph {
	return &Graph{
		nodes:  immutable.NewMap[string, bool](nil),
		edges:  immutable.NewMap[string, *DirectedEdge](nil),
		rels:   immutable.NewMap[string, *state.StateEdgeTypeRegister](nil),
		kinds:  immutable.NewMap[string, *state.StateNodeTypeRegister](nil),
		labels: immutable.NewMap[string, *immutable.Map[string, string]](nil),
		ann:    immutable.NewMap[string, string](nil),
		block:  block,
	}
}

func (g *Graph) BlockNumber() uint64 {
	return g.block
}

func (g *Graph) SetRelData(relData *state.StateEdgeTypeRegister) *Graph {
	return &Graph{
		nodes:  g.nodes,
		edges:  g.edges,
		rels:   g.rels.Set(hexutil.Encode(relData.Id[:]), relData),
		kinds:  g.kinds,
		labels: g.labels,
		ann:    g.ann,
		block:  g.block,
	}
}

func (g *Graph) SetKindData(kindData *state.StateNodeTypeRegister) *Graph {
	return &Graph{
		nodes:  g.nodes,
		edges:  g.edges,
		rels:   g.rels,
		kinds:  g.kinds.Set(hexutil.Encode(kindData.Id[:]), kindData),
		labels: g.labels,
		ann:    g.ann,
		block:  g.block,
	}
}

func (g *Graph) SetAnnotationData(nodeID string, label string, ref string, data string, block uint64) *Graph {
	// update ann data
	ann := g.ann.Set(ref, data)
	// update the label data
	labels, ok := g.labels.Get(nodeID)
	if !ok {
		labels = immutable.NewMap[string, string](nil)
	}
	labels = labels.Set(label, ref)

	// update the node data to mark the seen nodes
	nodes := g.nodes
	nodes = nodes.Set(nodeID, true)

	// build our new graph
	newGraph := &Graph{
		nodes:  nodes,
		edges:  g.edges,
		rels:   g.rels,
		kinds:  g.kinds,
		labels: g.labels.Set(nodeID, labels),
		ann:    ann,
		block:  block,
	}

	return newGraph
}

func (g *Graph) SetEdge(relID string, relKey uint8, srcNodeID string, dstNodeID string, weight *big.Int, block uint64) *Graph {

	e := &DirectedEdge{
		from:   srcNodeID,
		to:     dstNodeID,
		weight: weight,
		key:    relKey,
		rel:    relID,
	}

	// set the edge going in both directions
	edges := g.edges.Set(e.ID(), e)

	// update the node data to mark the seen nodes
	nodes := g.nodes
	nodes = nodes.Set(srcNodeID, true)
	nodes = nodes.Set(dstNodeID, true)

	// build our new graph
	newGraph := &Graph{
		nodes:  nodes,
		edges:  edges,
		rels:   g.rels,
		kinds:  g.kinds,
		labels: g.labels,
		ann:    g.ann,
		block:  block,
	}

	return newGraph
}

func (g *Graph) RemoveEdge(relID string, relKey uint8, srcNodeID string, block uint64) *Graph {

	// remove edge
	e := &DirectedEdge{
		from: srcNodeID,
		rel:  relID,
		key:  relKey,
	}
	edges := g.edges.Delete(e.ID())

	// TODO: should we remove nodes from the node list that have no edges connected?

	// build our new graph
	newGraph := &Graph{
		nodes:  g.nodes,
		edges:  edges,
		rels:   g.rels,
		kinds:  g.kinds,
		labels: g.labels,
		ann:    g.ann,
		block:  block,
	}

	return newGraph
}

func (g *Graph) GetRelByName(relName string) string {
	itr := g.rels.Iterator()
	for !itr.Done() {
		relID, relData, ok := itr.Next()
		if !ok {
			continue
		}
		if relData.Name == relName {
			return relID
		}
	}
	return ""
}

func (g *Graph) GetKindByName(kindName string) string {
	itr := g.kinds.Iterator()
	for !itr.Done() {
		kindID, kindData, ok := itr.Next()
		if !ok {
			continue
		}
		if kindData.Name == kindName {
			return kindID
		}
	}
	return ""
}

func (g *Graph) get(id string) *Node {
	return &Node{
		g:  g,
		ID: id,
	}
}

func (g *Graph) GetNode(match *Match) *Node {
	if match == nil {
		match = &Match{}
	}
	limit := 1
	match.Limit = &limit
	nodes := g.GetNodes(match)
	if len(nodes) == 0 {
		return nil
	}
	return nodes[0]
}

func (g *Graph) NodeExists(id string) bool {
	_, exists := g.nodes.Get(id)
	return exists
}

func (g *Graph) GetNodes(match *Match) []*Node {
	nodes := []*Node{}
	itr := g.nodes.Iterator()
	for !itr.Done() {
		id, _, exists := itr.Next()
		if !exists {
			continue
		}
		node := g.get(id)
		if match != nil && !match.MatchNode(node) {
			continue
		}
		nodes = append(nodes, node)
	}
	return nodes
}

type Node struct {
	g *Graph

	ID   string `json:"id"`
	kind string
}

func (n *Node) Keys() ([]*big.Int, error) {
	id, err := hexutil.Decode(n.ID)
	if err != nil {
		return nil, fmt.Errorf("keys: failed to decode node id %v: %v", n.ID, err)
	}
	// find the compound key type
	kindID := hexutil.Encode(id[:4])
	kindData, ok := n.g.kinds.Get(kindID)
	if !ok {
		return nil, fmt.Errorf("keys: no kind type data for kind id %v", kindID)
	}
	switch CompoundKeyKind(kindData.KeyKind) {
	case CompoundKeyKindUint64Array, CompoundKeyKindInt64Array:
		return n.splitKeys(id, 1, 8)
	case CompoundKeyKindUint32Array, CompoundKeyKindInt32Array:
		return n.splitKeys(id, 2, 4)
	case CompoundKeyKindUint16Array, CompoundKeyKindInt16Array:
		return n.splitKeys(id, 4, 2)
	case CompoundKeyKindUint8Array, CompoundKeyKindInt8Array:
		return n.splitKeys(id, 8, 1)
	default:
		return n.splitKeys(id, 1, 20)
	}
}

func (n *Node) splitKeys(id []byte, split int, nbytes int) ([]*big.Int, error) {
	key := id[4:]
	keys := []*big.Int{}
	for len(key) > 0 {
		keys = append([]*big.Int{
			big.NewInt(0).SetBytes(key[len(key)-nbytes:]),
		}, keys...)
		key = key[:len(key)-nbytes]
	}
	return keys[len(keys)-split:], nil
}

func (n *Node) Key() (*big.Int, error) {
	id, err := hexutil.Decode(n.ID)
	if err != nil {
		return nil, fmt.Errorf("keys: failed to decode node id %v: %v", n.ID, err)
	}
	return big.NewInt(0).SetBytes(id[4:]), nil
}

func (n *Node) Annotation(name string) *Annotation {
	for _, ann := range n.Annotations() {
		if ann.Name == name {
			return ann
		}
	}
	return nil
}

func (n *Node) Annotations() []*Annotation {
	annotations := []*Annotation{}
	labels, ok := n.g.labels.Get(n.ID)
	if !ok {
		return annotations
	}
	itr := labels.Iterator()
	for !itr.Done() {
		label, annotationRef, ok := itr.Next()
		if !ok {
			continue
		}
		data, ok := n.g.ann.Get(annotationRef)
		if !ok {
			continue
		}
		annotations = append(annotations, &Annotation{
			ID:    fmt.Sprintf("%s-%s", n.ID, label),
			Ref:   annotationRef,
			Name:  label,
			Value: data,
		})

	}
	return annotations
}

func (n *Node) Kind() string {
	if n.kind != "" {
		return n.kind
	}
	kindID, err := hexutil.Decode(n.ID)
	if err != nil {
		// FIXME: stop panicing here
		panic(fmt.Sprintf("failed to decode node id: %v", err))
	}
	kindData, ok := n.g.kinds.Get(hexutil.Encode(kindID[:4]))
	if !ok {
		return n.ID
	}
	n.kind = kindData.Name
	return n.kind
}

func (n *Node) Count(match *Match) (int, error) {
	edges, err := n.Edges(match)
	if err != nil {
		return 0, err
	}
	return len(edges), nil
}

func (n *Node) Sum(match *Match) (int, error) {
	edges, err := n.Edges(match)
	if err != nil {
		return 0, err
	}
	sum := 0
	for _, edge := range edges {
		sum += edge.Weight()
	}
	return sum, nil
}

func (n *Node) Value(match *Match) (*int, error) {
	edge, err := n.Edge(match)
	if err != nil {
		return nil, err
	}
	if edge == nil {
		return nil, nil
	}
	w := edge.Weight()
	return &w, nil
}

func (n *Node) Node(match *Match) (*Node, error) {
	nodes, err := n.Nodes(match)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, nil
	}
	return nodes[0], nil
}

// defaults to Out edges
func (n *Node) Nodes(match *Match) ([]*Node, error) {
	seen := map[string]bool{}
	nodes := []*Node{}
	for _, edge := range n.matchEdges(match, 0, map[string]bool{}) {
		node := edge.Node()
		if seen[node.ID] {
			continue
		}
		seen[node.ID] = true
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (n *Node) Edge(match *Match) (*Edge, error) {
	if match == nil {
		match = &Match{}
	}
	limit := 1
	match.Limit = &limit
	edges := n.matchEdges(match, 0, map[string]bool{})
	if len(edges) == 0 {
		return nil, nil
	}
	return edges[0], nil
}

func (n *Node) Edges(match *Match) ([]*Edge, error) {
	return n.matchEdges(match, 0, map[string]bool{}), nil
}

func (n *Node) getDirectEdges(match *Match) []*Edge {
	result := []*Edge{}

	edgesItr := n.g.edges.Iterator()
	for !edgesItr.Done() {
		_, directedEdge, exists := edgesItr.Next()
		if !exists {
			continue
		}
		edge := &Edge{
			g:            n.g,
			DirectedEdge: directedEdge,
		}
		if edge.from == n.ID {
			// match normal
			edge.Dir = RelMatchDirectionOut
		} else if edge.to == n.ID {
			// match reverse
			edge.Dir = RelMatchDirectionIn
		} else {
			// no direct match
			continue
		}
		if !match.MatchEdge(edge) {
			continue
		}
		result = append(result, edge)
	}
	return result
}

func (n *Node) matchEdges(match *Match, depth int, seen map[string]bool) []*Edge {
	// add ourself to the seen list
	seen[n.ID] = true
	// get all the connections from this node
	edges := n.getDirectEdges(match)
	// for each edge collect the edges of their node
	// do this until we hit match.MaxDepth
	if match != nil && match.MaxDepth != nil && depth < *(match.MaxDepth) {
		for _, e := range edges {
			next := e.Node()
			if seen[next.ID] {
				continue
			}
			seen[next.ID] = true
			for _, moreEdges := range next.matchEdges(match, depth+1, seen) {
				edges = append(edges, moreEdges)
			}
		}
	}

	// return what we got
	return edges
}

type Edge struct {
	g   *Graph
	Dir RelMatchDirection
	*DirectedEdge
}

func (e *Edge) Node() *Node {
	if e.Dir == RelMatchDirectionOut {
		return e.g.get(e.to)
	} else {
		return e.g.get(e.from)
	}
}

type DirectedEdge struct {
	from   string
	to     string
	key    uint8
	weight *big.Int
	rel    string
}

func (e *Edge) Rel() string {
	relData, ok := e.g.rels.Get(e.rel)
	if !ok {
		return e.rel
	}
	return relData.Name
}

func (e *Edge) Src() *Node {
	return e.g.get(e.from)
}

func (e *Edge) Dst() *Node {
	return e.g.get(e.to)
}

func (e *DirectedEdge) ID() string {
	return fmt.Sprintf("%s-%s-%d", e.from, e.rel, e.key)
}

func (e *DirectedEdge) Key() int {
	return int(e.key)
}

func (e *DirectedEdge) Weight() int {
	if e.weight == nil {
		return 0
	}
	return int(e.weight.Uint64())
}

func (match *Match) MatchNode(n *Node) bool {
	if match == nil {
		return true
	}
	if len(match.Ids) > 0 && !contains(match.Ids, n.ID) {
		return false
	}
	if len(match.Kinds) > 0 && !contains(match.Kinds, n.Kind()) {
		return false
	}
	return true
}

func (match *Match) MatchEdge(e *Edge) bool {
	if ok := match.MatchVia(e); !ok {
		return false
	}
	if ok := match.MatchNode(e.Node()); !ok {
		return false
	}
	return true
}

func (match *Match) MatchVia(e *Edge) bool {
	if match == nil {
		return true
	}
	if !match.viaRels(e) {
		return false
	}
	return true
}

func (match *Match) viaRels(e *Edge) bool {
	if match.Via == nil {
		return true
	}
	numRels := 0
	for _, via := range match.Via {
		if via.Rel != "" {
			matchRelID := e.g.GetRelByName(via.Rel)
			if matchRelID == "" {
				continue // invalid
			}
			dir := RelMatchDirectionOut
			if via.Dir != nil {
				dir = *(via.Dir)
			}
			numRels++
			if matchRelID == e.rel &&
				(via.Key == nil || *(via.Key) == e.Key()) &&
				(dir == RelMatchDirectionBoth || dir == e.Dir) {
				return true
			}
		}
	}
	if numRels == 0 {
		return true
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		fmt.Println("a vs e", a, e)
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
