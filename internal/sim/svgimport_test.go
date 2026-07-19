package sim

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseSVGPathLineAndClose(t *testing.T) {
	curves, err := ParseSVGPathData("M 0,0 L 10,0 L 10,10 Z")
	if err != nil {
		t.Fatal(err)
	}
	if len(curves) != 3 {
		t.Fatalf("curves=%d want 3", len(curves))
	}
	if curves[0].P0 != (Vec2{}) || curves[0].P1 != (Vec2{X: 10}) {
		t.Fatalf("first seg %+v", curves[0])
	}
	last := curves[len(curves)-1]
	if last.P1 != (Vec2{}) {
		t.Fatalf("close end %+v", last.P1)
	}
}

func TestParseDevMapSVG(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "assets", "maps", "dev-map.svg"))
	if err != nil {
		t.Skip(err)
	}
	pieces, err := ParseSVG(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(pieces) != 1 {
		t.Fatalf("pieces=%d want 1", len(pieces))
	}
	if len(pieces[0].Curves) != 8 {
		t.Fatalf("curves=%d want 8", len(pieces[0].Curves))
	}
	mf := PiecesToMapFile("dev-map", pieces, DefaultSVGMergeEps)
	if len(mf.Nodes) != 8 {
		t.Fatalf("nodes=%d want 8", len(mf.Nodes))
	}
	if len(mf.Edges) != 8 {
		t.Fatalf("edges=%d want 8", len(mf.Edges))
	}
	net := newNetwork()
	if err := ApplyMapFile(mf, net); err != nil {
		t.Fatal(err)
	}
	if net.EdgeCount() != 8 || net.NodeCount() != 8 {
		t.Fatalf("net nodes=%d edges=%d", net.NodeCount(), net.EdgeCount())
	}
}

func TestStampPieceMultiple(t *testing.T) {
	piece := PathPiece{
		Name: "seg",
		Curves: []CubicBezier{
			lineCubic(Vec2{}, Vec2{X: 100}),
			lineCubic(Vec2{X: 100}, Vec2{X: 100, Y: 50}),
		},
	}.Recentered()

	net := newNetwork()
	g1, e1 := StampPiece(net, piece, Vec2{X: 200, Y: 200}, 1)
	g2, e2 := StampPiece(net, piece, Vec2{X: 500, Y: 200}, 1)
	if g1 == 0 || g2 == 0 || g1 == g2 {
		t.Fatalf("groups %d %d", g1, g2)
	}
	if len(e1) != 2 || len(e2) != 2 {
		t.Fatalf("edges %d %d", len(e1), len(e2))
	}
	if net.EdgeCount() != 4 {
		t.Fatalf("edge count %d", net.EdgeCount())
	}
}
