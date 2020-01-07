package graph

import (
	"fmt"
	"strconv"
	"testing"
)

var (
	// simple heurisitc function – the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	h = func(key, otherKey string) int {
		m := map[string]int{
			"0": 0,
			"1": 1,
			"2": 2,
			"3": 3,
			"4": 4,
			"5": 5,
			"6": 6,
			"7": 7,
			"8": 8,
			"9": 9,
		}

		diff := m[key] - m[otherKey]

		if diff < 0 {
			diff = -diff
		}

		return diff
	}
)

func TestShortestPathWithHeuristic(t *testing.T) {
	g := New()

	// add nodes
	for i := range [10]int{} {
		g.Add(strconv.Itoa(i))
	}

	// connect nodes
	g.Connect("0", "1", 1)
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

	_, err := g.ShortestPathWithHeuristic("0", "9", h)
	if err != nil {
		t.Error(err)
	}

	// test impossible path

	g = New()

	// add nodes
	for i := range [10]int{} {
		g.Add(strconv.Itoa(i))
	}

	// connect nodes
	g.Connect("0", "1", 1)
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

	path, err := g.ShortestPathWithHeuristic("0", "9", h)
	if err != ErrNoPath {
		t.Errorf("expected to find no path, but found %v", path)
	}
}

func ExampleShortestPathWithHeuristic() {
	g := New()

	// add nodes
	for i := range [10]int{} {
		g.Add(strconv.Itoa(i))
	}

	// connect nodes
	g.Connect("0", "1", 1)
	g.Connect("1", "2", 1)
	g.Connect("1", "3", 2) // these two lines make it cheaper to go 1→3
	g.Connect("2", "3", 2) // than 1→2→3
	g.Connect("3", "4", 1)
	g.Connect("4", "5", 1)
	g.Connect("5", "6", 1)
	g.Connect("6", "7", 1)
	g.Connect("6", "8", 2) // these two lines make it cheaper to go 6→8
	g.Connect("7", "8", 2) // than 6→7→8
	g.Connect("8", "9", 1)

	// the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	path, err := g.ShortestPathWithHeuristic("0", "9", h)
	if err != nil {
		fmt.Println("something went wrong:", err)
	}

	for _, key := range path {
		fmt.Print(key, " ")
	}
	fmt.Println()

	// Output:
	// 9 8 6 5 4 3 1 0
}
