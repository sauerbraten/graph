package graph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

func TestGraph(t *testing.T) {
	g := New()

	g.Set("1", 123)
	g.Set("2", 678)
	g.Set("3", "abc")
	g.Set("4", "xyz")

	ok := g.Connect("1", "2", 5)
	if !ok {
		t.Fail()
	}

	ok = g.Connect("1", "3", 1)
	if !ok {
		t.Fail()
	}

	ok = g.Connect("2", "3", 9)
	if !ok {
		t.Fail()
	}

	ok = g.Connect("4", "2", 3)
	if !ok {
		t.Fail()
	}

	ok = g.Delete("1")
	if !ok {
		t.Fail()
	}
}

func ExampleGraph() {
	g := New()

	// set key → value pairs
	g.Set("1", 123)
	g.Set("2", 678)
	g.Set("3", "abc")
	g.Set("4", "xyz")

	// connect vertexes/nodes
	g.Connect("1", "2", 5)
	g.Connect("1", "3", 1)
	g.Connect("2", "3", 9)
	g.Connect("4", "2", 3)

	printVertexes(g.vertexes)
	fmt.Println(" - - - - - - ")

	// delete a node, and all connections to it
	g.Delete("1")

	printVertexes(g.vertexes)
	fmt.Println()

	// test gob encoding and decoding
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	err := enc.Encode(g)
	if err != nil {
		fmt.Println(err)
	}

	// now decode into new graph
	dec := gob.NewDecoder(buf)
	newG := New()
	err = dec.Decode(newG)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("new graph:")
	printVertexes(newG.vertexes)

	// output differs everytime you run this example, because of the map being used for vertexes (maps don't have a fixed order)
}

func printVertexes(vSlice map[string]*Vertex) {
	for _, v := range vSlice {
		fmt.Printf("%v\n", v.value)
		for otherV, _ := range v.edges {
			fmt.Printf("  → %v\n", otherV.value)
		}
	}
}
