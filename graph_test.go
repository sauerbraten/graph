package graph

import (
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

	g.Connect("1", "3", 1)
	if !ok {
		t.Fail()
	}

	g.Connect("2", "3", 9)
	if !ok {
		t.Fail()
	}

	g.Connect("4", "2", 3)
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

	g.Set("1", 123)
	g.Set("2", 678)
	g.Set("3", "abc")
	g.Set("4", "xyz")

	g.Connect("1", "2", 5)
	g.Connect("1", "3", 1)
	g.Connect("2", "3", 9)
	g.Connect("4", "2", 3)

	printVertexes(g.Vertexes)
	fmt.Println(" - - - - - - ")

	g.Delete("1")

	printVertexes(g.Vertexes)
}

func printVertexes(vSlice map[string]*Vertex) {
	for _, v := range vSlice {
		fmt.Printf("%v\n", v.value)
		for otherV, _ := range v.edges {
			fmt.Printf("  → %v\n", otherV.value)
		}
	}
}
