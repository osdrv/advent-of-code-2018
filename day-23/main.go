package main

import (
	"os"
	"strings"

	"github.com/Tom-Johnston/mamba/graph"
)

// This solution is taken from https://www.reddit.com/r/adventofcode/comments/a8s17l/2018_day_23_solutions/
// in particular, I'm using this impl: https://github.com/blu3r4y/AdventOfCode2018/blob/master/src/day23.py
// the original code was giving me an off-by-1 error

// incorrect: 112997633
// correct:   112997634

func parseRobot(s string) [4]int {
	ss := strings.SplitN(s[5:], ">, r=", 2)
	coords := parseInts(ss[0])
	r := parseInt(ss[1])
	return [4]int{coords[0], coords[1], coords[2], r}
}

func computeDist(r1, r2 [4]int) int {
	return abs(r1[0]-r2[0]) + abs(r1[1]-r2[1]) + abs(r1[2]-r2[2])
}

func isRobotsConnected(robots [][4]int, i, j int) bool {
	return computeDist(robots[i], robots[j]) <= (robots[i][3] + robots[j][3])
}

func buildGraph(robots [][4]int) [][]bool {
	graph := make([][]bool, len(robots))
	for i := 0; i < len(robots); i++ {
		graph[i] = make([]bool, len(robots))
	}
	for i := 0; i < len(robots); i++ {
		for j := i + 1; j < len(robots); j++ {
			if isRobotsConnected(robots, i, j) {
				graph[i][j] = true
				graph[j][i] = true
			}
		}
	}
	return graph
}

type Set[T comparable] struct {
	items map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		items: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(item T) {
	s.items[item] = struct{}{}
}

func (s *Set[T]) Remove(item T) {
	delete(s.items, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.items[item]
	return ok
}

func (s *Set[T]) Size() int {
	return len(s.items)
}

func (s *Set[T]) Union(other *Set[T]) *Set[T] {
	ns := NewSet[T]()
	for item := range s.items {
		ns.Add(item)
	}
	for item := range other.items {
		ns.Add(item)
	}
	return ns
}

func (s *Set[T]) Copy() *Set[T] {
	ns := NewSet[T]()
	for item := range s.items {
		ns.Add(item)
	}
	return ns
}

func (s *Set[T]) Items() []T {
	items := make([]T, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}

func (s *Set[T]) Intersect(other *Set[T]) *Set[T] {
	ns := NewSet[T]()
	one, another := s, other
	if another.Size() < s.Size() {
		one, another = another, one
	}
	for item := range one.items {
		if _, ok := another.items[item]; ok {
			ns.Add(item)
		}
	}
	return ns
}

func computeCliques(graph [][]bool) [][]int {
	R := NewSet[int]()
	X := NewSet[int]()
	P := NewSet[int]()
	for i := 0; i < len(graph); i++ {
		P.Add(i)
	}
	N := make(map[int]*Set[int])
	for i := 0; i < len(graph); i++ {
		N[i] = NewSet[int]()
		for j := 0; j < len(graph); j++ {
			if i == j {
				continue
			}
			if graph[i][j] {
				N[i].Add(j)
			}
		}
	}
	cliques := BronKerbosch(N, R, P, X, 0)
	res := make([][]int, 0, len(cliques))
	for _, clique := range cliques {
		res = append(res, clique.Items())
	}
	return res
}

func BronKerbosch(N map[int]*Set[int], R, P, X *Set[int], depth int) []*Set[int] {
	debugf("depth: %d", depth)
	// Continue until P is empty
	if P.Size() == 0 {
		// if X is empty then report the content of R as a new maximal clique
		if X.Size() == 0 {
			return []*Set[int]{R.Copy()}
		}
		// if it’s not then R contains a subset of an already found clique
		return nil
	}
	var res []*Set[int]
	// Pick a vertex v from P to expand.
	vertices := P.Items()
	for _, vertex := range vertices {
		// Add v to R and remove its non-neighbors from P and X
		R.Add(vertex)
		PuN := P.Intersect(N[vertex])
		XuN := X.Intersect(N[vertex])
		res = append(res, BronKerbosch(N, R, PuN, XuN, depth+1)...)
		// Now backtrack to the last vertex picked and restore P ,R and X as they were before the choice
		R.Remove(vertex)
		// remove the vertex from P and add it to X
		P.Remove(vertex)
		X.Add(vertex)
		// then expand the next vertex
	}
	// If there are no more vertexes in P then backtrack to the superior level.
	return res
}

func buildGraph2(robots [][4]int) graph.Graph {
	gg := graph.NewDense(len(robots), nil)
	for i := 0; i < len(robots); i++ {
		for j := i + 1; j < len(robots); j++ {
			if isRobotsConnected(robots, i, j) {
				gg.AddEdge(i, j)
			}
		}
	}
	return gg
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	robots := make([][4]int, 0, len(lines))
	for _, s := range lines {
		robots = append(robots, parseRobot(s))
	}

	debugf("robots: %+v", robots)

	maxr := -ALOT
	maxix := -1
	for ix, robot := range robots {
		if robot[3] > maxr {
			maxr = robot[3]
			maxix = ix
		}
	}

	debugf("The strongest robot is: %+v", robots[maxix])

	cnt := 0
	for _, robot := range robots {
		dist := computeDist(robots[maxix], robot)
		if dist <= robots[maxix][3] {
			debugf("robot %v is in range", robot)
			cnt++
		}
	}

	printf("The total number of robots in range: %d", cnt)

	//graph := buildGraph(robots)

	gg := buildGraph2(robots)

	//cliques := computeCliques(graph)
	cliques := make(chan []int)
	go func() {
		graph.AllMaximalCliques(gg, cliques)
	}()

	var maxclique []int
	for clique := range cliques {
		printf("clique: %+v", clique)
		if len(clique) > len(maxclique) {
			maxclique = clique
		}
	}

	//var maxclique []int
	//for _, clique := range cliques {
	//	if len(clique) > len(maxclique) {
	//		maxclique = clique
	//	}
	//}

	printf("max clique: %+v", maxclique)

	origin := [4]int{0, 0, 0, 0}
	maxsurf := -ALOT
	for _, ix := range maxclique {
		surf := computeDist(robots[ix], origin) - robots[ix][3] + 1
		if surf > maxsurf {
			maxsurf = surf
		}
	}
	printf("equidistant point distance: %d", maxsurf)
}
