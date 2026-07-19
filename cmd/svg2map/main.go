// Command svg2map converts an SVG with <path d="..."> elements into map JSON.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ridespirals/this-city/internal/sim"
)

func main() {
	outPath := flag.String("o", "", "output JSON path (default: stdout)")
	name := flag.String("name", "", "map name (default: from input filename)")
	eps := flag.Float64("eps", float64(sim.DefaultSVGMergeEps), "node merge epsilon")
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: svg2map [-o out.json] [-name name] [-eps 1.5] input.svg\n")
		os.Exit(2)
	}
	inPath := flag.Arg(0)
	data, err := os.ReadFile(inPath)
	if err != nil {
		fail(err)
	}
	mapName := *name
	if mapName == "" {
		mapName = sim.PieceNameFromPath(inPath)
	}
	mf, err := sim.MapFileFromSVG(mapName, data, float32(*eps))
	if err != nil {
		fail(err)
	}
	raw, err := sim.MarshalMapFile(mf)
	if err != nil {
		fail(err)
	}
	raw = append(raw, '\n')
	if *outPath == "" {
		_, _ = os.Stdout.Write(raw)
		return
	}
	if err := os.WriteFile(*outPath, raw, 0o644); err != nil {
		fail(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s (%d nodes, %d edges)\n", *outPath, len(mf.Nodes), len(mf.Edges))
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "svg2map: %v\n", err)
	os.Exit(1)
}
