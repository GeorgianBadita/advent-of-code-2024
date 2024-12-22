package main

import (
	"bufio"
	"math"
	"os"
	"strconv"
)

type Sequence struct {
	x1 int
	x2 int
	x3 int
	x4 int
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

func solvePartTwo(nums []uint64) int {
	buyers := [][]int{}
	for _, num := range nums {
		buyer := []int{int(num) % 10}
		for idx := 0; idx < 2000; idx++ {
			num = getNextSecretNumber(num)
			buyer = append(buyer, int(num)%10)
		}
		buyers = append(buyers, buyer)
	}

	sequences := map[Sequence]int{}
	for _, buyer := range buyers {
		seen := map[Sequence]bool{}
		for idx := 0; idx < len(buyer)-4; idx++ {
			slice := buyer[idx : idx+5]
			a := slice[0]
			b := slice[1]
			c := slice[2]
			d := slice[3]
			e := slice[4]

			seq := Sequence{b - a, c - b, d - c, e - d}
			if _, ok := seen[seq]; ok {
				continue
			}
			seen[seq] = true
			sequences[seq] += e
		}
	}

	mx := math.MinInt64
	for _, res := range sequences {
		mx = max(mx, res)
	}
	return mx
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
		println(solvePartTwo(nums))
	}
}
