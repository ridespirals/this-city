package sim

import (
	"encoding/xml"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// DefaultSVGMergeEps is the node-snap distance when importing SVG paths.
const DefaultSVGMergeEps float32 = 1.5

// PathPiece is a reusable local-space path fragment (from SVG or elsewhere)
// that can be stamped into a Network multiple times.
type PathPiece struct {
	Name   string
	Curves []CubicBezier
}

// ParseSVG extracts every <path d="..."> into a PathPiece (absolute coordinates).
func ParseSVG(data []byte) ([]PathPiece, error) {
	dec := xml.NewDecoder(strings.NewReader(string(data)))
	dec.Strict = false
	var pieces []PathPiece
	var n int
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("svg xml: %w", err)
		}
		se, ok := tok.(xml.StartElement)
		if !ok || !strings.EqualFold(se.Name.Local, "path") {
			continue
		}
		var id, d string
		for _, a := range se.Attr {
			switch strings.ToLower(a.Name.Local) {
			case "id":
				id = a.Value
			case "d":
				d = a.Value
			}
		}
		if strings.TrimSpace(d) == "" {
			continue
		}
		curves, err := ParseSVGPathData(d)
		if err != nil {
			return nil, fmt.Errorf("path %q: %w", id, err)
		}
		if len(curves) == 0 {
			continue
		}
		n++
		name := id
		if name == "" {
			name = fmt.Sprintf("path-%d", n)
		}
		pieces = append(pieces, PathPiece{Name: name, Curves: curves})
	}
	if len(pieces) == 0 {
		return nil, fmt.Errorf("svg: no path elements with d= found")
	}
	return pieces, nil
}

// MapFileFromSVG parses an SVG document into a MapFile (absolute coords, nodes snapped).
func MapFileFromSVG(name string, data []byte, mergeEps float32) (MapFile, error) {
	pieces, err := ParseSVG(data)
	if err != nil {
		return MapFile{}, err
	}
	if name == "" {
		name = "svg-import"
	}
	return PiecesToMapFile(name, pieces, mergeEps), nil
}

// PiecesToMapFile stamps every piece at the origin into a fresh MapFile.
func PiecesToMapFile(name string, pieces []PathPiece, mergeEps float32) MapFile {
	net := newNetwork()
	for _, p := range pieces {
		StampPiece(net, p, Vec2{}, mergeEps)
	}
	return ExportMapFile(name, net)
}

// Bounds returns axis-aligned min/max of all curve endpoints and controls.
func (p PathPiece) Bounds() (min, max Vec2) {
	if len(p.Curves) == 0 {
		return Vec2{}, Vec2{}
	}
	min = p.Curves[0].P0
	max = min
	grow := func(v Vec2) {
		if v.X < min.X {
			min.X = v.X
		}
		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.X > max.X {
			max.X = v.X
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
	}
	for _, c := range p.Curves {
		grow(c.P0)
		grow(c.C0)
		grow(c.C1)
		grow(c.P1)
	}
	return min, max
}

// Translated returns a copy shifted by delta.
func (p PathPiece) Translated(delta Vec2) PathPiece {
	out := PathPiece{Name: p.Name, Curves: make([]CubicBezier, len(p.Curves))}
	for i, c := range p.Curves {
		out.Curves[i] = CubicBezier{
			P0: c.P0.Add(delta),
			C0: c.C0.Add(delta),
			C1: c.C1.Add(delta),
			P1: c.P1.Add(delta),
		}
	}
	return out
}

// Recentered returns a copy with bounding-box center at the origin (good for stamping).
func (p PathPiece) Recentered() PathPiece {
	min, max := p.Bounds()
	center := Vec2{X: (min.X + max.X) * 0.5, Y: (min.Y + max.Y) * 0.5}
	out := p.Translated(Vec2{X: -center.X, Y: -center.Y})
	out.Name = p.Name
	return out
}

// StampPiece places piece so its local origin lands at `at`, merging nodes within mergeEps.
// Returns the new editor group id and created edge ids.
func StampPiece(net *Network, piece PathPiece, at Vec2, mergeEps float32) (uint32, []EdgeID) {
	if net == nil || len(piece.Curves) == 0 {
		return 0, nil
	}
	if mergeEps <= 0 {
		mergeEps = DefaultSVGMergeEps
	}
	placed := piece.Translated(at)
	group := net.NewGroup()
	var edges []EdgeID
	for _, c := range placed.Curves {
		if segNearZero(c) {
			continue
		}
		from := net.EnsureNode(c.P0, mergeEps)
		to := net.EnsureNode(c.P1, mergeEps)
		if from == to {
			// Degenerate after snap; skip unless it still has meaningful length.
			mid := c.Point(0.5)
			d := mid.Sub(c.P0)
			if d.Dot(d) < mergeEps*mergeEps {
				continue
			}
		}
		id := net.AddEdgeGrouped(from, to, c.C0, c.C1, group)
		if id != NilEdge {
			edges = append(edges, id)
		}
	}
	return group, edges
}

// SampleStroke returns polyline samples per curve (for ghosts/previews; committed
// network edges still own Edge.Poly).
func (p PathPiece) SampleStroke(stepsPerSeg int) [][]Vec2 {
	if stepsPerSeg < 1 {
		stepsPerSeg = 1
	}
	out := make([][]Vec2, 0, len(p.Curves))
	for _, c := range p.Curves {
		out = append(out, c.Sample(stepsPerSeg))
	}
	return out
}

// PieceNameFromPath derives a piece/map name from a file path.
func PieceNameFromPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
