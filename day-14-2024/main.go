package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

type Robot struct {
	startPoint Point
	velocity   Point
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

func parseLine(line string) Robot {
	parts := strings.Split(line, "=")
	posStr := strings.Split(strings.Split(parts[1], " ")[0], ",")
	velocityStr := strings.Split(parts[2], ",")

	posX, _ := strconv.Atoi(posStr[0])
	posY, _ := strconv.Atoi(posStr[1])

	velocityX, _ := strconv.Atoi(velocityStr[0])
	velocityY, _ := strconv.Atoi(velocityStr[1])

	return Robot{
		startPoint: Point{x: posY, y: posX},
		velocity:   Point{x: velocityY, y: velocityX},
	}
}

func paraseInputToRobots(data []string) []Robot {
	robots := []Robot{}

	for _, line := range data {
		robots = append(robots, parseLine(line))
	}

	return robots
}

func printInverted(robots []Robot, numRows int, numCols int) {
	coords := []Point{}

	for _, r := range robots {
		coords = append(coords, Point{r.startPoint.x, r.startPoint.y})
	}

	for idx := 0; idx < numRows; idx++ {
		for jdx := 0; jdx < numCols; jdx++ {
			numRobots := 0
			for _, r := range coords {
				if r.x == idx && r.y == jdx {
					numRobots++
				}
			}
			if numRobots > 0 {
				fmt.Print(strconv.Itoa(numRobots))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func positiveMod(n int, mod int) int {
	modRes := n % mod
	if modRes < 0 {
		modRes += mod
	}

	return modRes
}

func advanceOneSecond(robots *[]Robot, numRows int, numCols int) {
	for idx := 0; idx < len(*robots); idx++ {
		cx := (*robots)[idx].startPoint.x
		cy := (*robots)[idx].startPoint.y

		cx = positiveMod(cx+(*robots)[idx].velocity.x, numRows)
		cy = positiveMod(cy+(*robots)[idx].velocity.y, numCols)

		(*robots)[idx].startPoint.x = cx
		(*robots)[idx].startPoint.y = cy
	}
}

func solvePartOne(robots []Robot, seconds int, numRows int, numCols int) int {
	for seconds > 0 {
		advanceOneSecond(&robots, numRows, numCols)
		seconds--
	}

	halfRows := numRows / 2
	halfCols := numCols / 2

	areasRowBounds := []Point{{0, halfRows}, {0, halfRows}, {halfRows + 1, numRows}, {halfRows + 1, numRows}}
	areasColBounds := []Point{{0, halfCols}, {halfCols + 1, numCols}, {0, halfCols}, {halfCols + 1, numCols}}

	safeArea := 0

	for idx := 0; idx < len(areasColBounds); idx++ {
		rowStart := areasRowBounds[idx].x
		rowEnd := areasRowBounds[idx].y

		colStart := areasColBounds[idx].x
		colEnd := areasColBounds[idx].y

		numRobots := 0
		for row := rowStart; row < rowEnd; row++ {
			for col := colStart; col < colEnd; col++ {
				for _, r := range robots {
					if r.startPoint.x == row && r.startPoint.y == col {
						numRobots++
					}
				}
			}
		}

		if safeArea == 0 {
			safeArea = numRobots
		} else {
			safeArea *= numRobots
		}
	}

	return safeArea
}

func solvePartTwo(robots []Robot, seconds int, numRows int, numCols int) {
	for idx := 1; idx <= seconds; idx++ {
		advanceOneSecond(&robots, numRows, numCols)
		perX := map[int][]Point{}

		for _, r := range robots {
			perX[r.startPoint.x] = append(perX[r.startPoint.x], r.startPoint)
		}

		for _, r := range perX {
			if len(r) > 30 {
				fmt.Printf("After %d seconds\n", idx)
				printInverted(robots, numRows, numCols)
			}
		}
	}
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	robots := paraseInputToRobots(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(robots, 100, 103, 101))
	} else {
		solvePartTwo(robots, 10000, 103, 101)
	}
}
