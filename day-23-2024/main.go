package main

import (
	"aoc-2024/datastructures"
	"bufio"
	"os"
	"slices"
	"strings"
)

type Node string
type Graph map[Node][]Node

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

func makeGraph(data []string) map[Node][]Node {
	graph := map[Node][]Node{}

	for _, line := range data {
		split := strings.Split(line, "-")
		graph[Node(split[0])] = append(graph[Node(split[0])], Node(split[1]))
		graph[Node(split[1])] = append(graph[Node(split[1])], Node(split[0]))
	}

	return graph
}

func bfs(graph map[Node][]Node, start Node, visited map[Node]bool) []Node {
	visited[start] = true
	nodesVisitedNow := []Node{start}

	queue := datastructures.Queue{start}

	for !queue.IsEmpty() {
		n := queue.Dequeue()
		node := n.(Node)

		neighbours := graph[node]

		for idx := 0; idx < len(neighbours); idx++ {
			adj := neighbours[idx]
			if _, ok := visited[adj]; !ok {
				visited[adj] = true
				nodesVisitedNow = append(nodesVisitedNow, adj)
				queue.Enqueue(adj)
			}
		}
	}
	return nodesVisitedNow
}

func contains(a Node, b Node, c Node, graph map[Node][]Node) bool {
	return slices.Contains(graph[a], b) && slices.Contains(graph[a], c) &&
		slices.Contains(graph[b], a) && slices.Contains(graph[b], c) &&
		slices.Contains(graph[c], a) && slices.Contains(graph[c], b)
}

func bronKerbosch(R Graph, P Graph, X Graph, original Graph, cliques *[]Graph) {
	if len(P) == 0 && len(X) == 0 {
		(*cliques) = append((*cliques), R)
		return
	} else {
		for node := range P {
			RCp := Graph{}
			RCp[node] = []Node{}
			for n := range R {
				RCp[n] = []Node{}
			}

			PInt := Graph{}
			XInt := Graph{}
			for _, adj := range original[node] {
				if _, ok := P[adj]; ok {
					PInt[adj] = []Node{}
				}
				if _, ok := X[adj]; ok {
					XInt[adj] = []Node{}
				}
			}
			bronKerbosch(RCp, PInt, XInt, original, cliques)
			delete(P, node)
			X[node] = []Node{}
		}
	}
}

func solvePartOne(graph map[Node][]Node) int {
	res := 0
	visited := map[Node]bool{}
	for node := range graph {
		if _, ok := visited[node]; !ok {
			component := bfs(graph, node, visited)
			for idx := 0; idx < len(component)-2; idx++ {
				for jdx := idx + 1; jdx < len(component)-1; jdx++ {
					for kdx := jdx + 1; kdx < len(component); kdx++ {

						if !contains(component[idx], component[jdx], component[kdx], graph) {
							continue
						}

						if strings.HasPrefix(string(component[idx]), "t") ||
							strings.HasPrefix(string(component[jdx]), "t") ||
							strings.HasPrefix(string(component[kdx]), "t") {
							res++
						}
					}
				}
			}
		}
	}
	return res
}

func solvePartTwo(graph map[Node][]Node) string {
	R := Graph{}
	P := Graph{}
	for node := range graph {
		P[node] = []Node{}
	}

	cliques := []Graph{}

	bronKerbosch(R, P, Graph{}, graph, &cliques)

	bestLen := 0
	bestLenStr := ""
	for _, clique := range cliques {
		if len(clique) >= bestLen {
			bestLenNodes := []string{}
			for node := range clique {
				bestLenNodes = append(bestLenNodes, string(node))
			}
			slices.Sort(bestLenNodes)
			currStr := strings.Join(bestLenNodes, ",")
			if bestLenStr == "" || currStr < bestLenStr || len(clique) > bestLen {
				bestLen = len(clique)
				bestLenStr = currStr
			}
		}
	}

	return bestLenStr
}

const IN_FILE_PATH = "./input.txt"

func main() {
	data, err := readFileAsString(IN_FILE_PATH)

	if err != nil {
		println(err)
		panic("Error reading from " + IN_FILE_PATH)
	}

	graph := makeGraph(data)

	if len(os.Args) != 2 {
		panic("Exactly one arg is expected")
	}
	arg := os.Args[1]

	if arg != "1" && arg != "2" {
		panic("Arg can only be 1 or 2 for part 1 ore part 2 of the problem respectively")
	}

	if arg == "1" {
		println(solvePartOne(graph))
	} else {
		println(solvePartTwo(graph))
	}
}
