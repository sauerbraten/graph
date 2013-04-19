package graph

import (
	"container/heap"
)

// Returns the shortest path from the vertex with key startKey to the vertex with key endKey as a string slice, and if such a path exists at all, using a function to calculate an estimated distance from a vertex to the endVertex. The heuristic function is passed the keys of a vertex and the end vertex. This function uses the A* search algorithm.
func (g *Graph) ShortestPathWithHeuristic(startKey, endKey string, heuristic func(key, endKey string) int) (path []string, exists bool) {
	g.RLock()
	defer g.RUnlock()

	// start and end vertex
	start := g.get(startKey)
	end := g.get(endKey)

	// priorityQueue for vertexes that have not yet been visited (open vertexes)
	openQueue := &priorityQueue{}

	// priorityQueue for vertexes that have not yet been visited (open vertexes)
	openList := map[*Vertex]*Item{}

	// list for vertexes that have been visited already (closed vertexes)
	closedList := map[*Vertex]*Item{}

	// add start vertex to list of open vertexes
	item := &Item{start, nil, 0, 0, 0}
	openList[start] = item

	heap.Push(openQueue, item)

	for openQueue.Len() > 0 {
		current := heap.Pop(openQueue).(*Item).v

		// current vertex was now visited; add to closed list
		closedList[current] = openList[current]
		delete(openList, current)

		// end vertex found?
		if current == end {
			// path exists
			exists = true

			// build path
			for current != nil {
				path = append(path, current.key)
				current = closedList[current].prev
			}

			return
		}

		// saved here for easy usage in following loop
		distance := closedList[current].distanceFromStart

		for _, n := range current.GetNeighbors() {
			if _, ok := closedList[n.V]; ok {

				continue
			}

			distanceToNeighbor := distance + n.EdgeWeight

			// skip neighbors that already have a better path leading to them
			if md, ok := openList[n.V]; ok {
				if md.distanceFromStart < distanceToNeighbor {
					continue
				} else {
					heap.Remove(openQueue, md.index)
				}
			}

			item := &Item{
				n.V,
				current,
				distanceToNeighbor,
				distanceToNeighbor + heuristic(n.V.key, endKey), // estimate (= priority)
				0,
			}

			// add neighbor vertex to list of open vertexes
			openList[n.V] = item

			// push into priority queue
			heap.Push(openQueue, item)
		}
	}

	return
}
