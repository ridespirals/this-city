package pathfind

import (
	"container/heap"

	"github.com/ridespirals/this-city/internal/sim"
)

// heapItem is an open-set entry for A*/Dijkstra.
type heapItem struct {
	node sim.NodeID
	f    float32 // priority (g+h)
	g    float32
}

type nodeHeap []heapItem

func (h nodeHeap) Len() int { return len(h) }
func (h nodeHeap) Less(i, j int) bool {
	return h[i].f < h[j].f || (h[i].f == h[j].f && h[i].g < h[j].g)
}
func (h nodeHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *nodeHeap) Push(x any) { *h = append(*h, x.(heapItem)) }

func (h *nodeHeap) Pop() any {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

// update decreases an existing node's priority if present (linear scan; fine for city graphs).
func (h *nodeHeap) update(node sim.NodeID, f, g float32) {
	for i := range *h {
		if (*h)[i].node == node {
			if f < (*h)[i].f {
				(*h)[i].f = f
				(*h)[i].g = g
				heap.Fix(h, i)
			}
			return
		}
	}
}
