package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Span struct {
	start int
	end   int
}

func (s1 Span) CompareWith(s datastructures.ElementWithPriority) int {
	s2 := s.(Span)
	if s1.start < s2.start {
		return -1
	}

	if s2.start < s1.start {
		return 1
	}

	if s1 == s2 {
		return 0
	}
	panic("There should be no spans starting from the same point, unless it's the same span")
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

func expandString(memory []rune) []int {
	// Meaning:
	// v[x] = -1, if xth space is empty
	// v[x] = y otherwise, where xth space has a memory page for id y

	res := []int{}
	memId := 0
	for idx, chr := range memory {
		num, _ := strconv.Atoi(string(chr))
		if idx%2 == 0 {
			for jdx := 0; jdx < num; jdx++ {
				res = append(res, memId)
			}
			memId += 1
		} else {
			for jdx := 0; jdx < num; jdx++ {
				res = append(res, -1)
			}
		}
	}

	return res
}

func printMemoryArr(memory []int) {
	str := ""
	for _, elem := range memory {
		if elem == -1 {
			str += "."
		} else {
			str += fmt.Sprintf("%d", elem)
		}
	}
	fmt.Println(str)
}

func solvePartOne(memory []rune) uint {
	expandedMemory := expandString(memory)
	pointPtr := 0
	for pointPtr < len(expandedMemory) && expandedMemory[pointPtr] != -1 {
		pointPtr++
	}

	memoryPtr := len(expandedMemory) - 1
	for memoryPtr >= 0 && expandedMemory[memoryPtr] == -1 {
		memoryPtr -= 1
	}

	for pointPtr < memoryPtr {
		expandedMemory[pointPtr], expandedMemory[memoryPtr] = expandedMemory[memoryPtr], expandedMemory[pointPtr]
		for pointPtr < len(expandedMemory) && expandedMemory[pointPtr] != -1 {
			pointPtr++
		}
		for memoryPtr >= 0 && expandedMemory[memoryPtr] == -1 {
			memoryPtr -= 1
		}
	}

	expandedMemory = expandedMemory[:pointPtr]
	total_res := uint(0)
	for idx, val := range expandedMemory {
		total_res += uint(idx) * uint(val)
	}

	return uint(total_res)
}

func assertZones(zones map[int]datastructures.PriorityQueue, mem []int) {
	for _, spans := range zones {
		for _, s := range spans {
			span := s.(Span)
			spanArr := mem[span.start : span.end+1]
			for _, v := range spanArr {
				if v != -1 {
					panic("Incorrect zone segments")
				}
			}
			if span.end+1 < len(mem) && mem[span.end+1] == -1 {
				panic("Incorrect zone segments")
			}
			if span.start-1 >= 0 && mem[span.start-1] == -1 {
				panic("Incorrect zone segments")
			}
		}
	}
}

func solvePartTwo(memory []rune) uint64 {
	emptyZones := map[int]datastructures.PriorityQueue{}
	expandedMemory := expandString(memory)

	pointPtr := 0
	for pointPtr < len(expandedMemory) && expandedMemory[pointPtr] != -1 {
		pointPtr++
	}

	zoneStart := pointPtr
	for idx := pointPtr + 1; idx < len(expandedMemory); idx++ {
		if expandedMemory[idx] == -1 {
			continue
		} else {
			zoneEnd := idx - 1
			size := zoneEnd - zoneStart + 1

			v, ok := emptyZones[size]
			if ok {
				v.Insert(Span{zoneStart, zoneEnd})
				emptyZones[size] = v
			} else {
				nq := datastructures.PriorityQueue{}
				nq.Insert(Span{zoneStart, zoneEnd})
				emptyZones[size] = nq
			}

			for idx < len(expandedMemory) && expandedMemory[idx] != -1 {
				idx++
			}
			zoneStart = idx
		}
	}

	assertZones(emptyZones, expandedMemory)
	// Normally we should do another check here, but the input has an odd length
	// si it cannot end with .....

	// fmt.Println(zoneStart)
	// printMemoryArr(expandedMemory)
	// fmt.Println(emptyZones)

	for memoryPtr := len(expandedMemory) - 1; memoryPtr >= 0; {
		regionEnd := memoryPtr + 1
		for memoryPtr >= 0 && expandedMemory[memoryPtr] == expandedMemory[regionEnd-1] {
			memoryPtr -= 1
		}

		printMemoryArr(expandedMemory[memoryPtr+1 : regionEnd])

		size := regionEnd - memoryPtr - 1
		bestDiff := datastructures.Pow(10, 9)
		bestSize := -1

		for zoneSize := range emptyZones {
			if zoneSize >= size && zoneSize-size < bestDiff {
				bestDiff = zoneSize - size
				bestSize = zoneSize
			}
		}

		if bestSize != -1 {
			zones := emptyZones[bestSize]
			z := *zones.Peek()
			zone := z.(Span)

			if zone.start >= memoryPtr+1 {
				for memoryPtr >= 0 && expandedMemory[memoryPtr] == -1 {
					memoryPtr -= 1
				}
				continue
			}

			if len(zones) == 1 {
				delete(emptyZones, bestSize)
			} else {
				zones.Remove()
				emptyZones[bestSize] = zones
			}

			if bestDiff > 0 {
				v, ok := emptyZones[bestDiff]
				if ok {
					newSpan := Span{zone.start + size, zone.end}
					v.Insert(newSpan)
					emptyZones[bestDiff] = v
				} else {
					nq := datastructures.PriorityQueue{}
					nq.Insert(Span{zone.start + size, zone.end})
					emptyZones[bestDiff] = nq
				}
			}

			for idx := 0; idx < size; idx++ {
				expandedMemory[zone.start+idx] = expandedMemory[memoryPtr+1+idx]
				expandedMemory[memoryPtr+1+idx] = -1
			}

		}

		for memoryPtr >= 0 && expandedMemory[memoryPtr] == -1 {
			memoryPtr -= 1
		}
	}

	printMemoryArr(expandedMemory)

	total_res := uint64(0)
	for idx, val := range expandedMemory {
		if val != -1 {
			total_res += uint64(idx) * uint64(val)
		}

	}

	// printMemoryArr(expandedMemory)
	return total_res
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	memory := []rune(data[0])

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(memory))
	} else {
		println(solvePartTwo(memory))
	}
}
