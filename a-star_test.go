package graph

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var (
	// simple heurisitc function – the heuristic function used here returns the absolute difference between the two ints as a simple guessing technique
	h func(string, string) int = func(key, otherKey string) int {
		diff := m[key] - m[otherKey]

		if diff < 0 {
			diff = -diff
		}

		return diff
	}

	m map[string]int = map[string]int{
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
)

func TestShortestPathWithHeuristic(t *testing.T) {
	g := New()

	// add nodes
	for key := range m {
		g.Add(key)
	}

	// connect nodes
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

	_, ok := g.ShortestPathWithHeuristic("1", "9", h)
	if !ok {
		t.Fail()
	}

	// test impossible path

	g = New()

	// add nodes
	for key := range m {
		g.Add(key)
	}

	// connect nodes
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

	_, ok = g.ShortestPathWithHeuristic("1", "9", h)
	if ok {
		t.Fail()
	}
}

func ExampleGraph_ShortestPathWithHeuristic() {
	g := New()

	// add nodes
	for key := range m {
		g.Add(key)
	}

	// connect nodes
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
	path, ok := g.ShortestPathWithHeuristic("1", "9", h)
	if !ok {
		fmt.Println("something went wrong")
	}

	for _, key := range path {
		fmt.Print(key, " ")
	}
	fmt.Println()

	// Output:
	// 9 8 6 5 4 3 1
}

func TestCreateGrid_100X100Nodes(t *testing.T) {
	expectedNodesConnections := []struct {
		connectedNode       string
		expectedConnections []string
	}{
		{
			connectedNode: "1",
			expectedConnections: []string{
				"2",     // right
				"102",   // right-down
				"101",   // down
				"200",   // left-down
				"100",   // left
				"10000", // left-up
				"9901",  // up
				"9902",  // up-right
			},
		},
		{
			connectedNode: "245",
			expectedConnections: []string{
				"246", // right
				"346", // right-down
				"345", // down
				"344", // left-down
				"244", // left
				"144", // left-up
				"145", // up
				"146", // up-right
			},
		},
		{
			connectedNode: "10000",
			expectedConnections: []string{
				"9901", // right
				"1",    // right-down
				"100",  // down
				"99",   // left-down
				"9999", // left
				"9899", // left-up
				"9900", // up
				"9901", // up-right
			},
		},
	}

	g := createGrid(100, 100)
	for _, e := range expectedNodesConnections {
		for _, expectedConn := range e.expectedConnections {
			connExist, _ := g.Adjacent(e.connectedNode, expectedConn)
			if !connExist {
				t.Fatalf("Expected grid to have a connection between node %s and node %s", e.connectedNode, expectedConn)
			}
		}
	}
}

func createGrid(rows, columns int) *Graph {
	g := New()

	// add nodes
	totalNumberOfNodes := rows * columns
	for i := 1; i <= totalNumberOfNodes; i++ {
		g.Add(strconv.Itoa(i))
	}

	// connect each node to 8 surrounding nodes
	for i := 1; i <= totalNumberOfNodes; i++ {
		nodeKey := strconv.Itoa(i)

		right := getRight(i, columns)
		left := getLeft(i, columns)
		down := getDown(i, rows, columns)

		rightDown := getDown(right, rows, columns)
		downLeft := getDown(left, rows, columns)

		g.Connect(nodeKey, strconv.Itoa(right), rand.Intn(5))
		g.Connect(nodeKey, strconv.Itoa(rightDown), rand.Intn(5))
		g.Connect(nodeKey, strconv.Itoa(down), rand.Intn(5))
		g.Connect(nodeKey, strconv.Itoa(downLeft), rand.Intn(5))
	}

	return g
}

func getRight(nodePos, columns int) int {
	if nodePos%columns == 0 {
		return nodePos - (columns - 1)
	}
	return nodePos + 1
}

func getLeft(nodePos, columns int) int {
	if nodePos%columns == 1 {
		return nodePos + (columns - 1)
	}
	return nodePos - 1
}

func getDown(nodePos, rows, columns int) int {
	FirstNodeOfLastRow := rows*columns - columns + 1
	if nodePos >= FirstNodeOfLastRow {
		return nodePos - (rows*columns - columns)
	}
	return nodePos + columns
}

func BenchmarkShortestPathWithHeuristic_100X100GridOfNodes(b *testing.B) {
	rows, columns := 100, 100
	grid := createGrid(rows, columns)
	randomNodes := createLongListOfRandomKeyNodes(rows*columns, 10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		j := i % len(randomNodes)
		_, ok := grid.ShortestPathWithHeuristic(randomNodes[j].fromNode, randomNodes[j].toNode, heuristicFor100X100Grid)
		if !ok {
			b.Fatal(`ok is false`, randomNodes[j])
		}
	}
}

type randFromTo struct {
	fromNode, toNode string
}

func createLongListOfRandomKeyNodes(totalNumberOfNodes, resultLength int) (res []randFromTo) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < resultLength; i++ {
		randRes := randFromTo{"0", "0"}
		for randRes.fromNode == "0" || randRes.toNode == "0" {
			randRes = randFromTo{
				fromNode: strconv.Itoa(rand.Intn(totalNumberOfNodes)),
				toNode:   strconv.Itoa(rand.Intn(totalNumberOfNodes)),
			}
		}

		res = append(res, randRes)
	}

	return res
}

func TestHeuristicCalculationFor100X100Grid(t *testing.T) {
	expected := []struct {
		startNode, endNode string
		distance           int
	}{
		{
			startNode: "0",
			endNode:   "1",
			distance:  1,
		},
		{
			startNode: "1",
			endNode:   "3",
			distance:  2,
		},
		{
			startNode: "1",
			endNode:   "101",
			distance:  1,
		},
		{
			startNode: "220",
			endNode:   "450",
			distance:  30,
		},
		{
			startNode: "50",
			endNode:   "850",
			distance:  8,
		},
		{
			startNode: "950",
			endNode:   "102",
			distance:  48,
		},
		{
			startNode: "940",
			endNode:   "7172",
			distance:  38, // 32 to the right + 62 down (38 up)
		},
		{
			startNode: "7172",
			endNode:   "940",
			distance:  38, // 32 to the left + 62 up (38 down)
		},
		{
			startNode: "7178",
			endNode:   "7920",
			distance:  42, // 58 to the left >
		},
	}

	for _, e := range expected {
		res := heuristicFor100X100Grid(e.startNode, e.endNode)
		if res != e.distance {
			t.Fatalf("Distance between startNode %v and endNode %v should have been %v but was %v",
				e.startNode, e.endNode, e.distance, res)
		}
	}
}

func heuristicFor100X100Grid(startKey, endKey string) int {
	rows, cols := 100, 100

	startPos, _ := strconv.Atoi(startKey)
	endPos, _ := strconv.Atoi(endKey)

	rowsDiff := calculateDiff((startPos / rows), (endPos / rows), cols)
	colsDiff := calculateDiff((startPos % cols), (endPos % cols), cols)

	if rowsDiff > colsDiff {
		return rowsDiff
	}
	return colsDiff
}

func calculateDiff(a, b, size int) int {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	if diff > size/2 {
		diff = size - diff
	}
	return diff
}
