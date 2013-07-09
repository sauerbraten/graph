package graph

import (
	"bytes"
	"encoding/gob"
	"errors"
)

type graphGob struct {
	Vertexes []string
	Edges    map[string]map[string]int
}

// Encodes the graph into a []byte. With this method, graph implements the gob.GobEncoder interface.
func (g *Graph) GobEncode() ([]byte, error) {
	gGob := graphGob{[]string{}, map[string]map[string]int{}}

	// add vertexes and edges to gGob
	for key, v := range g.vertexes {
		gGob.Vertexes = append(gGob.Vertexes, key)

		gGob.Edges[key] = map[string]int{}

		// for each neighbor...
		for neighbor, weight := range v.neighbors {
			// save the edge connection to the neighbor into the edges map
			gGob.Edges[key][neighbor.key] = weight
		}
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

	// add the vertexes
	for _, key := range gGob.Vertexes {
		g.Add(key)
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
