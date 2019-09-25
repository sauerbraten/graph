package graph

import "sync"

type Node struct {
	key       string
	neighbors map[*Node]int // maps the neighbor node to the weight of the connection to it
	sync.RWMutex
}

// Returns the map of neighbors.
func (n *Node) GetNeighbors() map[*Node]int {
	if n == nil {
		return nil
	}

	n.RLock()
	neighbors := n.neighbors
	n.RUnlock()

	return neighbors
}

// Returns the Nodees key.
func (n *Node) Key() string {
	if n == nil {
		return ""
	}

	n.RLock()
	key := n.key
	n.RUnlock()

	return key
}
