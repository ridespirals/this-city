package sim

// NodeID identifies a junction in the path network.
type NodeID uint32

// EdgeID identifies a directed-geometry edge (travel is bidirectional).
type EdgeID uint32

const (
	NilNode NodeID = 0
	NilEdge EdgeID = 0
)

// Node is a junction where edges meet.
type Node struct {
	ID  NodeID
	Pos Vec2
}

// Edge is a cubic Bézier between two nodes (P0=From, P1=To).
type Edge struct {
	ID    EdgeID
	From  NodeID
	To    NodeID
	Curve CubicBezier
	Poly  Polyline
	Group uint32 // editor chain id; 0 = ungrouped
}

// Network is the navigable street graph.
type Network struct {
	nextNode  NodeID
	nextEdge  EdgeID
	nextGroup uint32
	version   uint64 // bumps on topology / traversal-cost changes
	nodes     map[NodeID]*Node
	edges     map[EdgeID]*Edge
	// incident[node] = edge ids touching the node
	incident map[NodeID][]EdgeID
	// blocked edges are omitted from pathfinding (events, closures).
	blocked map[EdgeID]bool
	// costMul multiplies Poly.Length when > 0; missing or <=0 means 1.
	costMul map[EdgeID]float32
}

func newNetwork() *Network {
	return &Network{
		nextNode:  1,
		nextEdge:  1,
		nextGroup: 1,
		nodes:     make(map[NodeID]*Node),
		edges:     make(map[EdgeID]*Edge),
		incident:  make(map[NodeID][]EdgeID),
		blocked:   make(map[EdgeID]bool),
		costMul:   make(map[EdgeID]float32),
	}
}

// Version increments when connectivity or traversal costs change (for replanners).
func (n *Network) Version() uint64 {
	if n == nil {
		return 0
	}
	return n.version
}

func (n *Network) touch() {
	if n != nil {
		n.version++
	}
}

// SetEdgeBlocked marks an edge non-traversable (or clears the block).
func (n *Network) SetEdgeBlocked(id EdgeID, blocked bool) {
	if n == nil {
		return
	}
	if _, ok := n.edges[id]; !ok {
		return
	}
	if blocked {
		if !n.blocked[id] {
			n.blocked[id] = true
			n.touch()
		}
		return
	}
	if n.blocked[id] {
		delete(n.blocked, id)
		n.touch()
	}
}

// EdgeBlocked reports whether id is closed to pathfinding.
func (n *Network) EdgeBlocked(id EdgeID) bool {
	return n != nil && n.blocked[id]
}

// SetEdgeCostMul sets a travel-cost multiplier for id (1 = default length).
// Values <= 0 clear the override. Very large multipliers simulate congestion.
func (n *Network) SetEdgeCostMul(id EdgeID, mul float32) {
	if n == nil {
		return
	}
	if _, ok := n.edges[id]; !ok {
		return
	}
	if mul <= 0 {
		if _, ok := n.costMul[id]; ok {
			delete(n.costMul, id)
			n.touch()
		}
		return
	}
	if n.costMul[id] != mul {
		n.costMul[id] = mul
		n.touch()
	}
}

// TraversalCost returns the pathfinding weight for an edge, or false if blocked/missing.
func (n *Network) TraversalCost(id EdgeID) (float32, bool) {
	if n == nil {
		return 0, false
	}
	e, ok := n.edges[id]
	if !ok || n.blocked[id] {
		return 0, false
	}
	cost := e.Poly.Length
	if cost < 1e-6 {
		cost = 1e-6
	}
	if m := n.costMul[id]; m > 0 {
		cost *= m
	}
	return cost, true
}

// AddNode inserts a junction at pos.
func (n *Network) AddNode(pos Vec2) NodeID {
	if n == nil {
		return NilNode
	}
	id := n.nextNode
	n.nextNode++
	n.nodes[id] = &Node{ID: id, Pos: pos}
	n.touch()
	return id
}

// EnsureNode returns an existing node within eps of pos, or creates one.
func (n *Network) EnsureNode(pos Vec2, eps float32) NodeID {
	if n == nil {
		return NilNode
	}
	eps2 := eps * eps
	for id, node := range n.nodes {
		d := node.Pos.Sub(pos)
		if d.Dot(d) <= eps2 {
			return id
		}
	}
	return n.AddNode(pos)
}

// AddEdge connects from→to with the given curve (endpoints snapped to node positions).
func (n *Network) AddEdge(from, to NodeID, c0, c1 Vec2) EdgeID {
	return n.AddEdgeGrouped(from, to, c0, c1, 0)
}

// AddEdgeGrouped is AddEdge with an editor group id.
func (n *Network) AddEdgeGrouped(from, to NodeID, c0, c1 Vec2, group uint32) EdgeID {
	if n == nil {
		return NilEdge
	}
	a, okA := n.nodes[from]
	b, okB := n.nodes[to]
	if !okA || !okB {
		return NilEdge
	}
	id := n.nextEdge
	n.nextEdge++
	e := &Edge{
		ID:   id,
		From: from,
		To:   to,
		Curve: CubicBezier{
			P0: a.Pos,
			C0: c0,
			C1: c1,
			P1: b.Pos,
		},
		Group: group,
	}
	e.Poly = BuildPolyline([]CubicBezier{e.Curve}, DefaultPathSamples)
	n.edges[id] = e
	n.incident[from] = append(n.incident[from], id)
	n.incident[to] = append(n.incident[to], id)
	n.touch()
	return id
}

// NewGroup allocates an editor chain group id.
func (n *Network) NewGroup() uint32 {
	if n == nil {
		return 0
	}
	id := n.nextGroup
	n.nextGroup++
	return id
}

// GetEdge returns an edge by id.
func (n *Network) GetEdge(id EdgeID) (*Edge, bool) {
	if n == nil || id == NilEdge {
		return nil, false
	}
	e, ok := n.edges[id]
	return e, ok
}

// GetNode returns a node by id.
func (n *Network) GetNode(id NodeID) (*Node, bool) {
	if n == nil || id == NilNode {
		return nil, false
	}
	node, ok := n.nodes[id]
	return node, ok
}

// ForEachEdge invokes fn for every edge.
func (n *Network) ForEachEdge(fn func(e *Edge)) {
	if n == nil {
		return
	}
	for _, e := range n.edges {
		fn(e)
	}
}

// ForEachNode invokes fn for every node.
func (n *Network) ForEachNode(fn func(node *Node)) {
	if n == nil {
		return
	}
	for _, node := range n.nodes {
		fn(node)
	}
}

// EdgeCount returns the number of edges.
func (n *Network) EdgeCount() int {
	if n == nil {
		return 0
	}
	return len(n.edges)
}

// NodeCount returns the number of nodes.
func (n *Network) NodeCount() int {
	if n == nil {
		return 0
	}
	return len(n.nodes)
}

// IncidentEdges returns edges touching node.
func (n *Network) IncidentEdges(node NodeID) []EdgeID {
	if n == nil {
		return nil
	}
	return append([]EdgeID(nil), n.incident[node]...)
}

// OtherNode returns the endpoint of edge that is not node.
func (n *Network) OtherNode(edge EdgeID, node NodeID) NodeID {
	e, ok := n.GetEdge(edge)
	if !ok {
		return NilNode
	}
	if e.From == node {
		return e.To
	}
	if e.To == node {
		return e.From
	}
	return NilNode
}

// RemoveEdge deletes an edge and cleans adjacency.
func (n *Network) RemoveEdge(id EdgeID) bool {
	e, ok := n.GetEdge(id)
	if !ok {
		return false
	}
	delete(n.edges, id)
	delete(n.blocked, id)
	delete(n.costMul, id)
	n.incident[e.From] = removeEdgeID(n.incident[e.From], id)
	n.incident[e.To] = removeEdgeID(n.incident[e.To], id)
	n.touch()
	return true
}

// RemoveGroup deletes all edges in a group and orphan nodes with no edges.
func (n *Network) RemoveGroup(group uint32) {
	if n == nil || group == 0 {
		return
	}
	var doomed []EdgeID
	for id, e := range n.edges {
		if e.Group == group {
			doomed = append(doomed, id)
		}
	}
	for _, id := range doomed {
		n.RemoveEdge(id)
	}
	n.pruneOrphanNodes()
}

func (n *Network) pruneOrphanNodes() {
	for id := range n.nodes {
		if len(n.incident[id]) == 0 {
			delete(n.nodes, id)
			delete(n.incident, id)
		}
	}
}

func removeEdgeID(list []EdgeID, id EdgeID) []EdgeID {
	out := list[:0]
	for _, e := range list {
		if e != id {
			out = append(out, e)
		}
	}
	return out
}

// SetEdgeCurve updates control points and rebuilds the polyline (keeps endpoints on nodes).
func (n *Network) SetEdgeCurve(id EdgeID, c0, c1 Vec2) bool {
	e, ok := n.GetEdge(id)
	if !ok {
		return false
	}
	a := n.nodes[e.From]
	b := n.nodes[e.To]
	e.Curve = CubicBezier{P0: a.Pos, C0: c0, C1: c1, P1: b.Pos}
	e.Poly = BuildPolyline([]CubicBezier{e.Curve}, DefaultPathSamples)
	n.touch()
	return true
}

// MoveNode sets a node position and refreshes incident edge geometry endpoints.
func (n *Network) MoveNode(id NodeID, pos Vec2) bool {
	node, ok := n.GetNode(id)
	if !ok {
		return false
	}
	node.Pos = pos
	for _, eid := range n.incident[id] {
		e := n.edges[eid]
		a, b := n.nodes[e.From], n.nodes[e.To]
		e.Curve.P0 = a.Pos
		e.Curve.P1 = b.Pos
		e.Poly = BuildPolyline([]CubicBezier{e.Curve}, DefaultPathSamples)
	}
	n.touch()
	return true
}

// AnchorsToChain builds a linear chain of nodes/edges from anchors. Returns group id and edge ids.
func (n *Network) AnchorsToChain(anchors []Vec2, group uint32) (uint32, []EdgeID) {
	if n == nil || len(anchors) < 2 {
		return 0, nil
	}
	if group == 0 {
		group = n.NewGroup()
	} else {
		n.RemoveGroup(group)
	}
	nodes := make([]NodeID, len(anchors))
	for i, p := range anchors {
		nodes[i] = n.AddNode(p)
	}
	var edges []EdgeID
	for i := 0; i < len(anchors)-1; i++ {
		a, b := anchors[i], anchors[i+1]
		d := b.Sub(a)
		eid := n.AddEdgeGrouped(nodes[i], nodes[i+1], a.Add(d.Scale(1.0/3)), a.Add(d.Scale(2.0/3)), group)
		edges = append(edges, eid)
	}
	return group, edges
}

// NearestEdge returns the closest edge sample to pos within maxDist.
func (n *Network) NearestEdge(pos Vec2, maxDist float32) (EdgeID, float32, bool) {
	if n == nil {
		return NilEdge, 0, false
	}
	best := NilEdge
	bestD := maxDist * maxDist
	var bestDistAlong float32
	for _, e := range n.edges {
		for i, pt := range e.Poly.Points {
			d := pt.Sub(pos)
			dd := d.Dot(d)
			if dd <= bestD {
				bestD = dd
				best = e.ID
				if i < len(e.Poly.CumLen) {
					bestDistAlong = e.Poly.CumLen[i]
				}
			}
		}
	}
	if best == NilEdge {
		return NilEdge, 0, false
	}
	return best, bestDistAlong, true
}
