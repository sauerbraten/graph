// Package graph implements a weighted, undirected graph data structure.
// See http://en.wikipedia.org/wiki/Graph_(data_structure) for more information.
package graph

import (
	"sync"
)

type edge struct {
	vertexes [2]*Vertex
	weight   int
}

type Vertex struct {
	edges map[*Vertex]*edge // maps the neighbor node to the edge connecting this node to it
	value interface{}       // the stored value
	sync.RWMutex
}

// A Neighbor consists of a neighboring vertex and an edge weight.
type Neighbor struct {
	V          *Vertex
	EdgeWeight int
}

// Returns all adjacent vertexes and the respective edge's weight as a slice, which may be empty.
func (v *Vertex) GetNeighbors() []Neighbor {
	if v == nil {
		return nil
	}

	neighbors := []Neighbor{}

	v.RLock()
	for otherV, e := range v.edges {
		neighbors = append(neighbors, Neighbor{otherV, e.weight})
	}
	v.RUnlock()

	return neighbors
}

// Returns the Vertexes value.
func (v *Vertex) Value() interface{} {
	if v == nil {
		return nil
	}

	v.RLock()
	value := v.value
	v.RUnlock()

	return value
}

type Graph struct {
	vertexes map[string]*Vertex // A map of all the vertexes in this graph, indexed by their key.
	sync.RWMutex
}

// Initializes a new graph.
func New() *Graph {
	return &Graph{make(map[string]*Vertex), sync.RWMutex{}}
}

// Returns the amount of vertexes contained in the graph.
func (g *Graph) Len() int {
	return len(g.vertexes)
}

// If there is no vertex with the specified key yet, Set creates a new vertex and stores the value. Else, Set updates the value, but leaves all connections intact.
func (g *Graph) Set(key string, value interface{}) {
	// lock graph until this method is finished to prevent changes made by other goroutines
	g.Lock()
	defer g.Unlock()

	v := g.get(key)

	// if no such node exists
	if v == nil {
		// create a new one
		v = &Vertex{make(map[*Vertex]*edge), value, sync.RWMutex{}}

		// and add it to the graph
		g.vertexes[key] = v

		return
	}

	// else, just update the value
	v.Lock()
	v.value = value
	v.Unlock()
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

	// iterate over edges, remove edges from neighboring vertexes
	for _, e := range v.edges {
		ends := e.vertexes

		// choose other node, not v
		otherV := ends[0]
		if v == ends[0] {
			otherV = ends[1]
		}

		// delete edge to the to-be-deleted vertex
		otherV.Lock()
		delete(otherV.edges, v)
		otherV.Unlock()
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

// Returns the vertex with this key, or nil if there is no vertex with this key.
func (g *Graph) Get(key string) *Vertex {
	g.RLock()
	v := g.get(key)
	g.RUnlock()

	return v
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
	if v == nil {
		return false
	}

	otherV := g.get(otherKey)
	if otherV == nil {
		return false
	}

	// make a new edge
	e := &edge{[2]*Vertex{v, otherV}, weight}

	// add it to both vertexes
	v.Lock()
	otherV.Lock()

	v.edges[otherV] = e
	otherV.edges[v] = e

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
	if v == nil {
		return false
	}

	otherV := g.get(otherKey)
	if otherV == nil {
		return false
	}

	// delete the edge from both vertexes
	v.Lock()
	otherV.Lock()

	delete(v.edges, otherV)
	delete(otherV.edges, v)

	v.Unlock()
	otherV.Unlock()

	return true
}

func (g *Graph) Adjacent(key string, otherKey string) (bool, int) {
	g.RLock()
	v := g.get(key)
	otherV := g.get(otherKey)
	g.RUnlock()

	v.RLock()
	otherV.RLock()

	// choose vertex with less edges (easier to find 1 in 10 than to find 1 in 100)
	if len(v.edges) < len(otherV.edges) {
		// iterate over it's map of edges; when the right vertex is found, return
		for iteratingV, e := range v.edges {
			if iteratingV == otherV {
				return true, e.weight
			}
		}
	} else {
		// iterate over it's map of edges; when the right vertex is found, return
		for iteratingV, e := range otherV.edges {
			if iteratingV == v {
				return true, e.weight
			}
		}
	}

	v.RUnlock()
	otherV.RUnlock()

	return false, 0
}
