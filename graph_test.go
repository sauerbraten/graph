package graph

import (
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	g := New()

	// set some nodes
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

	// set some nodes
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

func ExampleGraph() {
	g := New()

	// set key → value pairs
	g.Add("1")
	g.Add("2")
	g.Add("3")
	g.Add("4")

	// connect nodes
	g.Connect("1", "2", 5)
	g.Connect("1", "3", 1)
	g.Connect("2", "3", 9)
	g.Connect("4", "2", 3)

	// delete a node, and all connections to it
	g.Delete("1")
}

func printNodes(nodes map[string]*Node) {
	for _, n := range nodes {
		fmt.Printf("%v\n", n.key)
		for neighbor, _ := range n.neighbors {
			fmt.Printf("  → %v\n", neighbor.key)
		}
	}
}
