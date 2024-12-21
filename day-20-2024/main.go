package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x int
	y int
}

type Cheat struct {
	first  Point
	second Point
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

const IN_FILE_PATH = "./input.txt"

func parseMatrix(data []string) ([][]int, Point, Point) {
	res := [][]int{}
	var s, e Point
	for idx, row := range data {
		rowInt := []int{}
		for jdx, chr := range row {
			if chr == '#' {
				rowInt = append(rowInt, -1)
				continue
			}
			if chr == 'S' {
				s = Point{idx, jdx}
			}
			if chr == 'E' {
				e = Point{idx, jdx}
			}
			rowInt = append(rowInt, 0)
		}
		res = append(res, rowInt)
	}
	return res, s, e
}

func printGrid(grid [][]int, start Point, end Point) {
	for idx, row := range grid {
		str := ""
		for jdx, chr := range row {
			if idx == start.x && jdx == start.y {
				str += "S"
				continue
			}
			if idx == end.x && jdx == end.y {
				str += "E"
				continue
			}
			if chr == -1 {
				str += "#"
				continue
			}
			if chr == 0 {
				str += "."
				continue
			}
			str += "O"
		}
		fmt.Println(str)
	}
}

func validCoords(grid [][]int, p Point) bool {
	return p.x >= 0 && p.x < len(grid) && p.y >= 0 && p.y < len(grid[0])
}

func getValidNeighbours(grid [][]int, p Point) []Point {
	neighbours := []Point{}
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for d := 0; d < len(dx); d++ {
		nP := Point{p.x + dx[d], p.y + dy[d]}
		if validCoords(grid, nP) {
			neighbours = append(neighbours, nP)
		}
	}
	return neighbours
}

func bfs(grid [][]int, start Point, end Point, validCheat *Cheat) (int, map[Point]bool) {
	queue := datastructures.Queue{start}
	grid[start.x][start.y] = 1

	potentialCheats := map[Point]bool{}
	potentialCheats[start] = true
	usedFirst := false

	for !queue.IsEmpty() {
		n := queue.Dequeue()
		node := n.(Point)

		if validCheat != nil {
			if node == validCheat.second && !usedFirst {
				return datastructures.Pow(10, 9), nil
			}
			if node == validCheat.first {
				usedFirst = true
			}
		}

		if node == end {
			potentialCheats[end] = true
			return grid[node.x][node.y] - 1, potentialCheats
		}

		for _, adj := range getValidNeighbours(grid, node) {
			if grid[adj.x][adj.y] == 0 {
				grid[adj.x][adj.y] = 1 + grid[node.x][node.y]
				queue.Enqueue(adj)
			} else {
				potentialCheats[adj] = true
			}
		}
	}
	panic(fmt.Sprintf("There is no path from %v to %v", start, end))
}

func solvePartOne(grid [][]int, start Point, end Point) int {
	matrixClone := [][]int{}
	for _, row := range grid {
		matrixClone = append(matrixClone, append([]int{}, row...))
	}

	validCheats := map[Cheat]bool{}

	baseLine, potentialCheats := bfs(matrixClone, start, end, nil)

	for potentialCheat := range potentialCheats {
		neighbours := getValidNeighbours(grid, potentialCheat)
		for _, adj := range neighbours {
			if grid[adj.x][adj.y] == -1 {
				continue
			}
			first := potentialCheat
			second := Point{adj.x, adj.y}
			cheat := Cheat{first, second}
			_, firstSecond := validCheats[Cheat{first, second}]
			_, secondFirst := validCheats[Cheat{second, first}]

			if firstSecond || secondFirst {
				continue
			}

			validCheats[cheat] = true

		}
	}

	bestCheats := 0
	for validCheat := range validCheats {
		matrixClone := [][]int{}
		for _, row := range grid {
			matrixClone = append(matrixClone, append([]int{}, row...))
		}

		matrixClone[validCheat.first.x][validCheat.first.y] = 0
		matrixClone[validCheat.second.x][validCheat.second.y] = 0

		currRes, _ := bfs(matrixClone, start, end, &validCheat)

		if baseLine-currRes >= 100 {
			bestCheats += 1
		}

	}
	return bestCheats
}

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	grid, start, end := parseMatrix(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(grid, start, end))
	} else {
		// placeholder for part 2
	}
}
