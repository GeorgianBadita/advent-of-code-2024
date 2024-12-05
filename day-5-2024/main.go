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

func makeGraphAndUpdates(data []string) (map[int][]int, [][]int) {
	graph := map[int][]int{}
	updates := [][]int{}
	for _, line := range data {
		if len(line) < 3 {
			continue
		}
		if strings.Contains(line, "|") {
			numbers := strings.Split(line, "|")
			source, _ := strconv.Atoi(numbers[0])
			dest, _ := strconv.Atoi(numbers[1])

			neighbours, ok := graph[source]
			if ok {
				neighbours = append(neighbours, dest)
				graph[source] = neighbours
			} else {
				graph[source] = []int{dest}
			}
		} else {
			numbers := strings.Split(line, ",")
			numList := []int{}
			for _, num := range numbers {
				n, _ := strconv.Atoi(num)
				numList = append(numList, n)
			}
			updates = append(updates, numList)
		}
	}

	return graph, updates
}

func getPosInArray(nodes []int, val int) int {
	for idx, n := range nodes {
		if n == val {
			return idx
		}
	}
	return -1
}

func checkIsTopoSorted(graph map[int][]int, nodes []int) bool {
	for idx, node := range nodes {
		neighbours, ok := graph[node]
		if !ok {
			continue
		}
		for _, neighbour := range neighbours {
			jdx := getPosInArray(nodes, neighbour)
			if jdx < 0 {
				continue
			}

			if jdx < idx {
				return false
			}
		}
	}

	return true
}

func solvePartOne(graph map[int][]int, updates [][]int) int {
	sum := 0
	for _, update := range updates {
		if checkIsTopoSorted(graph, update) {
			sum += update[len(update)/2]
		}
	}
	return sum
}

func getTopoSort(graph map[int][]int) []int {
	deg := map[int]int{}
	for node := range graph {
		deg[node] = 0
	}

	for _, edges := range graph {
		for _, neighbour := range edges {
			deg[neighbour] += 1
		}
	}

	queue := datastructures.Queue{}
	for node, inDegree := range deg {
		if inDegree == 0 {
			queue.Enqueue(node)
		}
	}

	topoSort := []int{}
	for len(queue) > 0 {
		node := queue.Dequeue()
		topoSort = append(topoSort, node.(int))
		for _, neighbour := range graph[node.(int)] {
			deg[neighbour] -= 1
			if deg[neighbour] == 0 {
				queue.Enqueue(neighbour)
			}
		}
	}
	return topoSort
}

func fixUpdate(graph map[int][]int, update []int) []int {
	filteredGraph := map[int][]int{}
	for node, edges := range graph {
		pos := getPosInArray(update, node)
		if pos < 0 {
			continue
		}
		// filteredGraph[node] = []int{}
		toKeep := []int{}
		for _, neighbour := range edges {
			pos := getPosInArray(update, neighbour)
			if pos < 0 {
				continue
			}
			toKeep = append(toKeep, neighbour)
		}
		filteredGraph[node] = toKeep
	}

	return getTopoSort(filteredGraph)
}

func solvePartTwo(graph map[int][]int, updates [][]int) int {
	sum := 0
	for _, update := range updates {
		if checkIsTopoSorted(graph, update) {
			continue
		}
		update = fixUpdate(graph, update)
		sum += update[len(update)/2]
	}
	return sum
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	graph, updates := makeGraphAndUpdates(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(graph, updates))
	} else {
		println(solvePartTwo(graph, updates))
	}
}
