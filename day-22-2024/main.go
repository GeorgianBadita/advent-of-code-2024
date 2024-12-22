package main

import (
	"bufio"
	"os"
	"strconv"
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

func parseToUint(data []string) []uint64 {
	res := []uint64{}
	for _, v := range data {
		val, _ := strconv.Atoi(v)
		res = append(res, uint64(val))
	}
	return res
}

func expandNextK(from uint64, k int) uint64 {
	for k > 0 {
		from = getNextSecretNumber(from)
		k--
	}
	return from
}

func getNextSecretNumber(n uint64) uint64 {
	m := ((n * (1 << 6)) ^ n) % MOD
	m = ((m >> 5) ^ m) % MOD
	m = ((m * (1 << 11)) ^ m) % MOD

	return m
}

func solvePartOne(nums []uint64) uint64 {
	allSum := uint64(0)
	for _, num := range nums {
		allSum += expandNextK(num, 2000)
	}
	return allSum
}

const MOD = 16777216

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	nums := parseToUint(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(nums))
	} else {
		// placeholder for part 2
	}
}
