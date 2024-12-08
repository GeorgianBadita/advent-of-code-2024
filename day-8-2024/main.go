package main

import (
	"bufio"
	"os"
)

type Point struct {
	x int
	y int
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

func getSymetricOfARelativeToB(a Point, b Point) Point {
	// (a.x + s.x) / 2 = b.x => a.x + s.x = 2*b.x => s.x = 2*b.x - a.x
	// (a.y + s.y) / 2 = b.y => a.y + s.y = 2*b.y => s.y = 2*b.y - a.y

	return Point{2*b.x - a.x, 2*b.y - a.y}
}

func testPointIsOnLineAB(a Point, b Point, c Point) bool {
	/*
	   y - yA = m(x - xA)
	   m = (yB - yA)/(xB - xA)

	   => y - yA - m(x - xA) = 0
	   => y - yA - (yB - yA)(x - xA)/(xB-xA) = 0
	*/

	ybya := b.y - a.y
	xxa := c.x - a.x
	xbxa := b.x - a.x

	return float64(c.y-a.y)-float64(ybya)*float64(xxa)/float64(xbxa) == 0
}

func validCoords(matrix [][]rune, p Point) bool {
	return p.x >= 0 && p.x < len(matrix) && p.y >= 0 && p.y < len(matrix[0])
}

func collectAntenas(matrix [][]rune) map[rune][]Point {
	res := map[rune][]Point{}
	for idx := 0; idx < len(matrix); idx++ {
		for jdx := 0; jdx < len(matrix[0]); jdx++ {
			charAt := matrix[idx][jdx]
			if charAt == '.' {
				continue
			}
			if val, ok := res[charAt]; ok {
				val = append(val, Point{idx, jdx})
				res[charAt] = val
			} else {
				res[charAt] = []Point{{idx, jdx}}
			}
		}
	}
	return res
}

func solvePartOne(matrix [][]rune) int {
	antenas := collectAntenas(matrix)
	seen := map[Point]bool{}
	for _, points := range antenas {
		for idx := 0; idx < len(points)-1; idx++ {
			for jdx := idx + 1; jdx < len(points); jdx++ {
				iSimJ := getSymetricOfARelativeToB(points[idx], points[jdx])
				jSimI := getSymetricOfARelativeToB(points[jdx], points[idx])

				if validCoords(matrix, iSimJ) {
					if _, ok := seen[iSimJ]; !ok {
						seen[iSimJ] = true
					}
				}

				if validCoords(matrix, jSimI) {
					if _, ok := seen[jSimI]; !ok {
						seen[jSimI] = true
					}
				}
			}
		}
	}

	return len(seen)
}

func solvePartTwo(matrix [][]rune) int {
	antenas := collectAntenas(matrix)
	seen := map[Point]bool{}
	for _, points := range antenas {
		for idx := 0; idx < len(points)-1; idx++ {
			for jdx := idx + 1; jdx < len(points); jdx++ {
				for r := 0; r < len(matrix); r++ {
					for c := 0; c < len(matrix[0]); c++ {
						pointOnLine := testPointIsOnLineAB(points[idx], points[jdx], Point{r, c})
						if pointOnLine {
							if _, ok := seen[Point{r, c}]; !ok {
								seen[Point{r, c}] = true
							}
						}
					}
				}
			}
		}
	}

	return len(seen)
}

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	matrix := [][]rune{}
	for _, line := range data {
		matrix = append(matrix, []rune(line))
	}

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(matrix))
	} else {
		println(solvePartTwo(matrix))
	}
}
