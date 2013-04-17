package graph

import (
	"fmt"
	"testing"
)

func TestShortestPathWithHeuristic(t *testing.T) {
	g := New()

	// set key → value pairs
	g.Set("1", 1)
	g.Set("2", 2)
	g.Set("3", 3)
	g.Set("4", 4)
	g.Set("5", 5)
	g.Set("6", 6)
	g.Set("7", 7)
	g.Set("8", 8)
	g.Set("9", 9)

	// connect vertexes/nodes
	g.Connect("1", "2", 1)
	g.Connect("1", "3", 2) // these two lines make it cheaper to go 1→3
	g.Connect("2", "3", 2) // than 1→2→3
	g.Connect("3", "4", 1)
	g.Connect("4", "5", 1) // cost of 4→5→6 is the same as
	g.Connect("4", "6", 2) // going 4→6
	g.Connect("5", "6", 1)
	g.Connect("6", "7", 1)
	g.Connect("7", "8", 1)
	g.Connect("8", "9", 1)

	// the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	_, ok := g.ShortestPathWithHeuristic("1", "9", func(key, otherKey string) int {
		diff := g.Get(key).value.(int) - g.Get(key).value.(int)

		if diff < 0 {
			diff = -diff
		}

		return diff
	})

	if !ok {
		t.Fail()
	}

	// test impossible path

	g = New()

	// set key → value pairs
	g.Set("1", 1)
	g.Set("2", 2)
	g.Set("3", 3)
	g.Set("4", 4)
	g.Set("5", 5)
	g.Set("6", 6)
	g.Set("7", 7)
	g.Set("8", 8)
	g.Set("9", 9)

	// connect vertexes/nodes
	g.Connect("1", "2", 1)
	g.Connect("1", "3", 2) // these two lines make it cheaper to go 1→3
	g.Connect("2", "3", 2) // than 1→2→3
	// missing connection 3→4
	g.Connect("4", "5", 1) // cost of 4→5→6 is the same as
	g.Connect("4", "6", 2) // going 4→6
	g.Connect("5", "6", 1)
	g.Connect("6", "7", 1)
	g.Connect("7", "8", 1)
	g.Connect("8", "9", 1)

	// the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	_, ok = g.ShortestPathWithHeuristic("1", "9", func(key, otherKey string) int {
		diff := g.Get(key).value.(int) - g.Get(key).value.(int)

		if diff < 0 {
			diff = -diff
		}

		return diff
	})

	if ok {
		t.Fail()
	}
}

func ExampleShortestPathWithHeuristic() {
	g := New()

	// set key → value pairs
	g.Set("1", 1)
	g.Set("2", 2)
	g.Set("3", 3)
	g.Set("4", 4)
	g.Set("5", 5)
	g.Set("6", 6)
	g.Set("7", 7)
	g.Set("8", 8)
	g.Set("9", 9)

	// connect vertexes/nodes
	g.Connect("1", "2", 1)
	g.Connect("1", "3", 2) // these two lines make it cheaper to go 1→3
	g.Connect("2", "3", 2) // than 1→2→3
	g.Connect("3", "4", 1)
	g.Connect("4", "5", 1) // cost of 4→5→6 is the same as
	g.Connect("4", "6", 2) // going 4→6
	g.Connect("5", "6", 1)
	g.Connect("6", "7", 1)
	g.Connect("7", "8", 1)
	g.Connect("8", "9", 1)

	// the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	path, ok := g.ShortestPathWithHeuristic("1", "9", func(key, otherKey string) int {
		diff := g.Get(key).value.(int) - g.Get(key).value.(int)

		if diff < 0 {
			diff = -diff
		}

		return diff
	})

	if !ok {
		fmt.Println("something went wrong")
	}

	for _, key := range path {
		fmt.Print(key, " ")
	}
	fmt.Println()

	// Output:
	// 9 8 7 6 4 3 1
}
