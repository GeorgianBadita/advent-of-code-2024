package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
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

func validCoord(grid [][]int, p Point) bool {
	return p.x >= 0 && p.x < len(grid) && p.y >= 0 && p.y < len(grid[0])
}

func findNeighbours(grid [][]int, p Point) []Point {
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	res := []Point{}
	for d := 0; d < len(dx); d++ {
		nP := Point{p.x + dx[d], p.y + dy[d]}
		if validCoord(grid, nP) && grid[nP.x][nP.y] == 0 {
			res = append(res, nP)
		}
	}
	return res
}

func makeGridAfterFirstNBytes(data []string, numRows int, numBytes int) [][]int {
	res := [][]int{}

	for idx := 0; idx <= numRows; idx++ {
		row := []int{}
		for jdx := 0; jdx <= numRows; jdx++ {
			row = append(row, 0)
		}
		res = append(res, row)
	}

	for idx := 0; idx < numBytes; idx++ {
		row := data[idx]
		coordsStr := strings.Split(row, ",")
		x, _ := strconv.Atoi(coordsStr[0])
		y, _ := strconv.Atoi(coordsStr[1])

		res[x][y] = -1
	}

	return res
}

func bfs(grid [][]int, start Point, target Point) int {
	queue := datastructures.Queue{start}

	for !queue.IsEmpty() {
		head := queue.Dequeue()
		point := head.(Point)

		if point == target {
			return grid[point.x][point.y]
		}

		neighbours := findNeighbours(grid, point)

		for _, adj := range neighbours {
			grid[adj.x][adj.y] = 1 + grid[point.x][point.y]
			queue.Enqueue(adj)
		}
	}
	return -1
}

func solvePartOne(data []string, start Point, target Point, numRows int, numBytes int) int {
	grid := makeGridAfterFirstNBytes(data, numRows, numBytes)
	return bfs(grid, start, target)
}

func solvePartTwo(data []string, start Point, target Point, numRows int, numBytes int) Point {
	grid := makeGridAfterFirstNBytes(data, numRows, numBytes)
	for idx := numBytes + 1; idx < len(data); idx++ {
		row := data[idx]
		coordsStr := strings.Split(row, ",")
		x, _ := strconv.Atoi(coordsStr[0])
		y, _ := strconv.Atoi(coordsStr[1])

		grid[x][y] = -1

		clonnedGrid := [][]int{}
		for _, row := range grid {
			clonnedGrid = append(clonnedGrid, append([]int{}, row...))
		}

		currSol := bfs(clonnedGrid, start, target)
		if currSol == -1 {
			return Point{x, y}
		}
	}
	panic("No byte falls on the critical path")
}

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
		println(solvePartOne(data, Point{0, 0}, Point{70, 70}, 70, 1024))
	} else {
		res := solvePartTwo(data, Point{0, 0}, Point{70, 70}, 70, 1024)
		println(fmt.Sprintf("%d,%d", res.x, res.y))
	}
}
