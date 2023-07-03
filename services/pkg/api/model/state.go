package model

import (
	"encoding/json"
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

type SerializableGraph struct {
	// cache of which Nodes we have seen edges between
	Nodes map[string]bool
	// nodeID => relID => relKey => EdgeData
	Edges map[string]*DirectedEdge
	// mapping of relID => friendly name
	Rels map[string]*state.StateEdgeTypeRegister
	// mapping of kindID => friendly name
	Kinds map[string]*state.StateNodeTypeRegister
	// nodeID => label => annotationRef
	Labels map[string]map[string]string
	// annotationRef => data
	Ann map[string]string
	// Block is the last seen update to the graph
	Block uint64
}

type Graph struct {
	// cache of which Nodes we have seen edges between
	Nodes *immutable.Map[string, bool]
	// nodeID => relID => relKey => EdgeData
	Edges *immutable.Map[string, *DirectedEdge]
	// mapping of relID => friendly name
	Rels *immutable.Map[string, *state.StateEdgeTypeRegister]
	// mapping of kindID => friendly name
	Kinds *immutable.Map[string, *state.StateNodeTypeRegister]
	// nodeID => label => annotationRef
	Labels *immutable.Map[string, *immutable.Map[string, string]]
	// annotationRef => data
	Ann *immutable.Map[string, string]
	// Block is the last seen update to the graph
	Block uint64

	directedEdgeCache map[string][]*Edge
}

func NewGraph(block uint64) *Graph {
	return &Graph{
		Nodes:             immutable.NewMap[string, bool](nil),
		Edges:             immutable.NewMap[string, *DirectedEdge](nil),
		Rels:              immutable.NewMap[string, *state.StateEdgeTypeRegister](nil),
		Kinds:             immutable.NewMap[string, *state.StateNodeTypeRegister](nil),
		Labels:            immutable.NewMap[string, *immutable.Map[string, string]](nil),
		Ann:               immutable.NewMap[string, string](nil),
		Block:             block,
		directedEdgeCache: map[string][]*Edge{},
	}
}

func (g *Graph) Dump() (string, error) {
	sg := &SerializableGraph{
		Nodes:  map[string]bool{},
		Edges:  map[string]*DirectedEdge{},
		Rels:   map[string]*state.StateEdgeTypeRegister{},
		Kinds:  map[string]*state.StateNodeTypeRegister{},
		Labels: map[string]map[string]string{},
		Ann:    map[string]string{},
		Block:  g.Block,
	}

	{
		itr := g.Nodes.Iterator()
		for !itr.Done() {
			k, v, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Nodes[k] = v
		}
	}

	{
		itr := g.Edges.Iterator()
		for !itr.Done() {
			k, v, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Edges[k] = v
		}
	}

	{
		itr := g.Rels.Iterator()
		for !itr.Done() {
			k, v, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Rels[k] = v
		}
	}

	{
		itr := g.Kinds.Iterator()
		for !itr.Done() {
			k, v, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Kinds[k] = v
		}
	}

	{
		itr := g.Labels.Iterator()
		for !itr.Done() {
			k1, v1, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Labels[k1] = map[string]string{}
			itr2 := v1.Iterator()
			for !itr2.Done() {
				k2, v2, ok := itr2.Next()
				if !ok {
					continue
				}
				sg.Labels[k1][k2] = v2
			}
		}
	}

	{
		itr := g.Ann.Iterator()
		for !itr.Done() {
			k, v, ok := itr.Next()
			if !ok {
				continue
			}
			sg.Ann[k] = v
		}
	}

	b, err := json.Marshal(sg)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func LoadGraph(sg *SerializableGraph) (*Graph, error) {
	g := NewGraph(sg.Block)
	for k, v := range sg.Nodes {
		g.Nodes = g.Nodes.Set(k, v)
	}
	for k, v := range sg.Edges {
		g.Edges = g.Edges.Set(k, v)
	}
	for k, v := range sg.Rels {
		g.Rels = g.Rels.Set(k, v)
	}
	for k, v := range sg.Kinds {
		g.Kinds = g.Kinds.Set(k, v)
	}
	for k1, v1 := range sg.Labels {
		label := immutable.NewMap[string, string](nil)
		for k2, v2 := range v1 {
			label = label.Set(k2, v2)
		}
		g.Labels = g.Labels.Set(k1, label)

	}
	for k, v := range sg.Ann {
		g.Ann = g.Ann.Set(k, v)
	}
	g.updateDirectEdgeCache()
	return g, nil
}

func (g *Graph) BlockNumber() uint64 {
	return g.Block
}

func (g *Graph) SetRelData(relData *state.StateEdgeTypeRegister) *Graph {
	return &Graph{
		Nodes:             g.Nodes,
		Edges:             g.Edges,
		Rels:              g.Rels.Set(hexutil.Encode(relData.Id[:]), relData),
		Kinds:             g.Kinds,
		Labels:            g.Labels,
		Ann:               g.Ann,
		Block:             g.Block,
		directedEdgeCache: g.directedEdgeCache,
	}
}

func (g *Graph) SetKindData(kindData *state.StateNodeTypeRegister) *Graph {
	return &Graph{
		Nodes:             g.Nodes,
		Edges:             g.Edges,
		Rels:              g.Rels,
		Kinds:             g.Kinds.Set(hexutil.Encode(kindData.Id[:]), kindData),
		Labels:            g.Labels,
		Ann:               g.Ann,
		Block:             g.Block,
		directedEdgeCache: g.directedEdgeCache,
	}
}

func (g *Graph) SetAnnotationData(nodeID string, label string, ref string, data string, block uint64) *Graph {
	// update ann data
	ann := g.Ann.Set(ref, data)
	// update the label data
	labels, ok := g.Labels.Get(nodeID)
	if !ok {
		labels = immutable.NewMap[string, string](nil)
	}
	labels = labels.Set(label, ref)

	// update the node data to mark the seen nodes
	nodes := g.Nodes
	nodes = nodes.Set(nodeID, true)

	// build our new graph
	newGraph := &Graph{
		Nodes:             nodes,
		Edges:             g.Edges,
		Rels:              g.Rels,
		Kinds:             g.Kinds,
		Labels:            g.Labels.Set(nodeID, labels),
		Ann:               ann,
		Block:             block,
		directedEdgeCache: g.directedEdgeCache,
	}

	return newGraph
}

func (g *Graph) SetEdge(relID string, relKey uint8, srcNodeID string, dstNodeID string, weight *big.Int, block uint64) *Graph {

	e := &DirectedEdge{
		From:   srcNodeID,
		To:     dstNodeID,
		Weight: weight,
		Key:    relKey,
		Rel:    relID,
	}

	// set the edge going in both directions
	edges := g.Edges.Set(e.ID(), e)

	// update the node data to mark the seen nodes
	nodes := g.Nodes
	nodes = nodes.Set(srcNodeID, true)
	nodes = nodes.Set(dstNodeID, true)

	// build our new graph
	newGraph := &Graph{
		Nodes:             nodes,
		Edges:             edges,
		Rels:              g.Rels,
		Kinds:             g.Kinds,
		Labels:            g.Labels,
		Ann:               g.Ann,
		Block:             block,
		directedEdgeCache: map[string][]*Edge{},
	}

	newGraph.updateDirectEdgeCache()

	return newGraph
}

func (g *Graph) RemoveEdge(relID string, relKey uint8, srcNodeID string, block uint64) *Graph {

	// remove edge
	e := &DirectedEdge{
		From: srcNodeID,
		Rel:  relID,
		Key:  relKey,
	}
	edges := g.Edges.Delete(e.ID())

	// TODO: should we remove nodes from the node list that have no edges connected?

	// build our new graph
	newGraph := &Graph{
		Nodes:             g.Nodes,
		Edges:             edges,
		Rels:              g.Rels,
		Kinds:             g.Kinds,
		Labels:            g.Labels,
		Ann:               g.Ann,
		Block:             block,
		directedEdgeCache: map[string][]*Edge{},
	}

	newGraph.updateDirectEdgeCache()

	return newGraph
}

func (g *Graph) GetRelByName(relName string) string {
	itr := g.Rels.Iterator()
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
	itr := g.Kinds.Iterator()
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
	_, exists := g.Nodes.Get(id)
	return exists
}

func (g *Graph) GetNodes(match *Match) []*Node {
	nodes := []*Node{}
	itr := g.Nodes.Iterator()
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

func (g *Graph) updateDirectEdgeCache() {
	edgesItr := g.Edges.Iterator()
	for !edgesItr.Done() {
		_, e, exists := edgesItr.Next()
		if !exists {
			continue
		}

		g.directedEdgeCache[e.From] = append(g.directedEdgeCache[e.From], &Edge{
			g:            g,
			DirectedEdge: e,
			Dir:          RelMatchDirectionOut,
		})
		g.directedEdgeCache[e.To] = append(g.directedEdgeCache[e.To], &Edge{
			g:            g,
			DirectedEdge: e,
			Dir:          RelMatchDirectionIn,
		})
	}
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
	kindData, ok := n.g.Kinds.Get(kindID)
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
	labels, ok := n.g.Labels.Get(n.ID)
	if !ok {
		return annotations
	}
	itr := labels.Iterator()
	for !itr.Done() {
		label, annotationRef, ok := itr.Next()
		if !ok {
			continue
		}
		data, ok := n.g.Ann.Get(annotationRef)
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
	kindData, ok := n.g.Kinds.Get(hexutil.Encode(kindID[:4]))
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
		sum += edge.WeightInt()
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
	w := edge.WeightInt()
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

	for _, edge := range n.g.directedEdgeCache[n.ID] {
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
		return e.g.get(e.To)
	} else {
		return e.g.get(e.From)
	}
}

type DirectedEdge struct {
	From   string
	To     string
	Key    uint8
	Weight *big.Int
	Rel    string
}

func (e *Edge) RelString() string {
	relData, ok := e.g.Rels.Get(e.Rel)
	if !ok {
		return e.Rel
	}
	return relData.Name
}

func (e *Edge) Src() *Node {
	return e.g.get(e.From)
}

func (e *Edge) Dst() *Node {
	return e.g.get(e.To)
}

func (e *DirectedEdge) ID() string {
	return fmt.Sprintf("%s-%s-%d", e.From, e.Rel, e.Key)
}

func (e *DirectedEdge) KeyInt() int {
	return int(e.Key)
}

func (e *DirectedEdge) WeightInt() int {
	if e.Weight == nil {
		return 0
	}
	return int(e.Weight.Uint64())
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
			if matchRelID == e.Rel &&
				(via.Key == nil || *(via.Key) == e.KeyInt()) &&
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
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
