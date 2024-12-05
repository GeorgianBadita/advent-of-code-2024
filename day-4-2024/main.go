package main

import (
	"bufio"
	"os"
)

type Pair struct {
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

func validCoords(matrix []string, x int, y int) bool {
	return x >= 0 && x < len(matrix) && y >= 0 && y < len(matrix[0])
}

func checkMAS(matrix []string, x int, y int) bool {
	if matrix[x][y] != 'M' && matrix[x][y] != 'S' {
		return false
	}

	nextHeadCoords := Pair{x, y + 2}
	if !validCoords(matrix, nextHeadCoords.x, nextHeadCoords.y) ||
		(matrix[nextHeadCoords.x][nextHeadCoords.y] != 'M' && matrix[nextHeadCoords.x][nextHeadCoords.y] != 'S') {
		return false
	}

	aCoords := Pair{x + 1, y + 1}

	if !validCoords(matrix, aCoords.x, aCoords.y) || matrix[aCoords.x][aCoords.y] != 'A' {
		return false
	}

	// first branch
	if matrix[x][y] == 'M' {
		sCoords := Pair{x + 2, y + 2}
		if !validCoords(matrix, sCoords.x, sCoords.y) || matrix[sCoords.x][sCoords.y] != 'S' {
			return false
		}
	}

	if matrix[x][y] == 'S' {
		mCoords := Pair{x + 2, y + 2}
		if !validCoords(matrix, mCoords.x, mCoords.y) || matrix[mCoords.x][mCoords.y] != 'M' {
			return false
		}
	}

	// second branch
	if matrix[nextHeadCoords.x][nextHeadCoords.y] == 'M' {
		sCoords := Pair{x + 2, y}
		if !validCoords(matrix, sCoords.x, sCoords.y) || matrix[sCoords.x][sCoords.y] != 'S' {
			return false
		}
	}

	if matrix[nextHeadCoords.x][nextHeadCoords.y] == 'S' {
		mCoords := Pair{x + 2, y}
		if !validCoords(matrix, mCoords.x, mCoords.y) || matrix[mCoords.x][mCoords.y] != 'M' {
			return false
		}
	}

	return true
}

func countWordsFrom(matrix []string, word string, x int, y int) int {
	dx := []int{1, -1, 0, 0, 1, -1, -1, 1}
	dy := []int{1, -1, 1, -1, 0, 0, 1, -1}

	count := 0
	for idx := 0; idx < len(dx); idx++ {
		wordPos := 1
		xx := x
		yy := y
		for wordPos < len(word) {
			nX := xx + dx[idx]
			nY := yy + dy[idx]
			if !validCoords(matrix, nX, nY) || matrix[nX][nY] != word[wordPos] {
				break
			}
			wordPos += 1
			xx = nX
			yy = nY
		}
		if wordPos == len(word) {
			count += 1
		}
	}

	return count
}

func solvePartOne(matrix []string) int {
	total_count := 0
	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[idx]); jdx++ {
			if matrix[idx][jdx] == 'X' {
				total_count += countWordsFrom(matrix, "XMAS", idx, jdx)
			}
		}
	}
	return total_count
}

func solvePartTwo(matrix []string) int {
	total_count := 0
	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[idx]); jdx++ {
			if checkMAS(matrix, idx, jdx) {
				total_count += 1
			}
		}
	}
	return total_count
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
		print(solvePartOne(data))
	} else {
		print(solvePartTwo(data))
	}
}
