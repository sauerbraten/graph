package graph

import (
	"container/heap"
)

// Returns the shortest path from the node with key startKey to the node with key endKey as a string slice, and if such a path exists at all, using a function to calculate an estimated distance from a node to the endNode. The heuristic function is passed the keys of a node and the end node. This function uses the A* search algorithm.
// If startKey or endKey (or both) are invalid, path will be empty and exists will be false.
func (g *Graph) ShortestPathWithHeuristic(startKey, endKey string, heuristic func(key, endKey string) int) (path []string, exists bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	// start and end node
	start := g.get(startKey)
	end := g.get(endKey)

	// check startKey and endKey for validity
	if start == nil || end == nil {
		return
	}

	// priorityQueue for nodes that have not yet been visited (open nodes)
	openQueue := &priorityQueue{}

	// list containing nodes that have not yet been visited (open nodes)
	openList := map[*Node]*Item{}

	// list containing nodes that have been visited already (closed nodes)
	closedList := map[*Node]*Item{}

	// add start node to list of open nodes
	item := &Item{start, nil, 0, 0, 0}
	openList[start] = item

	heap.Push(openQueue, item)

	for openQueue.Len() > 0 {
		current := heap.Pop(openQueue).(*Item).n

		// current node was now visited; add to closed list
		closedList[current] = openList[current]
		delete(openList, current)

		// end node found?
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

			// add neighbor node to list of open nodes
			openList[neighbor] = item

			// push into priority queue
			heap.Push(openQueue, item)
		}
	}

	return
}
