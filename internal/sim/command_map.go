package sim

// CommandKeyMap returns a ⌘-shaped network: two horizontal figure-8s stacked
// so they share left/right junctions (JL/JR), matching the Mac command symbol.
func CommandKeyMap() MapFile {
	// Nodes:
	//   Mt / Mb — figure-8 crossings (degree 4)
	//   JL / JR — stack intersections between the two 8s (degree 4)
	//   Lt, Rt, Lb, Rb — outer lobe tips
	const (
		nMt = 1
		nMb = 2
		nJL = 3
		nJR = 4
		nLt = 5
		nRt = 6
		nLb = 7
		nRb = 8
	)
	nodes := []MapFileNode{
		{ID: nMt, X: 640, Y: 250},
		{ID: nMb, X: 640, Y: 470},
		{ID: nJL, X: 470, Y: 360},
		{ID: nJR, X: 810, Y: 360},
		{ID: nLt, X: 380, Y: 210},
		{ID: nRt, X: 900, Y: 210},
		{ID: nLb, X: 380, Y: 510},
		{ID: nRb, X: 900, Y: 510},
	}

	// Helper: control points bulging toward a third point for a rounded lobe.
	edge := func(id, from, to uint32, c0x, c0y, c1x, c1y float32) MapFileEdge {
		return MapFileEdge{
			ID: id, From: from, To: to,
			C0: MapFileXY{X: c0x, Y: c0y},
			C1: MapFileXY{X: c1x, Y: c1y},
		}
	}

	edges := []MapFileEdge{
		// Top figure-8 left loop: Mt — Lt — JL — Mt
		edge(1, nMt, nLt, 520, 180, 420, 180),
		edge(2, nLt, nJL, 360, 260, 400, 320),
		edge(3, nJL, nMt, 520, 320, 580, 280),
		// Top figure-8 right loop: Mt — Rt — JR — Mt
		edge(4, nMt, nRt, 760, 180, 860, 180),
		edge(5, nRt, nJR, 920, 260, 880, 320),
		edge(6, nJR, nMt, 760, 320, 700, 280),
		// Bottom figure-8 left loop: Mb — Lb — JL — Mb
		edge(7, nMb, nLb, 520, 540, 420, 540),
		edge(8, nLb, nJL, 360, 460, 400, 400),
		edge(9, nJL, nMb, 520, 400, 580, 440),
		// Bottom figure-8 right loop: Mb — Rb — JR — Mb
		edge(10, nMb, nRb, 760, 540, 860, 540),
		edge(11, nRb, nJR, 920, 460, 880, 400),
		edge(12, nJR, nMb, 760, 400, 700, 440),
	}

	return MapFile{Name: "command-key", Nodes: nodes, Edges: edges}
}
