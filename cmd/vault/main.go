package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type coordinate []int

type node struct {
	Coordinates coordinate
	Value       int
	Operation   string
	Neighbors   [][]int
}

func (n node) validNeighbors() bool {
	for _, c := range n.Neighbors {
		cx, cy := n.Coordinates[0], n.Coordinates[1]
		nx, ny := c[0], c[1]

		absDistance := math.Abs(float64(cx-nx)) + math.Abs(float64(cy-ny))
		if absDistance > 1 {
			return false
		}
	}
	return true
}

type tree map[int]map[int]node

func newTree() tree {
	t := make(tree)
	for i := 0; i < 5; i++ {
		t[i] = make(map[int]node)
	}
	return t
}

func (t tree) add(n node) {
	x, y := n.Coordinates[0], n.Coordinates[1]
	t[x][y] = n
}

func (t tree) load(file string) {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(f)

	_, err = dec.Token()
	if err != nil {
		panic(err)
	}

	for dec.More() {
		var n node
		err := dec.Decode(&n)
		if err != nil {
			panic(err)
		}

		if n.validNeighbors() {
			t.add(n)
		} else {
			panic(fmt.Sprintf("Invalid neighbors: %v", n.Neighbors))
		}
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func (t tree) root() node {
	return t[0][0]
}

func add(x, y int) int {
	return x + y
}

func div(x, y int) int {
	return x / y
}

func mul(x, y int) int {
	return x * y
}

func nop(x, y int) int {
	return x
}

func sub(x, y int) int {
	return x - y
}

func main() {
	// create tree from file
	t := newTree()
	t.load("cmd/vault/vault.json")

	// initiate vars
	stack := []coordinate{t.root().Coordinates}
	path := []coordinate{}
	orb := t.root().Value

	// track visited nodes
	visited := make(map[int]map[int]bool)
	for i := 0; i < 5; i++ {
		visited[i] = make(map[int]bool)
	}

	// iterate through paths
	for len(stack) > 0 {
		// dequeue / shrinks
		c := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		fmt.Println("next node from queue:", c)

		// skip if visited
		if visited[c[0]][c[1]] == true {
			continue
		} else {
			visited[c[0]][c[1]] = true
		}

		// path grows
		n := t[c[0]][c[1]] // current node
		path = append(path, n.Coordinates)
		fmt.Println("current path:", path)

		if len(n.Neighbors) == 0 {
			fmt.Println("reached a reset room")
			fmt.Println("  orb value:", orb)
			fmt.Println("  path:", path)

			path = []coordinate{}
			stack = []coordinate{stack[0]}
			orb = t.root().Value
			for i := 0; i < 5; i++ {
				visited[i] = make(map[int]bool)
			}
			continue
		}

		// add neighbors to stack
		for _, neighbor := range n.Neighbors {
			// only visit neighbors not seen yet
			if visited[neighbor[0]][neighbor[1]] == false {
				stack = append(stack, neighbor)
			}
		}

		// not including the root node, path can't exceed 12 spaces
		if len(path) > 13 {
			fmt.Println("Path is too long, removing last item")
			// remove last 2 nodes (current node and the previous)
			path = path[0 : len(path)-2]
			// unvisit last node
			visited[3][3] = false
		}
	}
}
