package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"os"
	"strconv"
	"strings"
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

func parseLines(data []string) [][]int {
	res := [][]int{}

	for _, element := range data {
		row := []int{}
		split := strings.Split(element, " ")
		for _, number := range split {
			num, _ := strconv.Atoi(number)
			row = append(row, num)
		}
		res = append(res, row)
	}
	return res
}

func checkGrowthRate(num1 int, num2 int) bool {
	abs := datastructures.Abs(num1 - num2)
	return abs >= 1 && abs <= 3
}

func solvePartOne(lines [][]int) int {
	count := 0
	for _, row := range lines {
		dir := 0
		if row[0] == row[1] || !checkGrowthRate(row[0], row[1]) {
			continue
		}
		if row[0] > row[1] {
			dir = 1
		}

		goodRow := true
		for idx := 1; idx < len(row); idx++ {
			if dir == 0 {
				if row[idx-1] >= row[idx] || !checkGrowthRate(row[idx-1], row[idx]) {
					goodRow = false
					break
				}
			} else {
				if row[idx-1] <= row[idx] || !checkGrowthRate(row[idx-1], row[idx]) {
					goodRow = false
					break
				}
			}
		}

		if goodRow {
			count += 1
		}
	}
	return count
}

func checkAscending(row []int) int {
	for idx := 1; idx < len(row); idx += 1 {
		if row[idx] <= row[idx-1] || !checkGrowthRate(row[idx], row[idx-1]) {
			return idx
		}
	}
	return len(row)
}

func checkDescending(row []int) int {
	for idx := 1; idx < len(row); idx += 1 {
		if row[idx] >= row[idx-1] || !checkGrowthRate(row[idx], row[idx-1]) {
			return idx
		}
	}
	return len(row)
}

func removeIndex(line []int, idx int) []int {
	res := []int{}
	for jdx := 0; jdx < len(line); jdx++ {
		if jdx == idx {
			continue
		}
		res = append(res, line[jdx])
	}

	return res
}

func solvePartTwo(lines [][]int) int {
	for idx := 0; idx < len(lines); idx++ {
		line := lines[idx]
		asc := checkAscending(line)
		dsc := checkDescending(line)

		if asc == len(line) || dsc == len(line) {
			continue
		}

		removedAsc := removeIndex(line, asc)
		if checkAscending(removedAsc) == len(removedAsc) {
			lines[idx] = removedAsc
			continue
		}

		removedBeforeAsc := removeIndex(line, asc-1)
		if checkAscending(removedBeforeAsc) == len(removedBeforeAsc) {
			lines[idx] = removedBeforeAsc
			continue
		}

		removedBeforeDsc := removeIndex(line, dsc-1)
		if checkDescending(removedBeforeDsc) == len(removedBeforeDsc) {
			lines[idx] = removedBeforeDsc
			continue
		}
		removedDesc := removeIndex(line, dsc)
		if checkDescending(removedDesc) == len(removedDesc) {
			lines[idx] = removedDesc
		}
	}

	return solvePartOne(lines)
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	lines := parseLines(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(lines))
	} else {
		println(solvePartTwo(lines))
	}
}
