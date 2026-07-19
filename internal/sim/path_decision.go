package sim

import "math/rand"

// DecisionMode selects how an agent picks the next edge at a junction.
type DecisionMode int

const (
	// DecideRandom picks uniformly among legal outgoing choices (default).
	DecideRandom DecisionMode = iota
	// DecideRoute follows PathDecision.Route (A*/Dijkstra output or scripted path).
	DecideRoute
)

// PathDecision controls intersection behavior for a PathFollower.
type PathDecision struct {
	Mode DecisionMode
	// Route is a planned sequence of edge ids to traverse (ConsumeRoute advances it).
	Route []EdgeID
	// AvoidUTurn excludes the edge just traveled when other options exist.
	AvoidUTurn bool
}

// DefaultPathDecision returns random choice with U-turn avoidance.
func DefaultPathDecision() PathDecision {
	return PathDecision{Mode: DecideRandom, AvoidUTurn: true}
}

// Arrival is the context for a junction decision.
type Arrival struct {
	Node    NodeID
	ViaEdge EdgeID // edge just completed
	Forward bool   // travel direction on ViaEdge when arriving
}

// NextChoice is the edge and travel direction to take next.
type NextChoice struct {
	Edge    EdgeID
	Forward bool // true = travel From→To on Edge
}

// ChooseNext picks the next edge at a junction.
// If DecideRoute has remaining edges, those win; otherwise random (or reverse if stuck).
func ChooseNext(net *Network, d PathDecision, arr Arrival, rng *rand.Rand) (NextChoice, PathDecision) {
	if net == nil {
		return NextChoice{}, d
	}

	if d.Mode == DecideRoute && len(d.Route) > 0 {
		next := d.Route[0]
		d.Route = d.Route[1:]
		if e, ok := net.GetEdge(next); ok {
			fwd := e.From == arr.Node
			if e.From == arr.Node || e.To == arr.Node {
				return NextChoice{Edge: next, Forward: fwd}, d
			}
		}
		// Invalid route entry — fall through to random.
	}

	candidates := junctionChoices(net, arr, d.AvoidUTurn)
	if len(candidates) == 0 {
		// Dead end: reverse on the same edge.
		return NextChoice{Edge: arr.ViaEdge, Forward: !arr.Forward}, d
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	pick := candidates[rng.Intn(len(candidates))]
	return pick, d
}

func junctionChoices(net *Network, arr Arrival, avoidUTurn bool) []NextChoice {
	var out []NextChoice
	for _, eid := range net.IncidentEdges(arr.Node) {
		e, ok := net.GetEdge(eid)
		if !ok {
			continue
		}
		if avoidUTurn && eid == arr.ViaEdge {
			continue
		}
		fwd := e.From == arr.Node
		if e.From != arr.Node && e.To != arr.Node {
			continue
		}
		out = append(out, NextChoice{Edge: eid, Forward: fwd})
	}
	if len(out) == 0 && avoidUTurn {
		// Only U-turn available.
		return junctionChoices(net, arr, false)
	}
	return out
}
