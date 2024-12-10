package main

import (
	"aoc-2024/datastructures"
	"bufio"
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

func makeDigitMatrix(data []string) [][]uint8 {
	res := [][]uint8{}

	for _, line := range data {
		numLine := []uint8{}
		for _, elem := range line {
			numLine = append(numLine, uint8(int(elem)-'0'))
		}

		res = append(res, numLine)
	}
	return res
}

func validCoords(matrix [][]uint8, point Point) bool {
	return point.x >= 0 && point.x < len(matrix) && point.y >= 0 && point.y < len(matrix[0])
}

func getValidNeighbours(matrix [][]uint8, point Point) []Point {
	res := []Point{}
	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for idx := 0; idx < len(dx); idx++ {
		np := Point{point.x + dx[idx], point.y + dy[idx]}
		if validCoords(matrix, np) {
			res = append(res, np)
		}
	}

	return res
}

func getNumTrails(matrix [][]uint8, idx int, jdx int, uniqueTrails bool) uint {
	count := uint(0)
	queue := datastructures.Queue{Point{idx, jdx}}
	seenTrail := map[Point]bool{}

	for !queue.IsEmpty() {
		currPoint := queue.Dequeue()
		point := currPoint.(Point)

		val := matrix[point.x][point.y]
		if val == 9 {
			_, ok := seenTrail[point]
			if !ok || !uniqueTrails {
				count += 1
				seenTrail[point] = true
			}
			continue
		}

		neighbours := getValidNeighbours(matrix, point)

		for _, adjPoint := range neighbours {
			nextVal := matrix[adjPoint.x][adjPoint.y]

			if nextVal == val+1 {
				queue.Enqueue(adjPoint)
			}
		}
	}

	return count
}

func solvePartOne(matrix [][]uint8) uint {
	count := uint(0)
	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			if matrix[idx][jdx] == 0 {
				numTrails := getNumTrails(matrix, idx, jdx, true)
				count += numTrails
			}
		}
	}
	return count
}

func solvePartTwo(matrix [][]uint8) uint {
	count := uint(0)
	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			if matrix[idx][jdx] == 0 {
				numTrails := getNumTrails(matrix, idx, jdx, false)
				count += numTrails
			}
		}
	}
	return count
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	digitMatrix := makeDigitMatrix(data)

	// fmt.Println(digitMatrix)
	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(digitMatrix))
	} else {
		println(solvePartTwo(digitMatrix))
	}
}
