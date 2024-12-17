package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"math"
	"os"
)

type DirType int

const (
	N DirType = iota
	S
	E
	V
)

var INV_DIR = map[DirType]DirType{
	N: S,
	S: N,
	E: V,
	V: E,
}

type Point struct {
	x int
	y int
}

type PointWithDir struct {
	p Point
	d DirType
}

type Node struct {
	adj  PointWithDir
	cost int
}

func (a Node) CompareWith(e datastructures.ElementWithPriority) int {
	b := e.(Node)

	if a.cost < b.cost {
		return -1
	}

	if a.cost > b.cost {
		return 1
	}

	return 0
}

func readFileAsString(filePath string) ([]string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, nil
}

func validCoords(grid []string, p Point) bool {
	return p.x >= 0 && p.x < len(grid) && p.y >= 0 && p.y < len(grid[0])
}

func getAdjWithDir(grid []string, p Point) []PointWithDir {
	points := map[DirType]Point{
		N: {p.x - 1, p.y},
		S: {p.x + 1, p.y},
		V: {p.x, p.y - 1},
		E: {p.x, p.y + 1},
	}

	res := []PointWithDir{}

	for dir, point := range points {
		if validCoords(grid, point) {
			res = append(res, PointWithDir{p: point, d: dir})
		}
	}

	return res
}

func addTurningCost(graph *map[PointWithDir][]Node, p Point, from DirType, to DirType, cost int) {
	fromWithDirection := PointWithDir{p, from}
	toWithDirection := PointWithDir{p, to}

	if v, ok := (*graph)[fromWithDirection]; ok {
		if _, ok1 := (*graph)[toWithDirection]; ok1 {
			v = append(v, Node{toWithDirection, cost})
			(*graph)[fromWithDirection] = v
		}
	}
}

func makeGraph(data []string) (map[PointWithDir][]Node, PointWithDir, Point) {
	res := map[PointWithDir][]Node{}
	start := PointWithDir{Point{-1, -1}, N}
	end := Point{-1, -1}

	for idx := 0; idx < len(data); idx++ {
		for jdx := 0; jdx < len(data[0]); jdx++ {
			if data[idx][jdx] == 'S' {
				start = PointWithDir{p: Point{idx, jdx}, d: E}
			}

			if data[idx][jdx] == 'E' {
				end = Point{idx, jdx}
			}

			currPoint := Point{idx, jdx}
			neighbours := getAdjWithDir(data, currPoint)

			for _, adj := range neighbours {
				currNode := PointWithDir{p: currPoint, d: adj.d}

				if _, ok := res[currNode]; !ok {
					res[currNode] = []Node{}
				}

				if data[adj.p.x][adj.p.y] == '#' {
					continue
				}

				pWithDir := PointWithDir{p: adj.p, d: INV_DIR[adj.d]}
				res[currNode] = append(res[currNode], Node{pWithDir, 1})
			}

			addTurningCost(&res, currPoint, N, E, 1000)
			addTurningCost(&res, currPoint, N, V, 1000)

			addTurningCost(&res, currPoint, E, N, 1000)
			addTurningCost(&res, currPoint, E, S, 1000)

			addTurningCost(&res, currPoint, S, E, 1000)
			addTurningCost(&res, currPoint, S, V, 1000)

			addTurningCost(&res, currPoint, V, S, 1000)
			addTurningCost(&res, currPoint, V, N, 1000)

			addTurningCost(&res, currPoint, N, S, 0)
			addTurningCost(&res, currPoint, S, N, 0)
			addTurningCost(&res, currPoint, E, V, 0)
			addTurningCost(&res, currPoint, V, E, 0)
		}
	}

	return res, start, end
}

func dijkstra(graph map[PointWithDir][]Node, start PointWithDir) (map[PointWithDir]int, map[PointWithDir]map[PointWithDir]bool) {
	dist := map[PointWithDir]int{}
	prev := map[PointWithDir]map[PointWithDir]bool{}

	for node := range graph {
		dist[node] = math.MaxInt
	}

	dist[start] = 0
	pq := datastructures.PriorityQueue{Node{start, dist[start]}}

	for len(pq) != 0 {
		n := *pq.Remove()
		nodeWithCost := n.(Node)

		node := nodeWithCost.adj

		if nodeWithCost.cost != dist[nodeWithCost.adj] {
			continue
		}

		for _, adjNode := range graph[node] {
			if dist[adjNode.adj] > dist[node]+adjNode.cost {
				dist[adjNode.adj] = dist[node] + adjNode.cost
				pq.Insert(Node{adjNode.adj, dist[adjNode.adj]})

				if _, ok := prev[adjNode.adj]; !ok {
					prev[adjNode.adj] = map[PointWithDir]bool{}
				}
				prev[adjNode.adj][node] = true

			} else if dist[adjNode.adj] == dist[node]+adjNode.cost {
				if _, ok := prev[adjNode.adj]; !ok {
					prev[adjNode.adj] = map[PointWithDir]bool{}
				}
				prev[adjNode.adj][node] = true
			}
		}
	}

	return dist, prev
}

func solvePartOne(graph map[PointWithDir][]Node, start PointWithDir, end Point) int {
	dist, _ := dijkstra(graph, start)

	res := math.MaxInt

	for node, dst := range dist {
		if node.p == end {
			res = int(math.Min(float64(res), float64(dst)))
		}
	}

	if res == math.MaxInt {
		return -1
	}

	return res
}

func allPathsFromAtoB(dijkstraGraph map[PointWithDir]map[PointWithDir]bool, sols *[][]PointWithDir, usedNodes map[PointWithDir]bool, currPath []PointWithDir, targetNode PointWithDir) {
	lastNode := currPath[len(currPath)-1]
	if lastNode == targetNode {
		*sols = append(*sols, currPath)
	} else {
		for adj := range dijkstraGraph[lastNode] {
			if _, ok := usedNodes[adj]; !ok {
				currPath = append(currPath, adj)
				usedNodes[adj] = true
				allPathsFromAtoB(dijkstraGraph, sols, usedNodes, currPath, targetNode)
				usedNodes[adj] = false
				currPath = currPath[:len(currPath)-1]
			}
		}
	}
}

func printGrid(grid [][]rune) {
	for _, row := range grid {
		fmt.Println(string(row))
	}
}

func pathCost(graph map[PointWithDir][]Node, path []PointWithDir) int {
	cost := 0
	currNode := path[len(path)-1]

	for idx := len(path) - 2; idx >= 0; idx-- {
		nextNode := path[idx]
		adj := graph[currNode]
		found := false
		for _, node := range adj {
			if node.adj == nextNode {
				found = true
				cost += node.cost
				break
			}
		}
		if !found {
			// panic("Path is invalid!")
		}
		currNode = nextNode
	}

	return cost
}

func solvePartTwo(grid [][]rune, graph map[PointWithDir][]Node, start PointWithDir, end Point) int {
	_, prev := dijkstra(graph, start)

	// fmt.Println(prev)

	sols := [][]PointWithDir{}
	allPathsFromAtoB(prev, &sols, map[PointWithDir]bool{{end, S}: true}, []PointWithDir{{end, S}}, start)
	// allPathsFromAtoB(prev, &sols, map[PointWithDir]bool{{end, V}: true}, []PointWithDir{{end, V}}, start.p)
	// allPathsFromAtoB(prev, &sols, map[PointWithDir]bool{{end, E}: true}, []PointWithDir{{end, E}}, start.p)
	// allPathsFromAtoB(prev, &sols, map[PointWithDir]bool{{end, N}: true}, []PointWithDir{{end, N}}, start.p)

	flatten := map[Point]bool{}
	r := 0
	fmt.Println("Num of solfs: ", len(sols))
	for _, s := range sols {
		fmt.Println("Length of sol: ", len(s))
		fmt.Println("Path cost: ", pathCost(graph, s))
		copyGrid := [][]rune{}
		for _, row := range grid {
			copyGrid = append(copyGrid, append([]rune{}, row...))
		}
		for _, p := range s {
			if _, ok := flatten[p.p]; !ok {
				r++
				flatten[p.p] = true
			}
			copyGrid[p.p.x][p.p.y] = 'O'
		}
		printGrid(copyGrid)
		fmt.Println()
	}
	return r
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	graph, start, end := makeGraph(data)

	// fmt.Println(graph, start, end)

	// if len(os.Args) != 2 {
	// 	panic("Exactly one arg is expected")
	// }
	// arg := os.Args[1]

	// if arg != "1" && arg != "2" {
	// 	panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	// }

	// if arg == "1" {
	// 	println(solvePartOne(graph, start, end))
	// } else {
	grid := [][]rune{}
	for _, row := range data {
		grid = append(grid, []rune(row))
	}
	println(solvePartTwo(grid, graph, start, end))
	// }
}
