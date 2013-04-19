package graph

// An Item is something we manage in a priority queue.
type Item struct {
	v                 *Vertex // vertex this meta data belongs to
	prev              *Vertex // previous waypoint in the shortest path from start to here
	distanceFromStart int     // distance form start to this vertex using the shortest known path
	priority          int     // The priority of the item in the queue (= estimated distance from end vertex). Low value means high priority.
	index             int     // The index of the item in the heap. You do not need to set this, it's done automatically in Push(). DO NOT CHANGE!
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*Item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	item := (*pq)[len(*pq)-1]
	*pq = (*pq)[0 : len(*pq)-1]
	return item
}
