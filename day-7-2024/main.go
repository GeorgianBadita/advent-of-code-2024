package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	result    int
	operators []int
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

func parseInputToEq(lines []string) []Equation {
	eqs := []Equation{}

	for _, line := range lines {
		eqParts := strings.Split(line, ":")
		result, _ := strconv.Atoi(eqParts[0])
		operators := strings.Split(strings.Trim(eqParts[1], " "), " ")
		ops := []int{}
		for _, operator := range operators {
			op, _ := strconv.Atoi(operator)
			ops = append(ops, op)
		}
		eqs = append(eqs, Equation{result, ops})
	}
	return eqs
}

func solvePartOne(equations []Equation) int {
	res := 0
	for _, eq := range equations {
		pow := 1 << (len(eq.operators) - 1)
		for try := 0; try < pow; try++ {
			config := []int{}

			for p := 0; p < len(eq.operators)-1; p++ {
				pthBit := (1 << p) & try
				if pthBit > 0 {
					config = append(config, 1)
				} else {
					config = append(config, 0)
				}

			}

			result := eq.operators[0]
			for idx, elem := range config {
				if result > eq.result {
					break
				}
				if elem == 0 {
					result += eq.operators[idx+1]
				} else {
					result *= eq.operators[idx+1]
				}
			}
			if result == eq.result {
				res += eq.result
				break
			}
		}
	}

	return res
}

func solvePartTwo(equations []Equation) int {
	res := 0
	for _, eq := range equations {
		// pow := 1 << (len(eq.operators) - 1)
		pow := int(math.Pow(3, float64(len(eq.operators)-1)))
		for try := 0; try < pow; try++ {
			config := []int{}
			tryCpy := try
			for p := 0; p < len(eq.operators)-1; p++ {
				config = append(config, tryCpy%3)
				tryCpy = tryCpy / 3
			}

			result := eq.operators[0]
			for idx, elem := range config {
				if result > eq.result {
					break
				}
				if elem == 0 {
					result += eq.operators[idx+1]
				} else if elem == 1 {
					result *= eq.operators[idx+1]
				} else {
					left := strconv.Itoa(result)
					right := strconv.Itoa(eq.operators[idx+1])
					concat := left + right
					num, _ := strconv.Atoi(concat)
					result = num
				}
			}
			if result == eq.result {
				res += eq.result
				break
			}
		}
	}

	return res
}

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	equations := parseInputToEq(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(equations))
	} else {
		println(solvePartTwo(equations))
	}
}
