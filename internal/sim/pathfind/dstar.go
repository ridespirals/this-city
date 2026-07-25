package pathfind

import (
	"container/heap"
	"github.com/ridespirals/this-city/internal/sim"
	"math"
)

// pathInf is a stand-in for ∞ in D* Lite (must stay finite for arithmetic).
const pathInf float32 = math.MaxFloat32 / 8

// DStarLite is an incremental replanner for changing graphs and moving starts
// (Koenig & Likhachev). Prefer this over full A* when edges are blocked/unblocked
// often or the agent advances along a route toward a (mostly) fixed goal.
//
// Original “D*” is the older related algorithm; D* Lite is the usual practical form.
// Changing the goal reinitializes the search.
type DStarLite struct {
	net   *sim.Network
	start sim.NodeID
	goal  sim.NodeID
	last  sim.NodeID
	km    float32
	scale float32

	g   map[sim.NodeID]float32
	rhs map[sim.NodeID]float32

	open   dstarHeap
	inOpen map[sim.NodeID]bool

	netVer   uint64
	expanded int
}

type dstarKey struct {
	k1, k2 float32
}

func (a dstarKey) less(b dstarKey) bool {
	if a.k1 != b.k1 {
		return a.k1 < b.k1
	}
	return a.k2 < b.k2
}

type dstarItem struct {
	node sim.NodeID
	key  dstarKey
}

type dstarHeap []dstarItem

func (h dstarHeap) Len() int           { return len(h) }
func (h dstarHeap) Less(i, j int) bool { return h[i].key.less(h[j].key) }
func (h dstarHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *dstarHeap) Push(x any)        { *h = append(*h, x.(dstarItem)) }
func (h *dstarHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// NewDStarLite builds a planner from start → goal on net.
func NewDStarLite(net *sim.Network, start, goal sim.NodeID) *DStarLite {
	d := &DStarLite{
		net:    net,
		start:  start,
		goal:   goal,
		last:   start,
		scale:  1,
		g:      make(map[sim.NodeID]float32),
		rhs:    make(map[sim.NodeID]float32),
		inOpen: make(map[sim.NodeID]bool),
	}
	d.initialize()
	return d
}

func (d *DStarLite) initialize() {
	d.km = 0
	d.g = make(map[sim.NodeID]float32)
	d.rhs = make(map[sim.NodeID]float32)
	d.open = dstarHeap{}
	heap.Init(&d.open)
	d.inOpen = make(map[sim.NodeID]bool)
	d.last = d.start
	if d.net != nil {
		d.netVer = d.net.Version()
	}
	d.setRHS(d.goal, 0)
	d.insert(d.goal, d.calculateKey(d.goal))
}

func (d *DStarLite) getG(s sim.NodeID) float32 {
	if v, ok := d.g[s]; ok {
		return v
	}
	return pathInf
}

func (d *DStarLite) getRHS(s sim.NodeID) float32 {
	if v, ok := d.rhs[s]; ok {
		return v
	}
	return pathInf
}

func (d *DStarLite) setG(s sim.NodeID, v float32)   { d.g[s] = v }
func (d *DStarLite) setRHS(s sim.NodeID, v float32) { d.rhs[s] = v }

func (d *DStarLite) h(a, b sim.NodeID) float32 {
	if d.net == nil {
		return 0
	}
	return heuristic(d.net, a, b, d.scale)
}

func (d *DStarLite) calculateKey(s sim.NodeID) dstarKey {
	m := d.getG(s)
	if r := d.getRHS(s); r < m {
		m = r
	}
	return dstarKey{k1: m + d.h(d.start, s) + d.km, k2: m}
}

func (d *DStarLite) insert(s sim.NodeID, k dstarKey) {
	if d.inOpen[s] {
		for i := range d.open {
			if d.open[i].node == s {
				d.open[i].key = k
				heap.Fix(&d.open, i)
				return
			}
		}
	}
	heap.Push(&d.open, dstarItem{node: s, key: k})
	d.inOpen[s] = true
}

func (d *DStarLite) remove(s sim.NodeID) {
	if !d.inOpen[s] {
		return
	}
	for i := range d.open {
		if d.open[i].node == s {
			heap.Remove(&d.open, i)
			delete(d.inOpen, s)
			return
		}
	}
	delete(d.inOpen, s)
}

func (d *DStarLite) updateVertex(u sim.NodeID) {
	if d.net == nil {
		return
	}
	if u != d.goal {
		minRHS := pathInf
		for _, link := range linksFrom(d.net, u) {
			v := link.cost + d.getG(link.to)
			if v < minRHS {
				minRHS = v
			}
		}
		d.setRHS(u, minRHS)
	}
	if d.inOpen[u] {
		d.remove(u)
	}
	if d.getG(u) != d.getRHS(u) {
		d.insert(u, d.calculateKey(u))
	}
}

func (d *DStarLite) topKey() (dstarKey, bool) {
	if d.open.Len() == 0 {
		return dstarKey{}, false
	}
	return d.open[0].key, true
}

func (d *DStarLite) computeShortestPath() {
	if d.net == nil {
		return
	}
	const maxIter = 1_000_000
	for iter := 0; iter < maxIter; iter++ {
		tk, ok := d.topKey()
		sk := d.calculateKey(d.start)
		if !ok || (!tk.less(sk) && d.getRHS(d.start) == d.getG(d.start)) {
			return
		}
		u := heap.Pop(&d.open).(dstarItem)
		delete(d.inOpen, u.node)
		d.expanded++

		kNew := d.calculateKey(u.node)
		if u.key.less(kNew) {
			d.insert(u.node, kNew)
			continue
		}
		if d.getG(u.node) > d.getRHS(u.node) {
			d.setG(u.node, d.getRHS(u.node))
			for _, link := range linksFrom(d.net, u.node) {
				d.updateVertex(link.to)
			}
			continue
		}
		d.setG(u.node, pathInf)
		d.updateVertex(u.node)
		for _, link := range linksFrom(d.net, u.node) {
			d.updateVertex(link.to)
		}
	}
}

// SetStart moves the agent’s current node (cheap incremental update via km).
func (d *DStarLite) SetStart(start sim.NodeID) {
	if d == nil || start == sim.NilNode || start == d.start {
		return
	}
	d.km += d.h(d.last, start)
	d.last = start
	d.start = start
}

// SetGoal changes the destination and reinitializes the search.
func (d *DStarLite) SetGoal(goal sim.NodeID) {
	if d == nil || goal == sim.NilNode {
		return
	}
	if goal == d.goal {
		return
	}
	d.goal = goal
	d.initialize()
}

// SyncNetwork applies pending sim.Network.Version changes (blocks, cost muls, edits).
func (d *DStarLite) SyncNetwork() {
	if d == nil || d.net == nil {
		return
	}
	ver := d.net.Version()
	if ver == d.netVer {
		return
	}
	d.netVer = ver
	d.km += d.h(d.last, d.start)
	d.last = d.start
	seen := map[sim.NodeID]bool{}
	touch := func(id sim.NodeID) {
		if id == sim.NilNode || seen[id] {
			return
		}
		seen[id] = true
		d.updateVertex(id)
	}
	for id := range d.g {
		touch(id)
	}
	for id := range d.rhs {
		touch(id)
	}
	d.net.ForEachEdge(func(e *sim.Edge) {
		touch(e.From)
		touch(e.To)
	})
}

// Replan syncs network changes, computes, and extracts a path start→goal.
func (d *DStarLite) Replan() Result {
	if d == nil || d.net == nil {
		return Result{}
	}
	d.SyncNetwork()
	if d.start == d.goal {
		return Result{Found: true, Nodes: []sim.NodeID{d.start}, Cost: 0, Expanded: d.expanded}
	}
	before := d.expanded
	d.computeShortestPath()
	return d.extractPath(d.expanded - before)
}

// Path is an alias for Replan.
func (d *DStarLite) Path() Result { return d.Replan() }

func (d *DStarLite) extractPath(expanded int) Result {
	if d.getG(d.start) >= pathInf/2 {
		return Result{Expanded: expanded}
	}
	nodes := []sim.NodeID{d.start}
	var edges []sim.EdgeID
	cur := d.start
	seen := map[sim.NodeID]bool{cur: true}
	limit := d.net.NodeCount() + 2
	for cur != d.goal {
		bestTo := sim.NilNode
		bestEdge := sim.NilEdge
		best := pathInf
		var bestCost float32
		for _, link := range linksFrom(d.net, cur) {
			v := link.cost + d.getG(link.to)
			if v < best {
				best = v
				bestTo = link.to
				bestEdge = link.via
				bestCost = link.cost
			}
		}
		_ = bestCost
		if bestTo == sim.NilNode || seen[bestTo] {
			return Result{Expanded: expanded}
		}
		nodes = append(nodes, bestTo)
		edges = append(edges, bestEdge)
		seen[bestTo] = true
		cur = bestTo
		if len(nodes) > limit {
			return Result{Expanded: expanded}
		}
	}
	return Result{
		Found:    true,
		Nodes:    nodes,
		Edges:    edges,
		Cost:     pathCost(d.net, edges),
		Expanded: expanded,
	}
}
