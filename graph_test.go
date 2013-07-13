package graph

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	g := New()

	// set some vertexes
	g.Add("1")
	g.Add("2")
	g.Add("3")
	g.Add("4")

	// make some connections
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

	// test connections
	ok, weight := g.Adjacent("1", "2")
	if !ok || weight != 5 {
		t.Fail()
	}

	ok, weight = g.Adjacent("1", "3")
	if !ok || weight != 1 {
		t.Fail()
	}

	ok, weight = g.Adjacent("2", "3")
	if !ok || weight != 9 {
		t.Fail()
	}

	ok, weight = g.Adjacent("4", "2")
	if !ok || weight != 3 {
		t.Fail()
	}

	// test non-connections
	ok, _ = g.Adjacent("1", "4")
	if ok {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	g := New()

	// set some vertexes
	g.Add("1")
	g.Add("2")
	g.Add("3")
	g.Add("4")

	// make some connections
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

	// preserve a pointer to node "1"
	one := g.get("1")
	if one == nil {
		t.Fail()
	}

	// delete node
	ok = g.Delete("1")
	if !ok {
		t.Fail()
	}

	// make sure it's not in the graph anymore
	deletedOne := g.get("1")
	if deletedOne != nil {
		t.Fail()
	}

	// test for orphaned connections
	neighbors := g.get("2").GetNeighbors()
	for n, _ := range neighbors {
		if n == one {
			t.Fail()
		}
	}

	neighbors = g.get("3").GetNeighbors()
	for n, _ := range neighbors {
		if n == one {
			t.Fail()
		}
	}
}

func TestGob(t *testing.T) {
	g := New()

	// set key → value pairs
	g.Add("1")
	g.Add("2")
	g.Add("3")
	g.Add("4")

	// connect vertexes/nodes
	g.Connect("1", "2", 5)
	g.Connect("1", "3", 1)
	g.Connect("2", "3", 9)
	g.Connect("4", "2", 3)

	// encode
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

	// validate length of new graph
	if len(g.vertexes) != len(newG.vertexes) {
		t.Fail()
	}

	// validate contents of new graph
	for k, v := range g.vertexes {
		if newV := newG.get(k); newV.key != v.key {
			t.Fail()
		}
	}
}

func ExampleGraph() {
	g := New()

	// set key → value pairs
	g.Add("1")
	g.Add("2")
	g.Add("3")
	g.Add("4")

	// connect vertexes/nodes
	g.Connect("1", "2", 5)
	g.Connect("1", "3", 1)
	g.Connect("2", "3", 9)
	g.Connect("4", "2", 3)

	// delete a node, and all connections to it
	g.Delete("1")

	// encode into buffer
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
}

func printVertexes(vSlice map[string]*Vertex) {
	for _, v := range vSlice {
		fmt.Printf("%v\n", v.key)
		for otherV, _ := range v.neighbors {
			fmt.Printf("  → %v\n", otherV.key)
		}
	}
}
