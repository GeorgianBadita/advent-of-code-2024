package main

import (
	"bufio"
	"fmt"
	"os"
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

var DIR_MAPPING = map[rune]DirType{
	'^': UP,
	'>': RIGHT,
	'v': DOWN,
	'<': LEFT,
}

var DIR_TO_OFFSET = map[DirType]Point{
	UP:    {-1, 0},
	RIGHT: {0, 1},
	DOWN:  {1, 0},
	LEFT:  {0, -1},
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

func parseInput(data []string) ([][]rune, Point, string) {
	res := [][]rune{}
	coords := Point{-1, -1}
	instructions := ""
	parsingGrid := true

	for idx, line := range data {
		if len(line) == 0 || line == "\n" {
			parsingGrid = false
			continue
		}

		if parsingGrid {
			runes := []rune(line)
			res = append(res, runes)

			for jdx, chr := range runes {
				if chr == '@' {
					coords.x = idx
					coords.y = jdx
				}
			}
		} else {
			instructions += line
		}
	}

	return res, coords, instructions
}

func doubleGrid(grid [][]rune) ([][]rune, Point) {
	newGrid := [][]rune{}
	newStart := Point{-1, -1}
	for idx := 0; idx < len(grid); idx++ {
		newLine := []rune{}
		for jdx := 0; jdx < len(grid[0]); jdx++ {
			if grid[idx][jdx] == '#' {
				newLine = append(newLine, []rune{'#', '#'}...)
			} else if grid[idx][jdx] == 'O' {
				newLine = append(newLine, []rune{'[', ']'}...)
			} else if grid[idx][jdx] == '.' {
				newLine = append(newLine, []rune{'.', '.'}...)
			} else {
				newLine = append(newLine, []rune{'@', '.'}...)
				newStart = Point{idx, len(newLine) - 2}
			}
		}
		newGrid = append(newGrid, newLine)
	}
	return newGrid, newStart
}

func moveLineFromTo(line *[]rune, from int, to int, dir DirType) {
	if dir == UP || dir == DOWN {
		panic("Dir must be LEFT or RIGHT for line ops")
	}

	invDirOffset := DIR_TO_OFFSET[LEFT]

	if dir == LEFT {
		invDirOffset = DIR_TO_OFFSET[RIGHT]
	}

	for start := from; start != to; start += invDirOffset.y {
		(*line)[start] = (*line)[start+invDirOffset.y]
	}
}

func moveColFromTo(grid *[][]rune, colIdx int, from int, to int, dir DirType) {
	if dir == LEFT || dir == RIGHT {
		panic("Dir must be UP or DOWN for col ops")
	}

	invDirOffset := DIR_TO_OFFSET[UP]

	if dir == UP {
		invDirOffset = DIR_TO_OFFSET[DOWN]
	}

	for start := from; start != to; start += invDirOffset.x {
		(*grid)[start][colIdx] = (*grid)[start+invDirOffset.x][colIdx]
	}

}

func applyDir(grid *[][]rune, currPoint Point, instruction rune, advanceCondition func(rune, rune, DirType) bool) Point {
	dir := DIR_MAPPING[instruction]
	offset := DIR_TO_OFFSET[dir]

	cP := Point{currPoint.x + offset.x, currPoint.y + offset.y}
	firstVal := (*grid)[cP.x][cP.y]

	for advanceCondition((*grid)[cP.x][cP.y], firstVal, dir) {
		cP.x += offset.x
		cP.y += offset.y
	}

	if (*grid)[cP.x][cP.y] != '.' {
		return currPoint
	}

	if dir == RIGHT || dir == LEFT {
		moveLineFromTo(&(*grid)[currPoint.x], cP.y, currPoint.y, dir)
	} else {
		moveColFromTo(grid, currPoint.y, cP.x, currPoint.x, dir)
	}

	(*grid)[currPoint.x][currPoint.y] = '.'

	return Point{currPoint.x + offset.x, currPoint.y + offset.y}
}

func partOneAdvance(c rune, _ rune, _ DirType) bool {
	return c == 'O'
}

func partTwoAdvance(c rune, first rune, dir DirType) bool {
	if c != '[' && c != ']' {
		return false
	}

	if dir == LEFT || dir == RIGHT {
		return true
	}

	return c == first
}

func solvePartOne(grid [][]rune, startPoint Point, instructions string) int {
	for _, instruction := range instructions {
		startPoint = applyDir(&grid, startPoint, instruction, partOneAdvance)
	}

	res := 0
	for idx := 0; idx < len(grid); idx++ {
		for jdx := 0; jdx < len(grid[0]); jdx++ {
			if grid[idx][jdx] == 'O' {
				res += 100*idx + jdx
			}
		}
	}

	return res
}

func checkInvariant(grid [][]rune) bool {
	for _, line := range grid {
		for idx := 0; idx < len(line)-1; idx++ {
			if line[idx] == '[' && line[idx+1] != ']' {
				return false
			}
		}
	}

	return true
}

func solvePartTwo(grid [][]rune, startPoint Point, instructions string) int {
	fmt.Println("Initial")
	for _, row := range grid {
		fmt.Println(string(row))
	}

	// numDirs := 10

	for _, instruction := range instructions {
		fmt.Println("Dir " + string(instruction))
		startPoint = applyDir(&grid, startPoint, instruction, partTwoAdvance)
		for _, row := range grid {
			fmt.Println(string(row))
		}

		if !checkInvariant(grid) {
			break
		}
		// if numDirs < 0 {
		// 	break
		// }
		// numDirs -= 1
	}

	res := 0
	for idx := 0; idx < len(grid); idx++ {
		for jdx := 0; jdx < len(grid[0]); jdx++ {
			if grid[idx][jdx] == '[' {
				res += 100*idx + jdx
			}
		}
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

	grid, startPoint, instructions := parseInput(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(grid, startPoint, instructions))
	} else {
		newGrid, newStart := doubleGrid(grid)
		println(solvePartTwo(newGrid, newStart, instructions))
	}
}
