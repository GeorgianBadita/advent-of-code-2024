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

func bfs(grid [][]int, start Point, end Point) [][]int {
	clonedGrid := [][]int{}

	for _, row := range grid {
		clonedGrid = append(clonedGrid, append([]int{}, row...))
	}

	queue := datastructures.Queue{start}
	dist := [][]int{}

	for idx := 0; idx < len(clonedGrid); idx++ {
		row := []int{}
		for jdx := 0; jdx < len(clonedGrid[0]); jdx++ {
			row = append(row, -1)
		}
		dist = append(dist, row)
	}

	dist[start.x][start.y] = 0

	for {
		n := queue.Dequeue()
		node := n.(Point)
		if node == end {
			return dist
		}

		for _, adj := range getValidNeighbours(clonedGrid, node) {
			if clonedGrid[adj.x][adj.y] != -1 && dist[adj.x][adj.y] == -1 {
				clonedGrid[adj.x][adj.y] = -1
				dist[adj.x][adj.y] = 1 + dist[node.x][node.y]
				queue.Enqueue(adj)
			}
		}
	}
}

func solvePartOne(grid [][]int, start Point, end Point) int {
	dist := bfs(grid, start, end)
	cnt := 0
	for idx := 0; idx < len(grid); idx++ {
		for jdx := 0; jdx < len(grid[0]); jdx++ {
			if grid[idx][jdx] == -1 {
				continue
			}
			for _, adj := range []Point{{idx + 2, jdx}, {idx + 1, jdx + 1}, {idx, jdx + 2}, {idx - 1, jdx + 1}} {
				if !validCoords(grid, adj) || grid[adj.x][adj.y] == -1 {
					continue
				}
				if datastructures.Abs(dist[idx][jdx]-dist[adj.x][adj.y]) >= 102 {
					cnt += 1
				}
			}
		}
	}
	return cnt
}

func solvePartTwo(grid [][]int, start Point, end Point) int {
	dist := bfs(grid, start, end)
	cnt := 0
	for idx := 0; idx < len(grid); idx++ {
		for jdx := 0; jdx < len(grid[0]); jdx++ {
			if grid[idx][jdx] == -1 {
				continue
			}
			for rd := 2; rd <= 20; rd++ {
				for dr := 0; dr <= rd; dr++ {
					cd := rd - dr
					for adj := range map[Point]bool{{idx + dr, jdx + cd}: true, {idx - dr, jdx - cd}: true, {idx + dr, jdx - cd}: true, {idx - dr, jdx + cd}: true} {
						if !validCoords(grid, adj) || grid[adj.x][adj.y] == -1 {
							continue
						}
						if dist[idx][jdx]-dist[adj.x][adj.y] >= 100+rd {
							cnt += 1
						}
					}
				}
			}

		}
	}
	return cnt
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
		println(solvePartTwo(grid, start, end))
	}
}
