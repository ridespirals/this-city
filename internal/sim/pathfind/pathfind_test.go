package pathfind

import (
	"math"
	"testing"

	"github.com/ridespirals/this-city/internal/sim"
)

func lineGraph(t *testing.T) (*sim.Network, []sim.NodeID) {
	t.Helper()
	w := sim.NewWorld()
	net := w.Network
	ids := make([]sim.NodeID, 5)
	for i := 0; i < 5; i++ {
		ids[i] = net.AddNode(sim.Vec2{X: float32(i) * 100, Y: 0})
	}
	for i := 0; i < 4; i++ {
		a, b := ids[i], ids[i+1]
		pa, _ := net.GetNode(a)
		pb, _ := net.GetNode(b)
		d := pb.Pos.Sub(pa.Pos)
		net.AddEdge(a, b, pa.Pos.Add(d.Scale(1.0/3)), pa.Pos.Add(d.Scale(2.0/3)))
	}
	return net, ids
}

func TestAStarAndDijkstraShortest(t *testing.T) {
	net, ids := lineGraph(t)
	a := AStar(net, ids[0], ids[4])
	d := Dijkstra(net, ids[0], ids[4])
	if !a.Found || !d.Found {
		t.Fatalf("found a=%v d=%v", a.Found, d.Found)
	}
	if len(a.Edges) != 4 || len(d.Edges) != 4 {
		t.Fatalf("edges a=%d d=%d", len(a.Edges), len(d.Edges))
	}
	if math.Abs(float64(a.Cost-d.Cost)) > 1e-3 {
		t.Fatalf("cost a=%v d=%v", a.Cost, d.Cost)
	}
}

func TestBFSHopCount(t *testing.T) {
	net, ids := lineGraph(t)
	r := BFS(net, ids[0], ids[3])
	if !r.Found || len(r.Edges) != 3 {
		t.Fatalf("got %+v", r)
	}
}

func TestDFSFindsPath(t *testing.T) {
	net, ids := lineGraph(t)
	r := DFS(net, ids[0], ids[4])
	if !r.Found || r.Nodes[0] != ids[0] || r.Nodes[len(r.Nodes)-1] != ids[4] {
		t.Fatalf("%+v", r)
	}
}

func TestBidirectionalMeetInMiddle(t *testing.T) {
	net, ids := lineGraph(t)
	for _, algo := range []Algo{AlgoBidirectionalBFS, AlgoBidirectionalAStar, AlgoBidirectionalDijkstra} {
		r := Find(net, Query{From: ids[0], To: ids[4], Algo: algo, HeuristicScale: 1})
		if !r.Found || len(r.Edges) != 4 {
			t.Fatalf("algo %v %+v", algo, r)
		}
	}
}

func TestNoPath(t *testing.T) {
	net := sim.NewWorld().Network
	a := net.AddNode(sim.Vec2{})
	b := net.AddNode(sim.Vec2{X: 10})
	if AStar(net, a, b).Found {
		t.Fatal("expected no path")
	}
}

func TestSameNode(t *testing.T) {
	net, ids := lineGraph(t)
	r := AStar(net, ids[2], ids[2])
	if !r.Found || len(r.Edges) != 0 || r.Cost != 0 {
		t.Fatalf("%+v", r)
	}
}

func TestFigureEightPath(t *testing.T) {
	w := sim.NewWorld()
	if err := sim.ApplyMapFile(sim.FigureEightMap(), w.Network); err != nil {
		t.Fatal(err)
	}
	var left, right sim.NodeID
	w.Network.ForEachNode(func(n *sim.Node) {
		if n.Pos.X == 360 && n.Pos.Y == 360 {
			left = n.ID
		}
		if n.Pos.X == 920 && n.Pos.Y == 360 {
			right = n.ID
		}
	})
	if left == sim.NilNode || right == sim.NilNode {
		t.Fatal("tips not found")
	}
	r := AStar(w.Network, left, right)
	if !r.Found || len(r.Edges) != 2 {
		t.Fatalf("%+v", r)
	}
}

func TestSetRoute(t *testing.T) {
	var d sim.PathDecision
	d.SetRoute([]sim.EdgeID{1, 2, 3})
	if d.Mode != sim.DecideRoute || len(d.Route) != 3 {
		t.Fatalf("%+v", d)
	}
}

func TestBranchPrefersShorter(t *testing.T) {
	net := sim.NewWorld().Network
	a := net.AddNode(sim.Vec2{X: 0, Y: 0})
	b := net.AddNode(sim.Vec2{X: 50, Y: 10})
	c := net.AddNode(sim.Vec2{X: 50, Y: -100})
	d := net.AddNode(sim.Vec2{X: 100, Y: 0})
	add := func(from, to sim.NodeID) {
		pa, _ := net.GetNode(from)
		pb, _ := net.GetNode(to)
		delta := pb.Pos.Sub(pa.Pos)
		net.AddEdge(from, to, pa.Pos.Add(delta.Scale(1.0/3)), pa.Pos.Add(delta.Scale(2.0/3)))
	}
	add(a, b)
	add(b, d)
	add(a, c)
	add(c, d)
	r := Dijkstra(net, a, d)
	if !r.Found || r.Nodes[1] != b {
		t.Fatalf("%+v", r)
	}
}
