package graph

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type graphGob struct {
	inv      map[*Vertex]string
	Vertexes map[string]interface{}
	Edges    map[string]map[string]int
}

// adds a key - vertex pair to the graphGob
func (g graphGob) add(v *Vertex) {
	// set the key - vertex pair
	g.Vertexes[v.key] = v.value

	g.Edges[v.key] = map[string]int{}

	// for each neighbor...
	for neighbor, weight := range v.neighbors {
		/*
			// check if it already exists in the vertexes map
			if _, ok := g.Vertexes[g.inv[neighbor]]; !ok {
				// if not, recursively add it before proceeding
				g.add(g.inv[neighbor], neighbor)
			}
		*/

		// save the edge connection to the neighbor into the edges map
		g.Edges[v.key][neighbor.key] = weight
	}
}

// Encodes the graph into a []byte. With this method, graph implements the gob.GobEncoder interface.
func (g *Graph) GobEncode() ([]byte, error) {
	// build inverted map
	inv := map[*Vertex]string{}
	for key, v := range g.vertexes {
		if _, ok := inv[v]; !ok {
			inv[v] = key
		}
	}

	gGob := graphGob{inv, map[string]interface{}{}, map[string]map[string]int{}}

	// add vertexes and edges to gGob
	for _, v := range g.vertexes {
		gGob.add(v)
	}

	// encode gGob
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(gGob)

	return buf.Bytes(), err
}

// Decodes a []byte into the graphs vertexes and edges. With this method, graph implements the gob.GobDecoder interface.
func (g *Graph) GobDecode(b []byte) (err error) {
	// decode into graphGob
	gGob := &graphGob{}
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)

	err = dec.Decode(gGob)
	if err != nil {
		return
	}

	// set the vertexes
	for key, value := range gGob.Vertexes {
		g.Set(key, value)
	}

	// connect the vertexes
	for key, neighbors := range gGob.Edges {
		for otherKey, weight := range neighbors {
			if ok := g.Connect(key, otherKey, weight); !ok {
				return errors.New("invalid edge endpoints")
			}
		}
	}

	return
}
