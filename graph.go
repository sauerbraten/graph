// Package graph implements a graph data structure.
// See http://en.wikipedia.org/wiki/Graph_(data_structure) for more information.
package graph

import (
	"sync"
)

type edge struct {
	vertexes [2]*Vertex
	value    int
}

type Vertex struct {
	edges map[*Vertex]*edge
	value interface{}
	sync.RWMutex
}

// Returns all adjacent vertexes as a slice, which may be empty.
func (v *Vertex) GetNeighbors() []*Vertex {
	neighbors := []*Vertex{}

	for k, _ := range v.edges {
		neighbors = append(neighbors, k)
	}

	return neighbors
}

// Returns the Vertexes value.
func (v *Vertex) Value() interface{} {
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

// Sets the value of the vertex with the specified key.
func (g *Graph) Set(key string, value interface{}) {
	// lock graph until this method is finished to prevent changes made by other goroutines while this one is looping etc.
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

// Deletes the vertex with the specified key.
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
func (g *Graph) Connect(key string, otherKey string, value int) bool {
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
	e := &edge{[2]*Vertex{v, otherV}, value}

	// add it to both vertexes
	v.Lock()
	v.edges[otherV] = e
	v.Unlock()
	otherV.Lock()
	otherV.edges[v] = e
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

	v.Lock()
	delete(v.edges, otherV)
	v.Unlock()
	otherV.Lock()
	delete(otherV.edges, v)
	otherV.Unlock()

	return true
}
