package graph

import (
	"container/heap"
)

// Returns the shortest path from the vertex with key startKey to the vertex with key endKey as a string slice, and if such a path exists at all, using a function to calculate an estimated distance from a vertex to the endVertex. The heuristic function is passed the keys of a vertex and the end vertex. This function uses the A* search algorithm.
// If startKey or endKey (or both) are invalid, path will be empty and exists will be false.
func (g *Graph) ShortestPathWithHeuristic(startKey, endKey string, heuristic func(key, endKey string) int) (path []string, exists bool) {
	g.RLock()
	defer g.RUnlock()

	// start and end vertex
	start := g.get(startKey)
	end := g.get(endKey)

	// check startKey and endKey for validity
	if start == nil || end == nil {
		return
	}

	// priorityQueue for vertexes that have not yet been visited (open vertexes)
	openQueue := &priorityQueue{}

	// list containing vertexes that have not yet been visited (open vertexes)
	openList := map[*Vertex]*Item{}

	// list containing vertexes that have been visited already (closed vertexes)
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

		for neighbor, weight := range current.GetNeighbors() {
			if _, ok := closedList[neighbor]; ok {
				continue
			}

			distanceToNeighbor := distance + weight

			// skip neighbors that already have a better path leading to them
			if md, ok := openList[neighbor]; ok {
				if md.distanceFromStart < distanceToNeighbor {
					continue
				} else {
					heap.Remove(openQueue, md.index)
				}
			}

			item := &Item{
				neighbor,
				current,
				distanceToNeighbor,
				distanceToNeighbor + heuristic(neighbor.key, endKey), // estimate (= priority)
				0,
			}

			// add neighbor vertex to list of open vertexes
			openList[neighbor] = item

			// push into priority queue
			heap.Push(openQueue, item)
		}
	}

	return
}
