package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Pair struct {
	fist   int
	second int
}

type Range = Pair

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

func combineLines(lines []string) string {
	res := ""
	for _, line := range lines {
		res += line
	}
	return res
}

func findMulPairs(str string) []Pair {
	pairs := []Pair{}
	r := regexp.MustCompile(`mul\((\d{1,3},\d{1,3})\)`)

	matches := r.FindAllStringSubmatch(str, -1)

	for _, match := range matches {
		numbers := strings.Split(match[1], ",")
		num1, _ := strconv.Atoi(numbers[0])
		num2, _ := strconv.Atoi(numbers[1])

		pairs = append(pairs, Pair{num1, num2})
	}

	return pairs
}

func solvePartOne(str string) int {
	res := 0
	for _, pair := range findMulPairs(str) {
		res += pair.fist * pair.second
	}
	return res
}

func stringHasFrom(str string, has string, idx int) bool {
	if idx+len(has) >= len(str) {
		return false
	}
	return str[idx:idx+len(has)] == has
}

func isDo(str string, idx int) bool {
	return stringHasFrom(str, "do()", idx)
}

func isDont(str string, idx int) bool {
	return stringHasFrom(str, "don't()", idx)
}

func solvePartTwo(str string) int {
	goodRanges := []Range{}
	currStart := 0
	state := "OK"
	idx := 0
	doSize := 4
	dontSize := 7

	for idx < len(str) {
		if state == "OK" && isDont(str, idx) {
			idx += dontSize
			goodRanges = append(goodRanges, Range{currStart, idx})
			state = "NOK"
			continue
		}
		if state == "NOK" && isDo(str, idx) {
			idx += doSize
			currStart = idx
			state = "OK"
			continue
		}
		idx += 1
	}

	if state == "OK" {
		goodRanges = append(goodRanges, Range{currStart, idx})
	}

	res := 0
	for _, rang := range goodRanges {
		res += solvePartOne(str[rang.fist:rang.second])
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

	combinedLines := combineLines(data)
	if arg == "1" {
		println(solvePartOne(combinedLines))
	} else {
		println(solvePartTwo(combinedLines))
	}
}
