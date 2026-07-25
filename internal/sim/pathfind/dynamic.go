package pathfind

import "github.com/ridespirals/this-city/internal/sim"

// DynamicRoute is a lightweight replanner: full Find when start, goal, or
// sim.Network.Version changes. Use when changes are infrequent; prefer DStarLite
// when many edge costs flip while the agent is en route.
type DynamicRoute struct {
	Algo           Algo
	HeuristicScale float32

	from   sim.NodeID
	to     sim.NodeID
	netVer uint64
	cached Result
}

// Ensure returns a route from→to, recomputing only when needed.
func (d *DynamicRoute) Ensure(net *sim.Network, from, to sim.NodeID) Result {
	if d == nil {
		return Find(net, Query{From: from, To: to, Algo: AlgoAStar, HeuristicScale: 1})
	}
	ver := uint64(0)
	if net != nil {
		ver = net.Version()
	}
	scale := d.HeuristicScale
	algo := d.Algo
	if algo == 0 && d.HeuristicScale == 0 && d.from == sim.NilNode {
		// zero value → A*
		algo = AlgoAStar
		scale = 1
	}
	if d.cached.Found && from == d.from && to == d.to && ver == d.netVer {
		return d.cached
	}
	if scale == 0 && (algo == AlgoAStar || algo == AlgoBidirectionalAStar) {
		scale = 1
	}
	d.from, d.to, d.netVer = from, to, ver
	d.cached = Find(net, Query{From: from, To: to, Algo: algo, HeuristicScale: scale})
	return d.cached
}

// Invalidate forces the next Ensure to replan.
func (d *DynamicRoute) Invalidate() {
	if d == nil {
		return
	}
	d.netVer = ^uint64(0)
	d.cached = Result{}
}

// ApplyToDecision writes the ensured route into a sim.PathDecision when found.
func (d *DynamicRoute) ApplyToDecision(net *sim.Network, from, to sim.NodeID, dec *sim.PathDecision) bool {
	r := d.Ensure(net, from, to)
	if !r.Found || dec == nil {
		return false
	}
	dec.SetRoute(r.Edges)
	return true
}
