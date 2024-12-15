package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
)

type Point struct {
	x int
	y int
}

type DoubleCell struct {
	cell1 Point
	cell2 Point
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

func applyDir(grid *[][]rune, currPoint Point, dir DirType, advanceCondition func(rune) bool) Point {
	offset := DIR_TO_OFFSET[dir]

	cP := Point{currPoint.x + offset.x, currPoint.y + offset.y}

	for advanceCondition((*grid)[cP.x][cP.y]) {
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

func applyDirPartTwo(grid *[][]rune, currPoint Point, dir DirType) Point {
	offset := DIR_TO_OFFSET[dir]
	nextNode := Point{currPoint.x + offset.x, currPoint.y + offset.y}
	var box DoubleCell

	if (*grid)[nextNode.x][nextNode.y] == '.' {
		(*grid)[nextNode.x][nextNode.y] = '@'
		(*grid)[currPoint.x][currPoint.y] = '.'
		return nextNode
	}

	if (*grid)[nextNode.x][nextNode.y] == '[' {
		box = DoubleCell{nextNode, Point{nextNode.x, nextNode.y + 1}}
	} else {
		box = DoubleCell{Point{nextNode.x, nextNode.y - 1}, nextNode}
	}

	queue := datastructures.Queue{box}

	allNodesRespectCondition := true

	boxesToMove := map[DoubleCell]bool{}

	for !queue.IsEmpty() {
		n := queue.Dequeue()
		node := n.(DoubleCell)

		if isEmptySpace(*grid, node) {
			continue
		}

		fmt.Println(node, string((*grid)[node.cell1.x][node.cell1.y]), string((*grid)[node.cell2.x][node.cell2.y]))
		neighbours := getBoxNeighboursVertical(*grid, node, dir)
		fmt.Println(neighbours)

	out:
		for _, adj := range neighbours {
			if isBox(*grid, adj) {
				queue.Enqueue(adj)
			}

			if containsWall(*grid, adj) {
				allNodesRespectCondition = false
				break out
			}

			if isEmptySpace(*grid, adj) {
				boxesToMove[adj] = true
			}
			// fmt.Println(adj, string((*grid)[adj.cell1.x][adj.cell1.y]), string((*grid)[adj.cell2.x][adj.cell2.y]))
		}
	}

	if !allNodesRespectCondition {
		return currPoint
	}

	fmt.Println(allNodesRespectCondition)
	fmt.Println(boxesToMove)
	return nextNode
}

func isEmptySpace(grid [][]rune, b DoubleCell) bool {
	return isLeftEmptySpace(grid, b) && isRightEmptySpace(grid, b)
}

func isLeftEmptySpace(grid [][]rune, b DoubleCell) bool {
	return grid[b.cell1.x][b.cell1.y] == '.'
}

func isRightEmptySpace(grid [][]rune, b DoubleCell) bool {
	return grid[b.cell2.x][b.cell2.y] == '.'
}

func isBox(grid [][]rune, b DoubleCell) bool {
	return grid[b.cell1.x][b.cell1.y] == '[' && grid[b.cell2.x][b.cell2.y] == ']'
}

func containsWall(grid [][]rune, b DoubleCell) bool {
	return grid[b.cell1.x][b.cell1.y] == '#' && grid[b.cell2.x][b.cell2.y] == '#'
}

func getBoxNeighboursVertical(grid [][]rune, b DoubleCell, d DirType) []DoubleCell {
	if d == LEFT || d == RIGHT {
		panic("Direction should be UP or DOWN here")
	}
	cell1V := grid[b.cell1.x][b.cell1.y]
	cell2V := grid[b.cell2.x][b.cell2.y]

	offset := DIR_TO_OFFSET[d]

	directOffsetCell1 := Point{b.cell1.x + offset.x, b.cell1.y + offset.y}
	directOffsetCell2 := Point{b.cell2.x + offset.x, b.cell2.y + offset.y}

	// UP / DOWN simple case
	// .. [] []
	// [] .. []
	if cell1V == grid[directOffsetCell1.x][directOffsetCell1.y] && cell2V == grid[directOffsetCell2.x][directOffsetCell2.y] {
		return []DoubleCell{{cell1: directOffsetCell1, cell2: directOffsetCell2}}
	}

	// Complex cases
	// [][] []    []
	//  []   []  []
	leftNeighbour := DoubleCell{cell1: Point{directOffsetCell1.x, directOffsetCell1.y - 1}, cell2: directOffsetCell1}
	rightNeighbour := DoubleCell{cell1: directOffsetCell2, cell2: Point{directOffsetCell2.x, directOffsetCell2.y + 1}}
	return []DoubleCell{leftNeighbour, rightNeighbour}
}

func partOneAdvance(c rune) bool {
	return c == 'O'
}

func partTwoAdvance(c rune) bool {
	return c == '[' || c == ']'
}

func solvePartOne(grid [][]rune, startPoint Point, instructions string) int {
	for _, instruction := range instructions {
		dir := DIR_MAPPING[instruction]
		startPoint = applyDir(&grid, startPoint, dir, partOneAdvance)
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
		dir := DIR_MAPPING[instruction]
		if dir == LEFT || dir == RIGHT {
			startPoint = applyDir(&grid, startPoint, dir, partTwoAdvance)
		} else {
			startPoint = applyDirPartTwo(&grid, startPoint, dir)
		}

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

	grid, _, instructions := parseInput(data)

	// if len(os.Args) != 2 {
	// 	panic("Exactly one arg is expected")
	// }
	// arg := os.Args[1]

	// if arg != "1" && arg != "2" {
	// 	panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	// }

	// if arg == "1" {
	// 	println(solvePartOne(grid, startPoint, instructions))
	// } else {
	newGrid, newStart := doubleGrid(grid)
	println(solvePartTwo(newGrid, newStart, instructions))
	// }
}
