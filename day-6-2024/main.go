package main

import (
	"bufio"
	"fmt"
	"os"
)

type Pair struct {
	x int
	y int
}

type PairWithDirection struct {
	x         int
	y         int
	direction rune
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

func validCoords(matrix [][]rune, x int, y int) bool {
	return x >= 0 && x < len(matrix) && y >= 0 && y < len(matrix[0])
}

func getNext(matrix [][]rune, x int, y int) Pair {
	if matrix[x][y] == '^' {
		return Pair{x - 1, y}
	}
	if matrix[x][y] == '>' {
		return Pair{x, y + 1}
	}
	if matrix[x][y] == '<' {
		return Pair{x, y - 1}
	}
	if matrix[x][y] == 'v' {
		return Pair{x + 1, y}
	}
	panic(fmt.Sprintf("Current pointer at (%d, %d) has value %c instead of <>^v", x, y, matrix[x][y]))
}

func rotate(matrix [][]rune, x int, y int) [][]rune {
	if !validCoords(matrix, x, y) {
		panic(fmt.Sprintf("Current pointer at (%d, %d) is outside of matrix bounds", x, y))
	}
	// This is here just to check whether the function panics (i.e. the current pointer is <>^v)
	getNext(matrix, x, y)

	if (matrix)[x][y] == '>' {
		(matrix)[x][y] = 'v'
	} else if (matrix)[x][y] == 'v' {
		(matrix)[x][y] = '<'
	} else if (matrix)[x][y] == '<' {
		(matrix)[x][y] = '^'
	} else if (matrix)[x][y] == '^' {
		(matrix)[x][y] = '>'
	}
	return matrix
}

func getInitialCoords(matrix [][]rune) Pair {
	x, y := 0, 0
	for idx := 0; idx < len(matrix); idx++ {
		foundStart := false
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			if matrix[idx][jdx] == '^' {
				x, y = idx, jdx
				foundStart = true
				break
			}
		}
		if foundStart {
			break
		}
	}
	return Pair{x, y}
}

func getVisitedPath(matrix [][]rune) []Pair {
	startCoords := getInitialCoords(matrix)
	visited := map[Pair]bool{}

	path := []Pair{}
	currPair := Pair{startCoords.x, startCoords.y}

	for {
		if _, ok := visited[currPair]; !ok {
			path = append(path, currPair)
			visited[currPair] = true
		}
		nextP := getNext(matrix, currPair.x, currPair.y)
		if !validCoords(matrix, nextP.x, nextP.y) {
			break
		}
		if matrix[nextP.x][nextP.y] == '#' {
			matrix = rotate(matrix, currPair.x, currPair.y)
		} else {
			currVal := matrix[currPair.x][currPair.y]
			matrix[currPair.x][currPair.y] = '.'
			matrix[nextP.x][nextP.y] = currVal
			currPair = nextP
		}
	}

	return path
}

func solvePartOne(matrix [][]rune) int {
	return len(getVisitedPath(matrix))
}

func hasLoop(matrix [][]rune, x int, y int) bool {
	visited := map[PairWithDirection]bool{}
	currPair := PairWithDirection{x, y, matrix[x][y]}
	for {
		if _, ok := visited[currPair]; ok {
			return true
		}
		visited[currPair] = true
		nextP := getNext(matrix, currPair.x, currPair.y)
		if !validCoords(matrix, nextP.x, nextP.y) {
			return false
		}
		if matrix[nextP.x][nextP.y] == '#' {
			matrix = rotate(matrix, currPair.x, currPair.y)
			currPair.direction = matrix[currPair.x][currPair.y]
		} else {
			currVal := matrix[currPair.x][currPair.y]
			matrix[currPair.x][currPair.y] = '.'
			matrix[nextP.x][nextP.y] = currVal
			currPair = PairWithDirection{nextP.x, nextP.y, currVal}
		}
	}
}

func solvePartTwo(matrix [][]rune) int {
	matrixCopy := [][]rune{}
	for _, line := range matrix {
		matrixCopy = append(matrixCopy, append([]rune(nil), line...))
	}

	path := getVisitedPath(matrix)
	matrix = matrixCopy
	possibleSolutions := 0
	startCoords := getInitialCoords(matrix)

	for idx := 1; idx < len(path); idx++ {
		matrixCopy := [][]rune{}
		for _, line := range matrix {
			matrixCopy = append(matrixCopy, append([]rune(nil), line...))
		}

		x, y := path[idx].x, path[idx].y
		matrix[x][y] = '#'
		if hasLoop(matrix, startCoords.x, startCoords.y) {
			possibleSolutions += 1
		}
		matrix = matrixCopy
	}

	return possibleSolutions
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	matrix := [][]rune{}
	for _, line := range data {
		matrix = append(matrix, []rune(line))
	}

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(matrix))
	} else {
		println(solvePartTwo(matrix))
	}
}
