// Package graph implements a weighted, undirected graph data structure.
// See https://en.wikipedia.org/wiki/Graph_(abstract_data_type) for more information.
package graph

import (
	"errors"
	"sync"
)

type Vertex struct {
	key       string
	neighbors map[*Vertex]int // maps the neighbor node to the weight of the connection to it
	sync.RWMutex
}

// Returns the map of neighbors.
func (v *Vertex) GetNeighbors() map[*Vertex]int {
	if v == nil {
		return nil
	}

	v.RLock()
	neighbors := v.neighbors
	v.RUnlock()

	return neighbors
}

// Returns the vertexes key.
func (v *Vertex) Key() string {
	if v == nil {
		return ""
	}

	v.RLock()
	key := v.key
	v.RUnlock()

	return key
}

type Graph struct {
	vertexes map[string]*Vertex // A map of all the vertexes in this graph, indexed by their key.
	sync.RWMutex
}

// Initializes a new graph.
func New() *Graph {
	return &Graph{map[string]*Vertex{}, sync.RWMutex{}}
}

// Returns the amount of vertexes contained in the graph.
func (g *Graph) Len() int {
	return len(g.vertexes)
}

// Creates a new vertex. Returns true if the vertex was created, false if the key is already in use.
func (g *Graph) Add(key string) bool {
	// lock graph until this method is finished to prevent changes made by other goroutines
	g.Lock()
	defer g.Unlock()

	if g.get(key) != nil {
		return false
	}

	// create new vertex and add it to the graph
	g.vertexes[key] = &Vertex{key, map[*Vertex]int{}, sync.RWMutex{}}

	return true
}

// Deletes the vertex with the specified key. Returns false if key is invalid.
func (g *Graph) Delete(key string) bool {
	// lock graph until this method is finished to prevent changes made by other goroutines while this one is looping etc.
	g.Lock()
	defer g.Unlock()

	// get vertex in question
	v := g.get(key)
	if v == nil {
		return false
	}

	// iterate over neighbors, remove edges from neighboring vertexes
	for neighbor, _ := range v.neighbors {
		// delete edge to the to-be-deleted vertex
		neighbor.Lock()
		delete(neighbor.neighbors, v)
		neighbor.Unlock()
	}

	// delete vertex
	delete(g.vertexes, key)

	return true
}

// Returns a slice containing all vertexes. The slice is empty if the graph contains no nodes.
func (g *Graph) GetAll() (all []*Vertex) {
	g.RLock()
	for _, v := range g.vertexes {
		all = append(all, v)
	}
	g.RUnlock()

	return
}

// Returns the vertex with this key, or nil and an error if there is no vertex with this key.
func (g *Graph) Get(key string) (v *Vertex, err error) {
	g.RLock()
	v = g.get(key)
	g.RUnlock()

	if v == nil {
		err = errors.New("graph: invalid key")
	}

	return
}

// Internal function, does NOT lock the graph, should only be used in between RLock() and RUnlock() (or Lock() and Unlock()).
func (g *Graph) get(key string) *Vertex {
	return g.vertexes[key]
}

// Creates an edge between the vertexes specified by the keys. Returns false if one or both of the keys are invalid or if they are the same.
// If there already is a connection, it is overwritten with the new edge weight.
func (g *Graph) Connect(key string, otherKey string, weight int) bool {
	// recursive edges are forbidden
	if key == otherKey {
		return false
	}

	// lock graph for reading until this method is finished to prevent changes made by other goroutines while this one is running
	g.RLock()
	defer g.RUnlock()

	// get vertexes and check for validity of keys
	v := g.get(key)
	otherV := g.get(otherKey)

	if v == nil || otherV == nil {
		return false
	}

	// add connection to both vertexes
	v.Lock()
	otherV.Lock()

	v.neighbors[otherV] = weight
	otherV.neighbors[v] = weight

	v.Unlock()
	otherV.Unlock()

	// success
	return true
}

// Removes an edge connecting the two vertexes. Returns false if one or both of the keys are invalid or if they are the same.
func (g *Graph) Disconnect(key string, otherKey string) bool {
	// recursive edges are forbidden
	if key == otherKey {
		return false
	}

	// lock graph for reading until this method is finished to prevent changes made by other goroutines while this one is running
	g.RLock()
	defer g.RUnlock()

	// get vertexes and check for validity of keys
	v := g.get(key)
	otherV := g.get(otherKey)

	if v == nil || otherV == nil {
		return false
	}

	// delete the edge from both vertexes
	v.Lock()
	otherV.Lock()

	delete(v.neighbors, otherV)
	delete(otherV.neighbors, v)

	v.Unlock()
	otherV.Unlock()

	return true
}

// Returns true and the edge weight if there is an edge between the vertexes specified by their keys. Returns false if one or both keys are invalid, if they are the same, or if there is no edge between the vertexes.
func (g *Graph) Adjacent(key string, otherKey string) (exists bool, weight int) {
	// sanity check
	if key == otherKey {
		return
	}

	g.RLock()

	v := g.get(key)
	if v == nil {
		g.RUnlock()
		return
	}

	otherV := g.get(otherKey)
	if otherV == nil {
		g.RUnlock()
		return
	}

	g.RUnlock()

	v.RLock()
	defer v.RUnlock()
	otherV.RUnlock()
	defer otherV.RLock()

	// choose vertex with less edges (easier to find 1 in 10 than to find 1 in 100)
	if len(v.neighbors) < len(otherV.neighbors) {
		// iterate over it's map of edges; when the right vertex is found, return
		for iteratingV, weight := range v.neighbors {
			if iteratingV == otherV {
				return true, weight
			}
		}
	} else {
		// iterate over it's map of edges; when the right vertex is found, return
		for iteratingV, weight := range otherV.neighbors {
			if iteratingV == v {
				return true, weight
			}
		}
	}

	return
}
