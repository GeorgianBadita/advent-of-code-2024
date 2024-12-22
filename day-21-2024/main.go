package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Point struct {
	x int
	y int
}

type DirType int

const (
	UP DirType = iota
	DOWN
	LEFT
	RIGHT
)

var DIR_TO_ARROW = map[DirType]rune{
	UP:    '^',
	DOWN:  'v',
	LEFT:  '<',
	RIGHT: '>',
}

var DIR_TO_OFFSET = map[DirType]Point{
	UP:    {-1, 0},
	DOWN:  {1, 0},
	LEFT:  {0, -1},
	RIGHT: {0, 1},
}

var KEY_PAD = map[rune]Point{
	'7': {0, 0},
	'8': {0, 1},
	'9': {0, 2},
	'4': {1, 0},
	'5': {1, 1},
	'6': {1, 2},
	'1': {2, 0},
	'2': {2, 1},
	'3': {2, 2},
	'0': {3, 1},
	'A': {3, 2},
}

var ROBO_PAD = map[rune]Point{
	'^': {0, 1},
	'A': {0, 2},
	'<': {1, 0},
	'v': {1, 1},
	'>': {1, 2},
}

func align(curCoords Point, nextCoords Point, colDir DirType, lineDir DirType, invMap map[Point]rune, lineFirst bool) (string, bool) {
	res := ""
	if lineFirst {
		for curCoords.x != nextCoords.x {
			res += string(DIR_TO_ARROW[lineDir])
			curCoords.x += DIR_TO_OFFSET[lineDir].x
			if _, ok := invMap[curCoords]; !ok {
				return "", false
			}
		}

		for curCoords.y != nextCoords.y {
			res += string(DIR_TO_ARROW[colDir])
			curCoords.y += DIR_TO_OFFSET[colDir].y
			if _, ok := invMap[curCoords]; !ok {
				return "", false
			}
		}

		return res, true
	}

	for curCoords.y != nextCoords.y {
		res += string(DIR_TO_ARROW[colDir])
		curCoords.y += DIR_TO_OFFSET[colDir].y
		if _, ok := invMap[curCoords]; !ok {
			return "", false
		}
	}

	for curCoords.x != nextCoords.x {
		res += string(DIR_TO_ARROW[lineDir])
		curCoords.x += DIR_TO_OFFSET[lineDir].x
		if _, ok := invMap[curCoords]; !ok {
			return "", false
		}
	}

	return res, true
}

func getPathFrom(x rune, y rune, mapping map[rune]Point) []string {
	invMap := map[Point]rune{}

	for k, v := range mapping {
		invMap[v] = k
	}

	currCoords := mapping[x]
	nextCoords := mapping[y]

	colDir := RIGHT
	lineDir := UP

	if nextCoords.y > currCoords.y {
		colDir = RIGHT
	} else {
		colDir = LEFT
	}

	if nextCoords.x > currCoords.x {
		lineDir = DOWN
	} else {
		lineDir = UP
	}

	res1, correct1 := align(currCoords, nextCoords, colDir, lineDir, invMap, true)
	res2, correct2 := align(currCoords, nextCoords, colDir, lineDir, invMap, false)

	if correct1 && correct2 {
		return []string{res1 + "A", res2 + "A"}
	}

	if correct1 {
		return []string{res1 + "A"}
	}

	if correct2 {
		return []string{res2 + "A"}
	}
	panic("No correct answer")
}

func getCommandsToOutputStrings(str string, mapping map[rune]Point) []string {
	segments := map[string]bool{}
	currPoint := 'A'
	row := []rune(str)
	for idx := 0; idx < len(row); idx++ {
		if len(segments) == 0 {
			res := getPathFrom(currPoint, row[idx], mapping)
			for _, r := range res {
				segments[r] = true
			}
		} else {
			additions := getPathFrom(currPoint, row[idx], mapping)
			newSegments := map[string]bool{}
			for currSegment := range segments {
				for jdx := 0; jdx < len(additions); jdx++ {
					newSegments[currSegment+additions[jdx]] = true
				}
			}

			segments = newSegments
		}
		currPoint = row[idx]
	}
	strings := []string{}
	for r := range segments {
		strings = append(strings, r)
	}
	return strings
}

func getCommandsToOutputStringsMultipleStrings(strings []string, mapping map[rune]Point) []string {
	acc := map[string]bool{}
	for _, str := range strings {
		res := getCommandsToOutputStrings(str, mapping)
		for _, r := range res {
			acc[r] = true
		}
	}
	res := []string{}
	for r := range acc {
		res = append(res, r)
	}
	return res
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

func atoi(s string) int {
	res := ""
	for _, ch := range s {
		if ch != 'A' {
			res += string(ch)
		}
	}

	r, _ := strconv.Atoi(res)
	return r
}

const IN_FILE_PATH = "./input.txt"

func solvePartOne(rows []string) int {
	s := 0
	for _, r := range rows {
		firsts := getCommandsToOutputStrings(r, KEY_PAD)
		seconds := getCommandsToOutputStringsMultipleStrings(firsts, ROBO_PAD)
		thirds := getCommandsToOutputStringsMultipleStrings(seconds, ROBO_PAD)

		bestLen := datastructures.Pow(10, 9)

		for _, res := range thirds {
			if len(res) < bestLen {
				bestLen = len(res)
			}
		}
		s += bestLen * atoi(r)
	}
	return s
}

func solvePartTwo(rows []string) int {
	s := 0
	for i, r := range rows {
		fmt.Println("processing string ", i)
		firsts := getCommandsToOutputStrings(r, KEY_PAD)
		seconds := getCommandsToOutputStringsMultipleStrings(firsts, ROBO_PAD)
		for idx := 0; idx < 24; idx++ {
			seconds = getCommandsToOutputStringsMultipleStrings(seconds, ROBO_PAD)
			newSeconds := []string{}
			bestLen := datastructures.Pow(10, 9)
			for _, res := range seconds {
				if len(res) < bestLen {
					bestLen = len(res)
				}
			}

			for _, r := range seconds {
				if len(r) == bestLen {
					newSeconds = append(newSeconds, r)
				}
			}
			seconds = newSeconds
			fmt.Println(idx, len(seconds))
		}

		bestLen := datastructures.Pow(10, 9)

		for _, res := range seconds {
			if len(res) < bestLen {
				bestLen = len(res)
			}
		}
		s += bestLen * atoi(r)
	}
	return s
}

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

	if arg == "1" {
		println(solvePartOne(data))
	} else {
		println(solvePartTwo(data))
	}
}
