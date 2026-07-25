package pathfind

import (
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func TestDStarLiteBasicPath(t *testing.T) {
	net, ids := lineGraph(t)
	d := NewDStarLite(net, ids[0], ids[4])
	r := d.Replan()
	if !r.Found || len(r.Edges) != 4 {
		t.Fatalf("%+v", r)
	}
}

func TestDStarLiteReplansAroundBlock(t *testing.T) {
	net, ids := lineGraph(t)
	p1, _ := net.GetNode(ids[1])
	p3, _ := net.GetNode(ids[3])
	delta := p3.Pos.Sub(p1.Pos)
	bypass := net.AddEdge(ids[1], ids[3], p1.Pos.Add(delta.Scale(1.0/3)), p1.Pos.Add(delta.Scale(2.0/3)))

	d := NewDStarLite(net, ids[0], ids[4])
	if !d.Replan().Found {
		t.Fatal("initial")
	}

	var e12 sim.EdgeID
	for _, link := range linksFrom(net, ids[1]) {
		if link.to == ids[2] {
			e12 = link.via
			break
		}
	}
	if e12 == sim.NilEdge {
		t.Fatal("missing 1-2")
	}
	net.SetEdgeBlocked(e12, true)

	r2 := d.Replan()
	if !r2.Found {
		t.Fatal("expected replan around block")
	}
	usedBypass := false
	for _, e := range r2.Edges {
		if e == e12 {
			t.Fatal("path still uses blocked edge")
		}
		if e == bypass {
			usedBypass = true
		}
	}
	if !usedBypass {
		t.Fatalf("expected bypass in %v", r2.Edges)
	}
}

func TestDStarLiteSetStartAdvances(t *testing.T) {
	net, ids := lineGraph(t)
	d := NewDStarLite(net, ids[0], ids[4])
	_ = d.Replan()
	d.SetStart(ids[2])
	r := d.Replan()
	if !r.Found || r.Nodes[0] != ids[2] || len(r.Edges) != 2 {
		t.Fatalf("%+v", r)
	}
}

func TestDynamicRouteInvalidatesOnBlock(t *testing.T) {
	net, ids := lineGraph(t)
	var dr DynamicRoute
	if !dr.Ensure(net, ids[0], ids[4]).Found {
		t.Fatal("initial")
	}
	var mid sim.EdgeID
	for _, link := range linksFrom(net, ids[2]) {
		if link.to == ids[3] {
			mid = link.via
			break
		}
	}
	net.SetEdgeBlocked(mid, true)
	if dr.Ensure(net, ids[0], ids[4]).Found {
		t.Fatal("line should be disconnected")
	}
}

func TestEdgeCostMulAffectsDijkstra(t *testing.T) {
	net := sim.NewWorld().Network
	a := net.AddNode(sim.Vec2{X: 0, Y: 0})
	b := net.AddNode(sim.Vec2{X: 50, Y: 10})
	c := net.AddNode(sim.Vec2{X: 50, Y: -10})
	d := net.AddNode(sim.Vec2{X: 100, Y: 0})
	add := func(from, to sim.NodeID) sim.EdgeID {
		pa, _ := net.GetNode(from)
		pb, _ := net.GetNode(to)
		delta := pb.Pos.Sub(pa.Pos)
		return net.AddEdge(from, to, pa.Pos.Add(delta.Scale(1.0/3)), pa.Pos.Add(delta.Scale(2.0/3)))
	}
	ab := add(a, b)
	add(b, d)
	add(a, c)
	add(c, d)
	if Dijkstra(net, a, d).Nodes[1] != b {
		t.Fatal("want B")
	}
	net.SetEdgeCostMul(ab, 50)
	if Dijkstra(net, a, d).Nodes[1] != c {
		t.Fatal("want C after cost mul")
	}
}
