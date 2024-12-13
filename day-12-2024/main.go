package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"os"
	"sort"
)

type DirType int

const (
	UP DirType = iota
	DOWN
	LEFT
	RIGHT
)

type Point struct {
	x int
	y int
}

type GridEdgeNode struct {
	p   Point
	dir DirType
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

func makeMatrixFromStringList(data []string) [][]rune {
	res := [][]rune{}
	for _, line := range data {
		res = append(res, []rune(line))
	}
	return res
}

func validCoords(matrix [][]rune, point Point) bool {
	return point.x >= 0 && point.x < len(matrix) && point.y >= 0 && point.y < len(matrix[0])
}

func getAdjCoords(p Point) []Point {
	res := []Point{}
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for idx := 0; idx < len(dx); idx++ {
		res = append(res, Point{p.x + dx[idx], p.y + dy[idx]})
	}

	return res
}

func fill(matrix *[][]rune, p Point, val rune, points *[]Point) {
	if !validCoords(*matrix, p) || (*matrix)[p.x][p.y] != val {
		return
	}

	*points = append(*points, p)

	(*matrix)[p.x][p.y] = '.'
	adjCoords := getAdjCoords(p)

	for _, adj := range adjCoords {
		fill(matrix, adj, val, points)
	}
}

func getAreaAndPerimeter(matrix [][]rune, points []Point) (uint, uint) {
	area := uint(len(points))
	perimeter := uint(0)

	for _, point := range points {
		diff := uint(0)
		adjCoords := getAdjCoords(point)
		for _, adj := range adjCoords {
			if !validCoords(matrix, adj) {
				diff++
			} else if matrix[point.x][point.y] != matrix[adj.x][adj.y] {
				diff++
			}
		}
		perimeter += diff
	}

	return area, perimeter
}

func getNeighbourAtDirection(p Point, d DirType) Point {
	switch d {
	case UP:
		return Point{p.x - 1, p.y}
	case DOWN:
		return Point{p.x + 1, p.y}
	case LEFT:
		return Point{p.x, p.y - 1}
	case RIGHT:
		return Point{p.x, p.y + 1}
	default:
		panic("Dir must be UP/DOWN/LEFT/RIGHT")
	}
}

func bfs(graph map[GridEdgeNode][]GridEdgeNode, start GridEdgeNode, visited *map[GridEdgeNode]bool) {
	queue := datastructures.Queue{}
	queue.Enqueue(start)
	(*visited)[start] = true

	for !queue.IsEmpty() {
		nq := queue.Dequeue()
		node := nq.(GridEdgeNode)

		for _, adj := range graph[node] {
			if _, ok := (*visited)[adj]; !ok {
				queue.Enqueue(adj)
				(*visited)[adj] = true
			}
		}
	}
}

func getAreaAndSides(matrix [][]rune, points []Point) (uint, uint) {
	area := uint(len(points))
	pointsOnPerimeter := []Point{}
	edgesGraph := map[GridEdgeNode][]GridEdgeNode{}

	for _, point := range points {
		adjCoords := getAdjCoords(point)
		diff := 0
		for _, adj := range adjCoords {
			if !validCoords(matrix, adj) {
				diff++
			} else if matrix[point.x][point.y] != matrix[adj.x][adj.y] {
				diff++
			}
		}
		if diff > 0 {
			pointsOnPerimeter = append(pointsOnPerimeter, point)
		}
	}

	sort.Slice(pointsOnPerimeter, func(i, j int) bool {
		if pointsOnPerimeter[i].x < pointsOnPerimeter[j].x {
			return true
		} else if pointsOnPerimeter[i].x > pointsOnPerimeter[j].x {
			return false
		}
		return pointsOnPerimeter[i].y <= pointsOnPerimeter[j].y
	})

	for _, point := range pointsOnPerimeter {
		upPoint := getNeighbourAtDirection(point, UP)
		if !validCoords(matrix, upPoint) || matrix[upPoint.x][upPoint.y] != matrix[point.x][point.y] {
			// this node exists
			edgesGraph[GridEdgeNode{point, UP}] = []GridEdgeNode{}
			// check whether the left point has an up node
			leftPoint := getNeighbourAtDirection(point, LEFT)
			if val, ok := edgesGraph[GridEdgeNode{leftPoint, UP}]; ok {
				val = append(val, GridEdgeNode{point, UP})
				edgesGraph[GridEdgeNode{leftPoint, UP}] = val

				edgesGraph[GridEdgeNode{point, UP}] = append(edgesGraph[GridEdgeNode{point, UP}], GridEdgeNode{leftPoint, UP})
			}
		}

		downPoint := getNeighbourAtDirection(point, DOWN)
		if !validCoords(matrix, downPoint) || matrix[downPoint.x][downPoint.y] != matrix[point.x][point.y] {
			// this node exists
			edgesGraph[GridEdgeNode{point, DOWN}] = []GridEdgeNode{}
			// check whether the left point has an up node
			leftPoint := getNeighbourAtDirection(point, LEFT)
			if val, ok := edgesGraph[GridEdgeNode{leftPoint, DOWN}]; ok {
				val = append(val, GridEdgeNode{point, DOWN})
				edgesGraph[GridEdgeNode{leftPoint, DOWN}] = val

				edgesGraph[GridEdgeNode{point, DOWN}] = append(edgesGraph[GridEdgeNode{point, DOWN}], GridEdgeNode{leftPoint, DOWN})
			}
		}

		leftPoint := getNeighbourAtDirection(point, LEFT)
		if !validCoords(matrix, leftPoint) || matrix[leftPoint.x][leftPoint.y] != matrix[point.x][point.y] {
			// this node exists
			edgesGraph[GridEdgeNode{point, LEFT}] = []GridEdgeNode{}
			// check whether the left point has an up node
			upPoint := getNeighbourAtDirection(point, UP)
			if val, ok := edgesGraph[GridEdgeNode{upPoint, LEFT}]; ok {
				val = append(val, GridEdgeNode{point, LEFT})
				edgesGraph[GridEdgeNode{upPoint, LEFT}] = val

				edgesGraph[GridEdgeNode{point, LEFT}] = append(edgesGraph[GridEdgeNode{point, LEFT}], GridEdgeNode{upPoint, LEFT})
			}
		}

		rightPoint := getNeighbourAtDirection(point, RIGHT)
		if !validCoords(matrix, rightPoint) || matrix[rightPoint.x][rightPoint.y] != matrix[point.x][point.y] {
			// this node exists
			edgesGraph[GridEdgeNode{point, RIGHT}] = []GridEdgeNode{}
			// check whether the left point has an up node
			upPoint := getNeighbourAtDirection(point, UP)
			if val, ok := edgesGraph[GridEdgeNode{upPoint, RIGHT}]; ok {
				val = append(val, GridEdgeNode{point, RIGHT})
				edgesGraph[GridEdgeNode{upPoint, RIGHT}] = val

				edgesGraph[GridEdgeNode{point, RIGHT}] = append(edgesGraph[GridEdgeNode{point, RIGHT}], GridEdgeNode{upPoint, RIGHT})
			}
		}
	}

	connectedComponents := uint(0)
	visited := map[GridEdgeNode]bool{}

	for node := range edgesGraph {
		if _, ok := visited[node]; !ok {
			connectedComponents++
			bfs(edgesGraph, node, &visited)
		}
	}

	return area, connectedComponents
}

func solvePartOne(matrix [][]rune) uint {
	res := uint(0)

	originalMatrix := [][]rune{}

	for _, line := range matrix {
		originalMatrix = append(originalMatrix, append([]rune{}, line...))
	}

	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			if matrix[idx][jdx] != '.' {
				points := []Point{}
				fill(&matrix, Point{idx, jdx}, matrix[idx][jdx], &points)
				area, perimeter := getAreaAndPerimeter(originalMatrix, points)
				res += area * perimeter
			}
		}
	}
	return res
}

func solvePartTwo(matrix [][]rune) uint {
	res := uint(0)

	originalMatrix := [][]rune{}

	for _, line := range matrix {
		originalMatrix = append(originalMatrix, append([]rune{}, line...))
	}

	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			if matrix[idx][jdx] != '.' {
				points := []Point{}
				fill(&matrix, Point{idx, jdx}, matrix[idx][jdx], &points)
				area, numSides := getAreaAndSides(originalMatrix, points)
				res += area * numSides
			}
		}
	}
	return res
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(makeMatrixFromStringList(data)))
	} else {
		println(solvePartTwo(makeMatrixFromStringList(data)))
	}
}
