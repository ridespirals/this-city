package pathfind

import "github.com/ridespirals/this-city/internal/sim"

type netLink struct {
	to   sim.NodeID
	via  sim.EdgeID
	cost float32
}

func linksFrom(net *sim.Network, node sim.NodeID) []netLink {
	if net == nil {
		return nil
	}
	ids := net.IncidentEdges(node)
	out := make([]netLink, 0, len(ids))
	for _, eid := range ids {
		cost, ok := net.TraversalCost(eid)
		if !ok {
			continue
		}
		other := net.OtherNode(eid, node)
		if other == sim.NilNode {
			continue
		}
		out = append(out, netLink{to: other, via: eid, cost: cost})
	}
	return out
}

func heuristic(net *sim.Network, a, b sim.NodeID, scale float32) float32 {
	if net == nil {
		return 0
	}
	na, okA := net.GetNode(a)
	nb, okB := net.GetNode(b)
	if !okA || !okB {
		return 0
	}
	return na.Pos.Sub(nb.Pos).Len() * scale
}

func pathCost(net *sim.Network, edges []sim.EdgeID) float32 {
	var c float32
	for _, eid := range edges {
		if cost, ok := net.TraversalCost(eid); ok {
			c += cost
		}
	}
	return c
}
