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

func converString(str string) []uint64 {
	res := []uint64{}

	split := strings.Split(str, " ")
	for _, numStr := range split {
		num, _ := strconv.Atoi(numStr)
		res = append(res, uint64(num))
	}

	return res
}

func numDigits(num uint64) uint {
	if num == 0 {
		return 1
	}
	numDig := uint(0)
	for num > 0 {
		numDig++
		num /= 10
	}
	return numDig
}

func getFreq(arr []uint64) map[uint64]int {
	res := map[uint64]int{}

	for _, x := range arr {
		res[x]++
	}

	return res
}

func getNextFromFreq(freq map[uint64]int) map[uint64]int {
	countZeroes := freq[0]
	delete(freq, 0)

	newMap := map[uint64]int{}

	for k, v := range freq {
		numDigits := numDigits(k)

		if numDigits%2 == 1 {
			newMap[k*2024] += v
		} else {
			halfDigs := numDigits / 2
			halfPow := datastructures.Pow(10, halfDigs)
			left := k / uint64(halfPow)
			right := k % uint64(halfPow)
			newMap[left] += v
			newMap[right] += v
		}
	}

	if countZeroes > 0 {
		newMap[1] += countZeroes
	}

	return newMap
}

func solve(arr []uint64, steps int) uint64 {
	initialFreq := getFreq(arr)

	for idx := 0; idx < steps; idx++ {
		initialFreq = getNextFromFreq(initialFreq)
	}

	res := uint64(0)
	for _, v := range initialFreq {
		res += uint64(v)
	}

	return res
}

func solvePartOne(arr []uint64) uint64 {
	return solve(arr, 25)
}

func solvePartTwo(arr []uint64) uint64 {
	return solve(arr, 75)
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	arr := converString(data[0])

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(arr))
	} else {
		println(solvePartTwo(arr))
	}
}
