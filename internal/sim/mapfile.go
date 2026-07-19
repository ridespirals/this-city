package sim

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// MapFile is the on-disk / embeddable city map format (nodes + bezier edges).
type MapFile struct {
	Name  string        `json:"name"`
	Nodes []MapFileNode `json:"nodes"`
	Edges []MapFileEdge `json:"edges"`
}

// MapFileNode is a junction in the map file.
type MapFileNode struct {
	ID uint32  `json:"id"`
	X  float32 `json:"x"`
	Y  float32 `json:"y"`
}

// MapFileEdge connects two nodes with cubic control points.
type MapFileEdge struct {
	ID   uint32    `json:"id"`
	From uint32    `json:"from"`
	To   uint32    `json:"to"`
	C0   MapFileXY `json:"c0"`
	C1   MapFileXY `json:"c1"`
}

// MapFileXY is a 2D point in a map file.
type MapFileXY struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

// LoadNetworkJSON parses a MapFile and populates net (replacing contents).
func LoadNetworkJSON(data []byte, net *Network) error {
	var mf MapFile
	if err := json.Unmarshal(data, &mf); err != nil {
		return err
	}
	return ApplyMapFile(mf, net)
}

// LoadNetworkFile reads a map JSON file into net.
func LoadNetworkFile(path string, net *Network) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return LoadNetworkJSON(data, net)
}

// ApplyMapFile clears net and loads nodes/edges. File ids are remapped to fresh Network ids.
func ApplyMapFile(mf MapFile, net *Network) error {
	if net == nil {
		return fmt.Errorf("nil network")
	}
	*net = *newNetwork()
	idMap := make(map[uint32]NodeID, len(mf.Nodes))
	for _, n := range mf.Nodes {
		nid := net.AddNode(Vec2{X: n.X, Y: n.Y})
		idMap[n.ID] = nid
	}
	for _, e := range mf.Edges {
		from, okF := idMap[e.From]
		to, okT := idMap[e.To]
		if !okF || !okT {
			return fmt.Errorf("edge %d references missing node", e.ID)
		}
		net.AddEdge(from, to, Vec2{X: e.C0.X, Y: e.C0.Y}, Vec2{X: e.C1.X, Y: e.C1.Y})
	}
	return nil
}

// ExportMapFile serializes the network to a MapFile with stable id order.
func ExportMapFile(name string, net *Network) MapFile {
	mf := MapFile{Name: name}
	if net == nil {
		return mf
	}
	nodeIDs := make([]NodeID, 0, len(net.nodes))
	for id := range net.nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Slice(nodeIDs, func(i, j int) bool { return nodeIDs[i] < nodeIDs[j] })

	nodeNum := make(map[NodeID]uint32, len(nodeIDs))
	var next uint32 = 1
	for _, id := range nodeIDs {
		node := net.nodes[id]
		num := next
		next++
		nodeNum[id] = num
		mf.Nodes = append(mf.Nodes, MapFileNode{ID: num, X: node.Pos.X, Y: node.Pos.Y})
	}

	edgeIDs := make([]EdgeID, 0, len(net.edges))
	for id := range net.edges {
		edgeIDs = append(edgeIDs, id)
	}
	sort.Slice(edgeIDs, func(i, j int) bool { return edgeIDs[i] < edgeIDs[j] })

	var eNum uint32 = 1
	for _, id := range edgeIDs {
		e := net.edges[id]
		mf.Edges = append(mf.Edges, MapFileEdge{
			ID:   eNum,
			From: nodeNum[e.From],
			To:   nodeNum[e.To],
			C0:   MapFileXY{X: e.Curve.C0.X, Y: e.Curve.C0.Y},
			C1:   MapFileXY{X: e.Curve.C1.X, Y: e.Curve.C1.Y},
		})
		eNum++
	}
	return mf
}

// MarshalMapFile returns pretty JSON for a map.
func MarshalMapFile(mf MapFile) ([]byte, error) {
	return json.MarshalIndent(mf, "", "  ")
}
