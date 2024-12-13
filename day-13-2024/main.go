package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Button struct {
	xOffset int
	yOffset int
}

type Point struct {
	x uint64
	y uint64
}

type ChallengeConfig struct {
	aButton     Button
	bButton     Button
	prizeConfig Point
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

func parseOneButton(str string) Button {
	split := strings.Split(str, "+")
	xPart := strings.Split(split[1], ",")[0]
	yPart := split[2]

	x, _ := strconv.Atoi(xPart)
	y, _ := strconv.Atoi(yPart)

	return Button{xOffset: x, yOffset: y}
}

func parseInput(data []string) []ChallengeConfig {
	res := []ChallengeConfig{}
	idx := 0

	for idx < len(data) {
		aButton := parseOneButton(data[idx])
		bButton := parseOneButton(data[idx+1])
		prizeSplit := strings.Split(data[idx+2], ", ")

		xPart := prizeSplit[0][9:]
		yPart := prizeSplit[1][2:]

		x, _ := strconv.Atoi(xPart)
		y, _ := strconv.Atoi(yPart)

		res = append(res, ChallengeConfig{aButton: aButton, bButton: bButton, prizeConfig: Point{uint64(x), uint64(y)}})
		idx += 4
	}

	return res
}

func solvePartOne(challenges []ChallengeConfig) int {
	oneE9 := datastructures.Pow(10, 9)
	totalTokens := 0
	for _, challenge := range challenges {
		bestSol := oneE9
		for a := 1; a <= 100; a++ {
			for b := 1; b <= 100; b++ {
				xRes := a*challenge.aButton.xOffset + b*challenge.bButton.xOffset
				yRes := a*challenge.aButton.yOffset + b*challenge.bButton.yOffset

				if uint64(xRes) == challenge.prizeConfig.x && uint64(yRes) == challenge.prizeConfig.y {
					bestSol = min(bestSol, 3*a+b)
				}
			}
		}

		if bestSol < oneE9 {
			totalTokens += bestSol
		}
	}

	return totalTokens
}

func solvePartTwo(challenges []ChallengeConfig) int {
	totalTokens := 0
	for _, challenge := range challenges {
		challenge.prizeConfig.x += 10000000000000
		challenge.prizeConfig.y += 10000000000000

		// A.x, A.y
		// B.x, B.y

		// x*Ax + y*Bx = Cx | Ay
		// x*Ay + y*By = Cy | Ax

		// x*Ax*Ay + y*Bx*Ay = Cx*Ay
		// x*Ax*Ay + y*By*Ax = Cy*Ax

		// y(Bx*Ay - By*Ax) = Cx*Ay - Cy*Ax

		// y = (Cx*Ay - Cy*Ax) / (Bx*Ay - By*Ax)

		// x = (Cx - y*Bx)/Ax

		yNumerator := int(challenge.prizeConfig.x)*challenge.aButton.yOffset - int(challenge.prizeConfig.y)*challenge.aButton.xOffset
		yDenominator := challenge.aButton.yOffset*challenge.bButton.xOffset - challenge.bButton.yOffset*challenge.aButton.xOffset

		yNumerator = datastructures.Abs(yNumerator)
		yDenominator = datastructures.Abs(yDenominator)

		if yNumerator%yDenominator != 0 {
			continue
		}

		y := yNumerator / yDenominator

		xNumerator := int(challenge.prizeConfig.x) - y*(challenge.bButton.xOffset)
		xNumerator = datastructures.Abs(xNumerator)

		if xNumerator%challenge.aButton.xOffset != 0 {
			continue
		}

		x := xNumerator / challenge.aButton.xOffset

		totalTokens += 3*x + y
	}

	return totalTokens
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	challenges := parseInput(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(challenges))
	} else {
		println(solvePartTwo(challenges))
	}
}
