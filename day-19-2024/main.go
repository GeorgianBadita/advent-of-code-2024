package main

import (
	"bufio"
	"os"
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

func parseData(data []string) ([]string, []string) {
	return strings.Split(data[0], ", "), data[2:]
}

func matchOneNum(options []string, stringToMatch string, memo map[string]int) int {
	if len(stringToMatch) == 0 {
		return 1
	}

	if v, ok := memo[stringToMatch]; ok {
		return v
	}

	res := 0
	for idx := 0; idx < len(options); idx++ {
		if strings.HasPrefix(stringToMatch, options[idx]) {
			res += matchOneNum(options, stringToMatch[len(options[idx]):], memo)
		}
	}

	memo[stringToMatch] = res
	return res
}

func solvePartOne(options []string, stringsToMatch []string) int {
	res := 0
	memo := map[string]int{}

	for _, stringToMatch := range stringsToMatch {
		if matchOneNum(options, stringToMatch, memo) > 0 {
			res++
		}
	}
	return res
}

func solvePartTwo(options []string, stringsToMatch []string) int {
	res := 0
	memo := map[string]int{}

	for _, stringToMatch := range stringsToMatch {
		res += matchOneNum(options, stringToMatch, memo)
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

	options, stringsToMatch := parseData(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(options, stringsToMatch))
	} else {
		println(solvePartTwo(options, stringsToMatch))
	}
}
