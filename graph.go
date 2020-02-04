// Package graph implements a weighted, undirected graph data structure.
// See https://en.wikipedia.org/wiki/Graph_(abstract_data_type) for more information.
package graph

import (
	"errors"
	"sync"
)

type Graph struct {
	mu    sync.RWMutex
	nodes map[string]*Node // A map of all the nodes in this graph, indexed by their key.
}

// Initializes a new graph.
func New() *Graph {
	return &Graph{nodes: map[string]*Node{}}
}

// Returns the amount of nodes contained in the graph.
func (g *Graph) Len() int {
	return len(g.nodes)
}

// Creates a new node. Returns true if the node was created, false if the key is already in use.
func (g *Graph) Add(key string) bool {
	// lock graph until this method is finished to prevent changes made by other goroutines
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.get(key) != nil {
		return false
	}

	// create new node and add it to the graph
	g.nodes[key] = &Node{key, map[*Node]int{}, sync.RWMutex{}}

	return true
}

// Deletes the node with the specified key. Returns false if key is invalid.
func (g *Graph) Delete(key string) bool {
	// lock graph until this method is finished to prevent changes made by other goroutines while this one is looping etc.
	g.mu.Lock()
	defer g.mu.Unlock()

	// get node in question
	n := g.get(key)
	if n == nil {
		return false
	}

	// iterate over neighbors, remove edges from neighboring nodes
	for neighbor := range n.neighbors {
		// delete edge to the to-be-deleted node
		neighbor.Lock()
		delete(neighbor.neighbors, n)
		neighbor.Unlock()
	}

	// delete node
	delete(g.nodes, key)

	return true
}

// Returns a slice containing all nodes. The slice is empty if the graph contains no nodes.
func (g *Graph) GetAll() (all []*Node) {
	g.mu.RLock()
	for _, n := range g.nodes {
		all = append(all, n)
	}
	g.mu.RUnlock()

	return
}

// Returns the node with this key, or nil and an error if there is no node with this key.
func (g *Graph) Get(key string) (n *Node, err error) {
	g.mu.RLock()
	n = g.get(key)
	g.mu.RUnlock()

	if n == nil {
		err = errors.New("graph: invalid key")
	}

	return
}

// Internal function, does NOT lock the graph, should only be used in between RLock() and RUnlock() (or Lock() and Unlock()).
func (g *Graph) get(key string) *Node {
	return g.nodes[key]
}

// Creates an edge between the nodes specified by the keys. Returns false if one or both of the keys are invalid or if they are the same.
// If there already is a connection, it is overwritten with the new edge weight.
func (g *Graph) Connect(keyA string, keyB string, weight int) bool {
	// reflective edges are forbidden
	if keyA == keyB {
		return false
	}

	// lock graph for reading until this method is finished to prevent changes made by other goroutines while this one is running
	g.mu.RLock()
	defer g.mu.RUnlock()

	// get nodes and check for validity of keys
	nodeA := g.get(keyA)
	nodeB := g.get(keyB)

	if nodeA == nil || nodeB == nil {
		return false
	}

	// add connection to both nodes
	nodeA.Lock()
	nodeB.Lock()

	nodeA.neighbors[nodeB] = weight
	nodeB.neighbors[nodeA] = weight

	nodeA.Unlock()
	nodeB.Unlock()

	// success
	return true
}

// Removes an edge connecting the two nodes. Returns false if one or both of the keys are invalid or if they are the same.
func (g *Graph) Disconnect(keyA string, keyB string) bool {
	// recursive edges are forbidden
	if keyA == keyB {
		return false
	}

	// lock graph for reading until this method is finished to prevent changes made by other goroutines while this one is running
	g.mu.RLock()
	defer g.mu.RUnlock()

	// get nodes and check for validity of keys
	nodeA := g.get(keyA)
	nodeB := g.get(keyB)

	if nodeA == nil || nodeB == nil {
		return false
	}

	// delete the edge from both nodes
	nodeA.Lock()
	nodeB.Lock()

	delete(nodeA.neighbors, nodeB)
	delete(nodeB.neighbors, nodeA)

	nodeA.Unlock()
	nodeB.Unlock()

	return true
}

// Returns true and the edge weight if there is an edge between the nodes specified by their keys. Returns false and 0 if one or both keys are invalid, if they are the same, or if there is no edge between the nodes.
func (g *Graph) Adjacent(keyA string, keyB string) (exists bool, weight int) {
	// sanity check
	if keyA == keyB {
		return
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	nodeA := g.get(keyA)
	if nodeA == nil {
		return
	}

	nodeB := g.get(keyB)
	if nodeB == nil {
		return
	}

	nodeA.RLock()
	defer nodeA.RUnlock()
	nodeB.RLock()
	defer nodeB.RUnlock()

	// choose node with less edges (easier to find 1 in 10 than to find 1 in 100)
	if len(nodeA.neighbors) < len(nodeB.neighbors) {
		// iterate over it's map of edges; when the right node is found, return
		for neighbor, weight := range nodeA.neighbors {
			if neighbor == nodeB {
				return true, weight
			}
		}
	} else {
		// iterate over it's map of edges; when the right node is found, return
		for neighbor, weight := range nodeB.neighbors {
			if neighbor == nodeA {
				return true, weight
			}
		}
	}

	return
}
