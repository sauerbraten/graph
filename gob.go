package graph

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type graphGob struct {
	inv      map[*Vertex]string
	Vertexes map[string]interface{}
	Edges    map[[2]string]int
}

// adds a key→vertex pair to the graphGob
func (g graphGob) add(key string, v *Vertex) {
	// set the vertex key - value pair

	g.Vertexes[key] = v.value

	// for each neighbor...
	for neighbor, edge := range v.edges {
		// check if it already exists in the vertexes map
		if _, ok := g.Vertexes[g.inv[neighbor]]; !ok {
			// if not, recursively add it before proceeding
			g.add(g.inv[neighbor], neighbor)
		}

		// save the edge connection to the neighbor into the edges map

		endpoints := [2]string{g.inv[edge.vertexes[0]], g.inv[edge.vertexes[1]]}

		g.Edges[endpoints] = edge.value
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

	// make edges and vertexGobs map

	gGob := graphGob{inv, map[string]interface{}{}, map[[2]string]int{}}

	// add vertexes and edges to gGob
	for key, v := range g.vertexes {
		gGob.add(key, v)
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

	for endpoints, value := range gGob.Edges {
		if ok := g.Connect(endpoints[0], endpoints[1], value); !ok {
			return errors.New("invalid edge endpoints")
		}
	}

	return
}
