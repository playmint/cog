package model

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/benbjohnson/immutable"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/playmint/ds-node/pkg/contracts/state"
)

type EdgeData struct {
	To     string
	Weight *big.Int
}

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
	edges *immutable.Map[string, *immutable.Map[string, *immutable.Map[uint8, *EdgeData]]]
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
		edges:  immutable.NewMap[string, *immutable.Map[string, *immutable.Map[uint8, *EdgeData]]](nil),
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
	// update the edge data
	edgesForNode, ok := g.edges.Get(srcNodeID)
	if !ok {
		edgesForNode = immutable.NewMap[string, *immutable.Map[uint8, *EdgeData]](nil)
	}
	edgesOfKind, ok := edgesForNode.Get(relID)
	if !ok || edgesOfKind == nil {
		edgesOfKind = immutable.NewMap[uint8, *EdgeData](nil)
	}
	edgesOfKind = edgesOfKind.Set(relKey, &EdgeData{dstNodeID, weight})
	edgesForNode = edgesForNode.Set(relID, edgesOfKind)
	edges := g.edges
	edges = edges.Set(srcNodeID, edgesForNode)

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
	// update the edge data
	edgesForNode, ok := g.edges.Get(srcNodeID)
	if !ok {
		edgesForNode = immutable.NewMap[string, *immutable.Map[uint8, *EdgeData]](nil)
	}
	edgesOfKind, ok := edgesForNode.Get(relID)
	if !ok || edgesOfKind == nil {
		edgesOfKind = immutable.NewMap[uint8, *EdgeData](nil)
	}
	edgesOfKind = edgesOfKind.Delete(relKey)
	edgesForNode = edgesForNode.Set(relID, edgesOfKind)
	edges := g.edges
	edges = edges.Set(srcNodeID, edgesForNode)

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
		if match != nil && len(match.Kinds) > 0 {
			if !contains(match.Kinds, node.Kind()) {
				continue
			}
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
			ID:    annotationRef,
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
	edges, err := n.matchEdges(match)
	if err != nil {
		return nil, err
	}
	for _, edge := range edges {
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
	edges, err := n.matchEdges(match)
	if err != nil {
		return nil, err
	}
	if len(edges) == 0 {
		return nil, nil
	}
	return edges[0], nil
}

func (n *Node) matchEdges(match *Match) ([]*Edge, error) {
	if match == nil {
		match = &Match{}
	}
	edges := []*Edge{}
	// convert kind names to kind ids
	kindIDs := []string{}
	for _, kind := range match.Kinds {
		kindID := n.g.GetKindByName(kind)
		if kindID == "" {
			return nil, fmt.Errorf("no kind found with name %s", kind)
		}
		kindIDs = append(kindIDs, kindID)
	}
	// create a matcher for the nodes
	matchNode := func(nodeID string) bool {
		if len(match.Ids) > 0 && !contains(match.Ids, nodeID) {
			return false
		}
		b, _ := hexutil.Decode(nodeID)
		nodeKind := hexutil.Encode(b[:4])
		for len(kindIDs) > 0 && !contains(kindIDs, nodeKind) {
			return false
		}
		return true
	}
	// create a visitor
	seen := map[*EdgeData]bool{}
	visit := func(edge *Edge) bool {
		if !matchNode(edge.nodeID) {
			return false // don't stop
		}
		edges = append(edges, edge)
		if match.Limit != nil && len(edges) == *match.Limit {
			return true // stop
		}
		return false // don't stop
	}
	err := n.walk(match, visit, seen)
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func (n *Node) walk(match *Match, visit func(*Edge) bool, visited map[*EdgeData]bool) error {
	// walk along all the edges specified by Via
	for _, via := range match.Via {
		dir := RelMatchDirectionOut
		if via.Dir != nil {
			dir = *via.Dir
		}
		relID := n.g.GetRelByName(via.Rel)
		if relID == "" {
			return fmt.Errorf("no rel found with name %s", via.Rel)
		}

		if dir == RelMatchDirectionOut || dir == RelMatchDirectionBoth {
			// find outbound edges from this node
			out, edgesExist := n.g.edges.Get(n.ID)
			if edgesExist {
				outEdgesItr := out.Iterator()
				for !outEdgesItr.Done() {
					outRelID, outEdgesOfRel, outEdgesOfRelExist := outEdgesItr.Next()
					if !outEdgesOfRelExist {
						continue
					}
					if outRelID != relID {
						continue
					}
					outEdgesOfRelItr := outEdgesOfRel.Iterator()
					for !outEdgesOfRelItr.Done() {
						outRelKey, outEdgeData, outEdgeDataExists := outEdgesOfRelItr.Next()
						if !outEdgeDataExists {
							continue
						}
						// if via.Key != outRelKey {
						// 	continue
						// }
						if visited[outEdgeData] {
							continue
						}
						visited[outEdgeData] = true
						edge := &Edge{
							g:      n.g,
							nodeID: outEdgeData.To,
							weight: outEdgeData.Weight,
							key:    outRelKey,
							Dir:    RelMatchDirectionOut,
							RelID:  relID,
						}
						if stop := visit(edge); stop {
							return nil
						}
						if err := edge.Node().walk(match, visit, visited); err != nil {
							return err
						}
					}

				}
			}

		}
		if dir == RelMatchDirectionIn || dir == RelMatchDirectionBoth {
			// find inbound edges to this node
			allEdges := n.g.edges.Iterator()
			for !allEdges.Done() {
				srcNodeID, srcNodeEdges, exists := allEdges.Next()
				if !exists {
					continue
				}
				srcNodeEdgeItr := srcNodeEdges.Iterator()
				for !srcNodeEdgeItr.Done() {
					inRelID, inEdgesOfRel, exists := srcNodeEdgeItr.Next()
					if !exists {
						continue
					}
					if inRelID != relID {
						continue
					}
					inEdgesOfRelItr := inEdgesOfRel.Iterator()
					for !inEdgesOfRelItr.Done() {
						inRelKey, inEdgeData, exists := inEdgesOfRelItr.Next()
						if !exists {
							continue
						}
						if inEdgeData.To != n.ID {
							continue
						}
						if visited[inEdgeData] {
							continue
						}
						visited[inEdgeData] = true
						edge := &Edge{
							g:      n.g,
							nodeID: srcNodeID,
							weight: inEdgeData.Weight,
							key:    inRelKey,
							Dir:    RelMatchDirectionIn,
							RelID:  relID,
						}
						if stop := visit(edge); stop {
							return nil
						}
						if err := edge.Node().walk(match, visit, visited); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func (n *Node) Edges(match *Match) ([]*Edge, error) {
	return n.matchEdges(match)
}

type Edge struct {
	g      *Graph
	nodeID string
	key    uint8
	weight *big.Int
	Dir    RelMatchDirection `json:"dir"`
	RelID  string            `json:"relID"`
}

func (e *Edge) Rel() string {
	relData, ok := e.g.rels.Get(e.RelID)
	if !ok {
		return e.RelID
	}
	return relData.Name
}

func (e *Edge) Node() *Node {
	return e.g.get(e.nodeID)
}

func (e *Edge) Key() int {
	return int(e.key)
}

func (e *Edge) Weight() int {
	if e.weight == nil {
		return 0
	}
	return int(e.weight.Uint64())
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}
