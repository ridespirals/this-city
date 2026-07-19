// Package mapsvg embeds SVG path sources for maps and stampable pieces.
package mapsvg

import (
	"embed"
	"io/fs"
	"strings"

	"github.com/ridespirals/this-city/internal/sim"
)

//go:embed *.svg
var FS embed.FS

//go:embed dev-map.svg
var DevMapSVG []byte

// LoadStampPieces parses every embedded SVG path into recentered PathPieces
// suitable for editor stamping (origin at bounding-box center).
func LoadStampPieces() ([]sim.PathPiece, error) {
	entries, err := fs.ReadDir(FS, ".")
	if err != nil {
		return nil, err
	}
	var out []sim.PathPiece
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".svg") {
			continue
		}
		data, err := FS.ReadFile(e.Name())
		if err != nil {
			return nil, err
		}
		pieces, err := sim.ParseSVG(data)
		if err != nil {
			return nil, err
		}
		base := sim.PieceNameFromPath(e.Name())
		for _, p := range pieces {
			name := p.Name
			if name == "" || strings.HasPrefix(name, "path-") {
				name = base
			} else {
				name = base + "/" + name
			}
			p.Name = name
			out = append(out, p.Recentered())
		}
	}
	return out, nil
}

// DevMapFile converts the embedded dev-map.svg into a MapFile (absolute coords).
func DevMapFile() (sim.MapFile, error) {
	return sim.MapFileFromSVG("dev-map", DevMapSVG, sim.DefaultSVGMergeEps)
}
