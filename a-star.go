package graph

// meta data of vertexes in the closed list
type metaData struct {
	prev              *Vertex
	estimate          int
	distanceFromStart int
}

// Work in progress, do not use! See line 67 in a-star.go! Returns the shortest path from the vertex with key startKey to the vertex with key endKey as a string slice, and if such a path exists at all, using a function to calculate an estimated distance from a vertex to the endVertex. The heuristic function is passed the keys of a vertex and the end vertex. This function uses the A* search algorithm.
func (g *Graph) ShortestPathWithHeuristic(startKey, endKey string, heuristic func(key, endKey string) int) (path []string, exists bool) {
	g.RLock()
	defer g.RUnlock()

	// start and end vertex
	start := g.get(startKey)
	end := g.get(endKey)

	// list for vertexes that have not yet been visited (open vertexes)
	openList := map[*Vertex]metaData{}
	// list for vertexes that have been visited already (closed vertexes)
	closedList := map[*Vertex]metaData{}

	// add start vertex to list of open vertexes
	openList[start] = metaData{nil, 0, 0}

	for current, ok := chooseNext(openList); ok; current, ok = chooseNext(openList) {
		// current vertex was now visited, move it to list of closed vertexes
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
			if md, ok := openList[n.V]; ok && md.distanceFromStart < distanceToNeighbor {
				continue
			}
			openList[n.V] = metaData{current, distance + n.EdgeWeight + heuristic(n.V.key, endKey), distance + n.EdgeWeight}
		}
	}

	return
}

// WARNING: not finished! choosing the next vertex takes O(n) time; n being len(m)!
func chooseNext(m map[*Vertex]metaData) (next *Vertex, ok bool) {
	if len(m) == 0 {
		return
	}

	bestEstimate := 1000
	ok = true

	for v, md := range m {
		if md.estimate < bestEstimate {
			next = v
		}
	}

	return
}
