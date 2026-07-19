package sim

import (
	"math/rand"
	"testing"
)

func TestCommandKeyMapLoads(t *testing.T) {
	w := NewWorld()
	if err := ApplyMapFile(CommandKeyMap(), w.Network); err != nil {
		t.Fatal(err)
	}
	if w.Network.NodeCount() != 8 || w.Network.EdgeCount() != 12 {
		t.Fatalf("nodes=%d edges=%d", w.Network.NodeCount(), w.Network.EdgeCount())
	}
	var mt NodeID
	for id, n := range w.Network.nodes {
		if n.Pos.X == 640 && n.Pos.Y == 250 {
			mt = id
			break
		}
	}
	if mt == NilNode || len(w.Network.IncidentEdges(mt)) != 4 {
		t.Fatalf("Mt id=%d degree=%d want 4", mt, len(w.Network.IncidentEdges(mt)))
	}
}

func TestFollowerChoosesAtJunction(t *testing.T) {
	w := NewWorld()
	w.RNG = rand.New(rand.NewSource(42))
	if err := ApplyMapFile(CommandKeyMap(), w.Network); err != nil {
		t.Fatal(err)
	}
	var start EdgeID
	w.Network.ForEachEdge(func(e *Edge) {
		if start == NilEdge {
			start = e.ID
		}
	})
	ent := w.Create()
	w.Transforms.Set(ent, Transform2D{Scale: 1})
	w.Decisions.Set(ent, DefaultPathDecision())
	PlaceOnEdge(w, ent, start, 0, true, 200)

	seen := map[EdgeID]bool{}
	for i := 0; i < 300; i++ {
		TickPathFollowers(w, 1.0/60)
		f, _ := w.Followers.Get(ent)
		seen[f.Edge] = true
	}
	if len(seen) < 2 {
		t.Fatalf("expected multiple edges over time, got %d", len(seen))
	}
}

func TestDecideRouteOverridesRandom(t *testing.T) {
	w := NewWorld()
	n0 := w.Network.AddNode(Vec2{0, 0})
	n1 := w.Network.AddNode(Vec2{100, 0})
	n2 := w.Network.AddNode(Vec2{100, 100})
	e1 := w.Network.AddEdge(n0, n1, Vec2{33, 0}, Vec2{66, 0})
	e2 := w.Network.AddEdge(n1, n2, Vec2{100, 33}, Vec2{100, 66})
	_ = w.Network.AddEdge(n1, n0, Vec2{66, 10}, Vec2{33, 10}) // distractor

	dec := PathDecision{Mode: DecideRoute, Route: []EdgeID{e2}, AvoidUTurn: true}
	arr := Arrival{Node: n1, ViaEdge: e1, Forward: true}
	choice, dec2 := ChooseNext(w.Network, dec, arr, w.RNG)
	if choice.Edge != e2 {
		t.Fatalf("choice=%v want %v", choice.Edge, e2)
	}
	if len(dec2.Route) != 0 {
		t.Fatal("route should be consumed")
	}
}

func TestMapFileRoundTripJSON(t *testing.T) {
	mf := CommandKeyMap()
	data, err := MarshalMapFile(mf)
	if err != nil {
		t.Fatal(err)
	}
	w := NewWorld()
	if err := LoadNetworkJSON(data, w.Network); err != nil {
		t.Fatal(err)
	}
	if w.Network.EdgeCount() != 12 {
		t.Fatalf("edges=%d", w.Network.EdgeCount())
	}
}
