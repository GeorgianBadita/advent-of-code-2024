package main

import (
	"bufio"
	"os"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
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

func abs[T constraints.Signed | constraints.Float](v T) T {
	return max(-v, v)
}

func makeLists(lines []string) ([]int, []int) {
	list1 := []int{}
	list2 := []int{}

	for _, line := range lines {
		elements := strings.Split(line, "   ")
		num1, _ := strconv.Atoi(elements[0])
		num2, _ := strconv.Atoi(elements[1])

		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	return list1, list2
}

func solvePart1(list1 []int, list2 []int) int {
	slices.Sort(list1)
	slices.Sort(list2)

	res := 0

	for i := 0; i < len(list1); i++ {
		res += abs(list1[i] - list2[i])
	}
	return res
}

func solvePart2(list1 []int, list2 []int) int {
	res := 0

	for _, elem := range list1 {
		count := 0
		for _, elem2 := range list2 {
			if elem == elem2 {
				count += 1
			}
		}

		res += count * elem
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
	list1, list2 := makeLists(data)

	if arg == "1" {
		println(solvePart1(list1, list2))
	} else {
		println(solvePart2(list1, list2))
	}
}
