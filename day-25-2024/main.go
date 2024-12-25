package main

import (
	"bufio"
	"os"
)

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

func locksAndKeys(data []string) ([][][]rune, [][][]rune) {
	keys := [][][]rune{}
	locks := [][][]rune{}
	currMatrix := [][]rune{}
	for idx := 0; idx < len(data); idx++ {
		if len(data[idx]) == 0 || data[idx] == "\n" {
			isKey := true
			for jdx := 0; jdx < len(currMatrix[0]); jdx++ {
				if currMatrix[0][jdx] == '#' {
					isKey = false
					break
				}
			}
			if isKey {
				keys = append(keys, currMatrix)
			} else {
				locks = append(locks, currMatrix)
			}
			currMatrix = [][]rune{}
			continue
		}
		currMatrix = append(currMatrix, []rune(data[idx]))
	}

	return keys, locks
}

func makeArr(matrix [][]rune, isKey bool) []int {
	res := []int{}
	dir := 1
	if isKey {
		dir = -1
	}

	for jdx := 0; jdx < len(matrix[0]); jdx++ {
		numIter := 0
		r := 0
		if isKey {
			r = len(matrix) - 1
		}

		for r >= 0 && r < len(matrix) && matrix[r][jdx] == '#' {
			r += dir
			numIter++
		}
		res = append(res, numIter-1)
	}
	return res
}

func solvePartOne(keys [][][]rune, locks [][][]rune) int {
	cnt := 0
	for idx := 0; idx < len(keys); idx++ {
		for jdx := 0; jdx < len(locks); jdx++ {
			lockArr := makeArr(locks[jdx], false)
			keyArr := makeArr(keys[idx], true)

			goodCombo := true

			for kdx := 0; kdx < len(lockArr); kdx++ {
				if lockArr[kdx]+keyArr[kdx] > 5 {
					goodCombo = false
					break
				}
			}

			if goodCombo {
				cnt++
			}
		}
	}
	return cnt
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	keys, locks := locksAndKeys(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(keys, locks))
	} else {
		// placeholder for part 2
	}
}
