package pathfind

import (
	"container/heap"
	"github.com/ridespirals/this-city/internal/sim"
	"math"
)

// Algo selects a graph search strategy on the street sim.Network.
type Algo int

const (
	// AlgoAStar is best-first search with Euclidean heuristic (optimal if h admissible).
	AlgoAStar Algo = iota
	// AlgoDijkstra is uniform-cost search (A* with h=0); shortest weighted path.
	AlgoDijkstra
	// AlgoBFS is breadth-first search; shortest path in hop count (ignores lengths).
	AlgoBFS
	// AlgoDFS is depth-first search; finds some path, not necessarily short.
	AlgoDFS
	// AlgoBidirectionalBFS grows frontiers from start and goal until they meet.
	// Also called bidirectional search / "meet in the middle".
	AlgoBidirectionalBFS
	// AlgoBidirectionalAStar is bidirectional A* (meet-in-the-middle with heuristics).
	AlgoBidirectionalAStar
	// AlgoBidirectionalDijkstra is bidirectional Dijkstra (both sides uniform-cost).
	AlgoBidirectionalDijkstra
)

// Query requests a route between two network nodes.
type Query struct {
	From sim.NodeID
	To   sim.NodeID
	Algo Algo
	// HeuristicScale multiplies the A* heuristic (1 = default). Use 0 to force Dijkstra-like.
	HeuristicScale float32
}

// Result is a found route (or Found=false).
type Result struct {
	Found    bool
	Nodes    []sim.NodeID
	Edges    []sim.EdgeID // for sim.PathDecision.Route
	Cost     float32      // sum of edge lengths (BFS/DFS: hop-based or length along found path)
	Expanded int          // nodes popped / expanded (debug / tests)
}

// Find runs the selected algorithm on net.
func Find(net *sim.Network, q Query) Result {
	if net == nil || q.From == sim.NilNode || q.To == sim.NilNode {
		return Result{}
	}
	if _, ok := net.GetNode(q.From); !ok {
		return Result{}
	}
	if _, ok := net.GetNode(q.To); !ok {
		return Result{}
	}
	if q.From == q.To {
		return Result{Found: true, Nodes: []sim.NodeID{q.From}, Cost: 0}
	}
	if q.HeuristicScale == 0 && (q.Algo == AlgoAStar || q.Algo == AlgoBidirectionalAStar) {
		q.HeuristicScale = 1
	}
	switch q.Algo {
	case AlgoDijkstra:
		return searchAStar(net, q, true)
	case AlgoBFS:
		return searchBFS(net, q)
	case AlgoDFS:
		return searchDFS(net, q)
	case AlgoBidirectionalBFS:
		return searchBidirectionalBFS(net, q)
	case AlgoBidirectionalDijkstra:
		return searchBidirectionalAStar(net, q, true)
	case AlgoBidirectionalAStar:
		return searchBidirectionalAStar(net, q, false)
	default: // AlgoAStar
		return searchAStar(net, q, false)
	}
}

// AStar is Find with AlgoAStar.
func AStar(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoAStar, HeuristicScale: 1})
}

// Dijkstra is Find with AlgoDijkstra.
func Dijkstra(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoDijkstra})
}

// BFS is Find with AlgoBFS.
func BFS(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoBFS})
}

// DFS is Find with AlgoDFS.
func DFS(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoDFS})
}

// BidirectionalBFS is Find with AlgoBidirectionalBFS (meet-in-the-middle).
func BidirectionalBFS(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoBidirectionalBFS})
}

// BidirectionalAStar is Find with AlgoBidirectionalAStar.
func BidirectionalAStar(net *sim.Network, from, to sim.NodeID) Result {
	return Find(net, Query{From: from, To: to, Algo: AlgoBidirectionalAStar, HeuristicScale: 1})
}

func reconstruct(cameNode map[sim.NodeID]sim.NodeID, cameEdge map[sim.NodeID]sim.EdgeID, start, goal sim.NodeID) (nodes []sim.NodeID, edges []sim.EdgeID) {
	cur := goal
	for cur != start {
		nodes = append(nodes, cur)
		e := cameEdge[cur]
		edges = append(edges, e)
		prev, ok := cameNode[cur]
		if !ok {
			return nil, nil
		}
		cur = prev
	}
	nodes = append(nodes, start)
	// reverse
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
	for i, j := 0, len(edges)-1; i < j; i, j = i+1, j-1 {
		edges[i], edges[j] = edges[j], edges[i]
	}
	return nodes, edges
}

// --- A* / Dijkstra ---

func searchAStar(net *sim.Network, q Query, dijkstra bool) Result {
	scale := q.HeuristicScale
	if dijkstra {
		scale = 0
	}
	open := &nodeHeap{}
	heap.Init(open)
	gScore := map[sim.NodeID]float32{q.From: 0}
	cameNode := map[sim.NodeID]sim.NodeID{}
	cameEdge := map[sim.NodeID]sim.EdgeID{}
	inOpen := map[sim.NodeID]bool{q.From: true}
	heap.Push(open, heapItem{node: q.From, f: heuristic(net, q.From, q.To, scale)})

	expanded := 0
	for open.Len() > 0 {
		cur := heap.Pop(open).(heapItem)
		inOpen[cur.node] = false
		expanded++
		if cur.node == q.To {
			nodes, edges := reconstruct(cameNode, cameEdge, q.From, q.To)
			return Result{Found: true, Nodes: nodes, Edges: edges, Cost: gScore[q.To], Expanded: expanded}
		}
		for _, link := range linksFrom(net, cur.node) {
			tentative := gScore[cur.node] + link.cost
			if old, ok := gScore[link.to]; ok && tentative >= old {
				continue
			}
			cameNode[link.to] = cur.node
			cameEdge[link.to] = link.via
			gScore[link.to] = tentative
			f := tentative + heuristic(net, link.to, q.To, scale)
			if !inOpen[link.to] {
				heap.Push(open, heapItem{node: link.to, f: f, g: tentative})
				inOpen[link.to] = true
			} else {
				open.update(link.to, f, tentative)
			}
		}
	}
	return Result{Expanded: expanded}
}

// --- BFS ---

func searchBFS(net *sim.Network, q Query) Result {
	queue := []sim.NodeID{q.From}
	cameNode := map[sim.NodeID]sim.NodeID{}
	cameEdge := map[sim.NodeID]sim.EdgeID{}
	seen := map[sim.NodeID]bool{q.From: true}
	expanded := 0
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		expanded++
		if cur == q.To {
			nodes, edges := reconstruct(cameNode, cameEdge, q.From, q.To)
			return Result{Found: true, Nodes: nodes, Edges: edges, Cost: pathCost(net, edges), Expanded: expanded}
		}
		for _, link := range linksFrom(net, cur) {
			if seen[link.to] {
				continue
			}
			seen[link.to] = true
			cameNode[link.to] = cur
			cameEdge[link.to] = link.via
			queue = append(queue, link.to)
		}
	}
	return Result{Expanded: expanded}
}

// --- DFS ---

func searchDFS(net *sim.Network, q Query) Result {
	type frame struct {
		node sim.NodeID
		next int // index into links
	}
	linksCache := map[sim.NodeID][]netLink{}
	getLinks := func(n sim.NodeID) []netLink {
		if L, ok := linksCache[n]; ok {
			return L
		}
		L := linksFrom(net, n)
		linksCache[n] = L
		return L
	}

	stack := []frame{{node: q.From}}
	cameNode := map[sim.NodeID]sim.NodeID{}
	cameEdge := map[sim.NodeID]sim.EdgeID{}
	onPath := map[sim.NodeID]bool{q.From: true}
	expanded := 0

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		L := getLinks(top.node)
		if top.next == 0 {
			expanded++
		}
		if top.node == q.To {
			nodes, edges := reconstruct(cameNode, cameEdge, q.From, q.To)
			return Result{Found: true, Nodes: nodes, Edges: edges, Cost: pathCost(net, edges), Expanded: expanded}
		}
		if top.next >= len(L) {
			onPath[top.node] = false
			stack = stack[:len(stack)-1]
			continue
		}
		link := L[top.next]
		top.next++
		if onPath[link.to] {
			continue
		}
		// Allow revisiting nodes not on current path (standard graph DFS for path).
		if _, seen := cameNode[link.to]; seen && link.to != q.To {
			// Keep first discovery for simpler tree; skip if already reached.
			continue
		}
		cameNode[link.to] = top.node
		cameEdge[link.to] = link.via
		onPath[link.to] = true
		stack = append(stack, frame{node: link.to})
	}
	return Result{Expanded: expanded}
}

// --- Bidirectional BFS (meet-in-the-middle) ---

func searchBidirectionalBFS(net *sim.Network, q Query) Result {
	if q.From == q.To {
		return Result{Found: true, Nodes: []sim.NodeID{q.From}}
	}
	type side struct {
		queue    []sim.NodeID
		cameNode map[sim.NodeID]sim.NodeID
		cameEdge map[sim.NodeID]sim.EdgeID
		seen     map[sim.NodeID]bool
	}
	fwd := side{
		queue:    []sim.NodeID{q.From},
		cameNode: map[sim.NodeID]sim.NodeID{},
		cameEdge: map[sim.NodeID]sim.EdgeID{},
		seen:     map[sim.NodeID]bool{q.From: true},
	}
	bwd := side{
		queue:    []sim.NodeID{q.To},
		cameNode: map[sim.NodeID]sim.NodeID{},
		cameEdge: map[sim.NodeID]sim.EdgeID{},
		seen:     map[sim.NodeID]bool{q.To: true},
	}
	expanded := 0
	var meet sim.NodeID

	expand := func(s *side, other *side) bool {
		if len(s.queue) == 0 {
			return false
		}
		cur := s.queue[0]
		s.queue = s.queue[1:]
		expanded++
		for _, link := range linksFrom(net, cur) {
			if s.seen[link.to] {
				continue
			}
			s.seen[link.to] = true
			s.cameNode[link.to] = cur
			s.cameEdge[link.to] = link.via
			if other.seen[link.to] {
				meet = link.to
				return true
			}
			s.queue = append(s.queue, link.to)
		}
		// Also detect if cur already in other (start adjacent cases).
		if other.seen[cur] {
			meet = cur
			return true
		}
		return false
	}

	for len(fwd.queue) > 0 && len(bwd.queue) > 0 {
		// Expand the smaller frontier first (common optimization).
		if len(fwd.queue) <= len(bwd.queue) {
			if expand(&fwd, &bwd) {
				return stitchBidirectional(net, q.From, q.To, meet, fwd, bwd, expanded)
			}
		} else {
			if expand(&bwd, &fwd) {
				return stitchBidirectional(net, q.From, q.To, meet, fwd, bwd, expanded)
			}
		}
	}
	return Result{Expanded: expanded}
}

func stitchBidirectional(net *sim.Network, start, goal, meet sim.NodeID, fwd, bwd struct {
	queue    []sim.NodeID
	cameNode map[sim.NodeID]sim.NodeID
	cameEdge map[sim.NodeID]sim.EdgeID
	seen     map[sim.NodeID]bool
}, expanded int) Result {
	// Path start → meet
	var leftNodes []sim.NodeID
	var leftEdges []sim.EdgeID
	if meet != start {
		n, e := reconstruct(fwd.cameNode, fwd.cameEdge, start, meet)
		if n == nil {
			return Result{Expanded: expanded}
		}
		leftNodes, leftEdges = n, e
	} else {
		leftNodes = []sim.NodeID{start}
	}
	// Path meet → goal is reverse of goal → meet in bwd tree
	var rightNodes []sim.NodeID
	var rightEdges []sim.EdgeID
	if meet != goal {
		n, e := reconstruct(bwd.cameNode, bwd.cameEdge, goal, meet)
		if n == nil {
			return Result{Expanded: expanded}
		}
		// n is goal...meet; reverse to meet...goal and reverse edges
		for i, j := 0, len(n)-1; i < j; i, j = i+1, j-1 {
			n[i], n[j] = n[j], n[i]
		}
		for i, j := 0, len(e)-1; i < j; i, j = i+1, j-1 {
			e[i], e[j] = e[j], e[i]
		}
		// drop duplicate meet at start of right
		rightNodes, rightEdges = n[1:], e
	}
	nodes := append(leftNodes, rightNodes...)
	edges := append(leftEdges, rightEdges...)
	return Result{Found: true, Nodes: nodes, Edges: edges, Cost: pathCost(net, edges), Expanded: expanded}
}

// --- Bidirectional A* / Dijkstra (meet-in-the-middle) ---

func searchBidirectionalAStar(net *sim.Network, q Query, dijkstra bool) Result {
	scale := q.HeuristicScale
	if dijkstra {
		scale = 0
	}

	type half struct {
		open     *nodeHeap
		g        map[sim.NodeID]float32
		cameNode map[sim.NodeID]sim.NodeID
		cameEdge map[sim.NodeID]sim.EdgeID
		closed   map[sim.NodeID]bool
		inOpen   map[sim.NodeID]bool
		target   sim.NodeID
	}
	newHalf := func(start, target sim.NodeID) *half {
		h := &half{
			open:     &nodeHeap{},
			g:        map[sim.NodeID]float32{start: 0},
			cameNode: map[sim.NodeID]sim.NodeID{},
			cameEdge: map[sim.NodeID]sim.EdgeID{},
			closed:   map[sim.NodeID]bool{},
			inOpen:   map[sim.NodeID]bool{start: true},
			target:   target,
		}
		heap.Init(h.open)
		heap.Push(h.open, heapItem{node: start, f: heuristic(net, start, target, scale), g: 0})
		return h
	}
	fwd := newHalf(q.From, q.To)
	bwd := newHalf(q.To, q.From)

	expanded := 0
	bestCost := float32(math.MaxFloat32)
	var meet sim.NodeID
	found := false

	considerMeet := func(n sim.NodeID, total float32) {
		if total < bestCost {
			bestCost = total
			meet = n
			found = true
		}
	}

	expandSide := func(s, other *half) {
		if s.open.Len() == 0 {
			return
		}
		cur := heap.Pop(s.open).(heapItem)
		s.inOpen[cur.node] = false
		if s.closed[cur.node] {
			return
		}
		// Stale heap entry.
		if gs, ok := s.g[cur.node]; ok && cur.g > gs+1e-4 {
			return
		}
		s.closed[cur.node] = true
		expanded++

		if og, ok := other.g[cur.node]; ok {
			considerMeet(cur.node, s.g[cur.node]+og)
		}

		for _, link := range linksFrom(net, cur.node) {
			if s.closed[link.to] {
				continue
			}
			tentative := s.g[cur.node] + link.cost
			if old, ok := s.g[link.to]; ok && tentative >= old {
				continue
			}
			s.cameNode[link.to] = cur.node
			s.cameEdge[link.to] = link.via
			s.g[link.to] = tentative
			f := tentative + heuristic(net, link.to, s.target, scale)
			if !s.inOpen[link.to] {
				heap.Push(s.open, heapItem{node: link.to, f: f, g: tentative})
				s.inOpen[link.to] = true
			} else {
				s.open.update(link.to, f, tentative)
			}
			if og, ok := other.g[link.to]; ok {
				considerMeet(link.to, tentative+og)
			}
		}
	}

	for fwd.open.Len() > 0 && bwd.open.Len() > 0 {
		if found {
			fMin := (*fwd.open)[0].f
			bMin := (*bwd.open)[0].f
			if fMin+bMin >= bestCost {
				break
			}
		}
		if fwd.open.Len() <= bwd.open.Len() {
			expandSide(fwd, bwd)
		} else {
			expandSide(bwd, fwd)
		}
	}

	if !found {
		return Result{Expanded: expanded}
	}

	fwdSide := struct {
		queue    []sim.NodeID
		cameNode map[sim.NodeID]sim.NodeID
		cameEdge map[sim.NodeID]sim.EdgeID
		seen     map[sim.NodeID]bool
	}{cameNode: fwd.cameNode, cameEdge: fwd.cameEdge}
	bwdSide := struct {
		queue    []sim.NodeID
		cameNode map[sim.NodeID]sim.NodeID
		cameEdge map[sim.NodeID]sim.EdgeID
		seen     map[sim.NodeID]bool
	}{cameNode: bwd.cameNode, cameEdge: bwd.cameEdge}

	res := stitchBidirectional(net, q.From, q.To, meet, fwdSide, bwdSide, expanded)
	if res.Found {
		res.Cost = bestCost
	}
	return res
}
