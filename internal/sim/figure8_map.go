package sim

// FigureEightMap returns a classic horizontal figure-8 (∞) street network:
// one crossing in the middle and rounded left/right lobes.
func FigureEightMap() MapFile {
	const (
		nC = 1 // center crossing (degree 4)
		nL = 2 // left lobe tip
		nR = 3 // right lobe tip
	)
	// Centered in a 1280×720 view.
	nodes := []MapFileNode{
		{ID: nC, X: 640, Y: 360},
		{ID: nL, X: 360, Y: 360},
		{ID: nR, X: 920, Y: 360},
	}

	edge := func(id, from, to uint32, c0x, c0y, c1x, c1y float32) MapFileEdge {
		return MapFileEdge{
			ID: id, From: from, To: to,
			C0: MapFileXY{X: c0x, Y: c0y},
			C1: MapFileXY{X: c1x, Y: c1y},
		}
	}

	// Four arcs through the center: upper/lower left, upper/lower right.
	edges := []MapFileEdge{
		edge(1, nC, nL, 560, 220, 420, 220), // upper left
		edge(2, nL, nC, 420, 500, 560, 500), // lower left
		edge(3, nC, nR, 720, 220, 860, 220), // upper right
		edge(4, nR, nC, 860, 500, 720, 500), // lower right
	}

	return MapFile{Name: "figure-8", Nodes: nodes, Edges: edges}
}

// CommandKeyMap is a deprecated alias kept for older callers/tests.
// Prefer FigureEightMap.
func CommandKeyMap() MapFile { return FigureEightMap() }
