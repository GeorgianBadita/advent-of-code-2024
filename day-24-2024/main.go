package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Pair[T any, V any] struct {
	x T
	y V
}

type Operation struct {
	leftOperand  string
	rightOpenrad string
	operation    OperationType
	result       string
}

type OperationType int

const (
	AND OperationType = iota
	OR
	XOR
)

func and(x int, y int) int {
	return x & y
}

func or(x int, y int) int {
	return x | y
}

func xor(x int, y int) int {
	return x ^ y
}

var OPER_TO_FUNC = map[OperationType]func(int, int) int{
	AND: and,
	OR:  or,
	XOR: xor,
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

func parseInput(data []string) (map[string]int, []Operation) {
	initialVals := map[string]int{}
	operations := []Operation{}
	parseInitialVals := true
	for _, line := range data {
		if line == "\n" || len(line) == 0 {
			parseInitialVals = false
			continue
		}

		if parseInitialVals {
			valsSplit := strings.Split(line, ": ")
			value, _ := strconv.Atoi(valsSplit[1])
			initialVals[valsSplit[0]] = value
		} else {
			// x00 AND y00 -> wrs
			opsSplit := strings.Split(line, " ")
			operation := XOR
			if opsSplit[1] == "AND" {
				operation = AND
			} else if opsSplit[1] == "OR" {
				operation = OR
			}
			operations = append(operations, Operation{
				leftOperand:  opsSplit[0],
				rightOpenrad: opsSplit[2],
				operation:    operation,
				result:       opsSplit[4],
			})
		}
	}

	return initialVals, operations
}

func solvePartOne(initialVals map[string]int, operations []Operation) uint64 {
	zOperations := []Pair[string, int]{}
	for len(operations) > 0 {
		idxToDel := map[int]bool{}
		for idx := 0; idx < len(operations); idx++ {
			op := operations[idx]
			if op1Val, ok1 := initialVals[op.leftOperand]; ok1 {
				if op2Val, ok2 := initialVals[op.rightOpenrad]; ok2 {
					initialVals[op.result] = OPER_TO_FUNC[op.operation](op1Val, op2Val)
					idxToDel[idx] = true
				}
			}
		}

		newOperations := []Operation{}
		for idx := 0; idx < len(operations); idx++ {
			if _, ok := idxToDel[idx]; !ok {
				newOperations = append(newOperations, operations[idx])
			}
		}
		operations = newOperations
	}
	for k, v := range initialVals {
		if strings.HasPrefix(k, "z") {
			zOperations = append(zOperations, Pair[string, int]{k, v})
		}
	}

	return makeOutputNum(zOperations)
}

func makeOutputNum(bits []Pair[string, int]) uint64 {
	slices.SortFunc(bits, func(a Pair[string, int], b Pair[string, int]) int {
		return strings.Compare(a.x, b.x)
	})

	resNum := uint64(0)

	for idx := uint64(0); idx < uint64(len(bits)); idx++ {
		if bits[idx].y == 1 {
			resNum |= 1 << idx
		}
	}

	return resNum
}

func solvePartTwo(initialVals map[string]int, operations []Operation) string {
	xOpaeration := []Pair[string, int]{}
	yOpaeration := []Pair[string, int]{}

	for k, v := range initialVals {
		if strings.HasPrefix(k, "x") {
			xOpaeration = append(xOpaeration, Pair[string, int]{k, v})
		}
		if strings.HasPrefix(k, "y") {
			yOpaeration = append(yOpaeration, Pair[string, int]{k, v})
		}
	}

	x := makeOutputNum(xOpaeration)
	y := makeOutputNum(yOpaeration)

	fmt.Println(x, y, x+y)

	return ""
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	initialVals, operations := parseInput(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(initialVals, operations))
	} else {
		println(solvePartTwo(initialVals, operations))
	}
}
